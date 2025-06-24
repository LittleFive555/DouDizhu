package network

import (
	"DouDizhuServer/errordef"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/serialize"
	"DouDizhuServer/network/session"
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

var Server *GameServer

// GetServer 返回游戏服务器实例
func GetServer() *GameServer {
	return Server
}

type GameServer struct {
	listener   net.Listener
	sessionMgr *session.SessionManager
	dispatcher *message.MessageDispatcher
}

func NewGameServer() *GameServer {
	gameServer := &GameServer{}
	gameServer.sessionMgr = session.NewSessionManager()
	gameServer.dispatcher = message.NewMessageDispatcher(10, gameServer.handleMessage)
	return gameServer
}

// Start 启动TCP服务器
func (s *GameServer) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("启动服务器失败: %v", err)
	}
	s.listener = ln

	logger.InfoWith("TCP服务器启动成功", "addr", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.ErrorWith("接受连接失败", "error", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Stop 停止TCP服务器
func (s *GameServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *GameServer) RegisterHandler(msgId protodef.PMsgId, handler func(*message.MessageContext, *proto.Message) (*message.HandleResult, error)) {
	s.dispatcher.RegisterHandler(msgId, handler)
}

func HandleHandshake(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PHandshakeRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}

	session, err := Server.sessionMgr.GetSession(context.SessionId)
	if err != nil {
		return nil, err
	}

	serverPublicKeyStr, err := session.GenerateSecureKey(reqMsg.GetPublicKey(), reqMsg.GetSalt(), reqMsg.GetInfo())
	if err != nil {
		return nil, err
	}

	return &message.HandleResult{
		Resp: &protodef.PHandshakeResponse{
			PublicKey: serverPublicKeyStr,
		},
	}, nil
}

// handleConnection 处理单个连接
func (s *GameServer) handleConnection(conn net.Conn) {
	session, err := s.sessionMgr.CreatePlayerSession(conn)
	if err != nil {
		logger.ErrorWith("创建会话失败", "error", err)
		conn.Close()
		return
	}
	logger.InfoWith("创建会话成功，开始处理消息", "sessionId", session.Id)
	defer s.sessionMgr.CloseSession(session.Id)
	session.StartReading(s.dispatcher.EnqueueMessage)
}

// HandleMessage 实现Handler接口
func (s *GameServer) handleMessage(msg *message.Message) {
	// TODO 这里产生的error，需要返回给客户端
	clientMsg, err := serialize.Deserialize(msg.Data)
	if err != nil {
		logger.ErrorWith("解析消息失败", "error", err)
		return
	}

	session, err := s.sessionMgr.GetSession(msg.SessionId)
	if err != nil {
		logger.ErrorWith("获取会话失败", "error", err)
		return
	}

	// 消息体解密并反序列化
	enableEncryption := session.IsSecureKeyValid()
	msgHeader := clientMsg.GetHeader()
	msgId := msgHeader.GetMsgId()
	var reqPayloadBytes []byte
	if enableEncryption {
		logger.InfoWith("解密消息", "payload", base64.StdEncoding.EncodeToString(clientMsg.GetPayload()), "iv", msgHeader.Iv)
		reqPayloadBytes, err = session.DecryptPayload(clientMsg.GetPayload(), msgHeader.Iv)
		if err != nil {
			logger.ErrorWith("解密消息失败", "error", err)
			return
		}
		logger.InfoWith("解密消息成功", "payload", base64.StdEncoding.EncodeToString(reqPayloadBytes))
	} else {
		reqPayloadBytes = clientMsg.GetPayload()
	}
	reqPayload, err := serialize.DeserializePayload(msgId, reqPayloadBytes)
	if err != nil {
		logger.ErrorWith("反序列化消息失败", "error", err)
		return
	}

	if isSensitiveMessage(msgId) {
		logger.InfoWith("收到消息", "类型", msgId, "sessionId", msg.SessionId)
	} else {
		logger.InfoWith("收到消息", "类型", msgId, "sessionId", msg.SessionId, "payload", reqPayload)
	}

	// 处理消息
	handler := s.dispatcher.GetHandler(msgId)
	if handler == nil {
		logger.ErrorWith("未找到消息处理器", "type", msgId)
		return
	}
	context := &message.MessageContext{
		SessionId: msg.SessionId,
		PlayerId:  msgHeader.PlayerId,
	}
	result, err := handler(context, &reqPayload)

	respMessage := createResponseMsg(msgHeader)
	var respPayload proto.Message
	var notifyPayload proto.Message

	// 处理错误
	var gameError *errordef.GameError
	if gameError = errordef.AsGameError(err); gameError != nil {
		if gameError.Category == errordef.CategoryGameplay {
			logger.InfoWith("游戏逻辑错误", "errorCode", gameError.Code, "errorMessage", gameError.ClientMsg)
		} else {
			logger.ErrorWith("服务器错误", "errorCategory", gameError.Category, "errorCode", gameError.Code, "errorMessage", gameError.ClientMsg)
		}
		respPayload = message.CreateErrorPayload(gameError)
		respMessage.MsgType = protodef.PServerMsgType_PSERVER_MSG_TYPE_ERROR
		notifyPayload = nil
	} else {
		respMessage.MsgType = protodef.PServerMsgType_PSERVER_MSG_TYPE_RESPONSE
		respPayload = result.Resp
		notifyPayload = result.Notify
	}

	// 序列化和加密响应
	respPayloadBytes, err := serialize.SerializePayload(respPayload)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return
	}
	if enableEncryption {
		var iv []byte
		logger.InfoWith("加密响应前", "payload", base64.StdEncoding.EncodeToString(respPayloadBytes))
		respPayloadBytes, iv, err = session.EncryptPayload(respPayloadBytes)
		logger.InfoWith("加密响应后", "payload", base64.StdEncoding.EncodeToString(respPayloadBytes), "iv", iv)
		respMessage.Header.Iv = iv
	}
	if err != nil {
		logger.ErrorWith("加密响应失败", "error", err)
		return
	}
	respMessage.Payload = respPayloadBytes
	responseData, err := serialize.Serialize(respMessage)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return
	}
	err = session.SendMessage(responseData)
	if err != nil {
		logger.ErrorWith("发送消息失败", "error", err)
		return
	}
	if isSensitiveMessage(msgId) {
		logger.InfoWith("发送消息成功", "msgId", msgId, "sessionId", msg.SessionId)
	} else {
		logger.InfoWith("发送消息成功", "msgId", msgId, "sessionId", msg.SessionId, "payload", respPayload)
	}

	// 发送通知
	// 包装为 GameServerMessage
	if gameError == nil && notifyPayload != nil {
		notifyPayloadBytes, err := serialize.SerializePayload(notifyPayload)
		if err != nil {
			logger.ErrorWith("序列化通知失败", "error", err)
			return
		}
		allSessions := s.sessionMgr.GetAllSessions()
		for _, session := range allSessions {
			notificationPayloadBytes, iv, err := session.EncryptPayload(notifyPayloadBytes)
			if err != nil {
				logger.ErrorWith("加密通知失败", "error", err)
				continue
			}
			notificationMessage := createNotificationMsg(session, result.NofityMsgId)
			notificationMessage.Payload = notificationPayloadBytes
			notificationMessage.Header.Iv = iv
			notificationData, err := serialize.Serialize(notificationMessage)
			if err != nil {
				logger.ErrorWith("序列化响应失败", "error", err)
				continue
			}
			err = session.SendMessage(notificationData)
			if err != nil {
				logger.ErrorWith("发送消息失败", "error", err)
				continue
			}
		}
		if isSensitiveMessage(msgId) {
			logger.InfoWith("发送通知结束", "msgId", msgId, "sessionId", msg.SessionId)
		} else {
			logger.InfoWith("发送通知结束", "msgId", msgId, "sessionId", msg.SessionId, "payload", notifyPayload)
		}
	}
}

func isSensitiveMessage(msgId protodef.PMsgId) bool {
	if msgId == protodef.PMsgId_PMSG_ID_REGISTER ||
		msgId == protodef.PMsgId_PMSG_ID_LOGIN ||
		msgId == protodef.PMsgId_PMSG_ID_HANDSHAKE {
		return true
	}
	return false
}

func createResponseMsg(requestHeader *protodef.PMsgHeader) *protodef.PServerMsg {
	return &protodef.PServerMsg{
		Header: &protodef.PMsgHeader{
			UniqueId:  requestHeader.UniqueId,
			MsgId:     requestHeader.MsgId,
			Timestamp: requestHeader.Timestamp,

			SessionId: requestHeader.SessionId,
			PlayerId:  requestHeader.PlayerId,
		},
	}
}

func createNotificationMsg(session *session.PlayerSession, msgId protodef.PMsgId) *protodef.PServerMsg {
	return &protodef.PServerMsg{
		Header: &protodef.PMsgHeader{
			UniqueId:  time.Now().UnixNano(),
			MsgId:     msgId,
			Timestamp: time.Now().UnixMilli(),

			SessionId: session.Id,
			PlayerId:  session.PlayerId,
		},
		MsgType: protodef.PServerMsgType_PSERVER_MSG_TYPE_NOTIFICATION,
	}
}
