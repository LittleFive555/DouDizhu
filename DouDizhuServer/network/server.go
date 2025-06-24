package network

import (
	"DouDizhuServer/errors"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/serialize"
	"DouDizhuServer/network/session"
	"fmt"
	"net"

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
	clientMsg, err := serialize.Deserialize(msg.Data)
	if err != nil {
		logger.ErrorWith("解析消息失败", "error", err)
		return
	}

	// 查找对应的消息处理器
	content := clientMsg.GetPayload()
	if content == nil {
		logger.ErrorWith("消息内容为空")
		return
	}
	msgHeader := clientMsg.GetHeader()
	msgId := msgHeader.GetMsgId()
	msgPayload := serialize.GetMessage(msgId)
	if err := proto.Unmarshal(content, msgPayload); err != nil {
		logger.ErrorWith("解析消息体失败", "msgId", msgId, "error", err)
		return
	}

	if isSecretMessage(msgId) {
		logger.InfoWith("收到消息", "类型", msgId, "sessionId", msg.SessionId)
	} else {
		logger.InfoWith("收到消息", "类型", msgId, "sessionId", msg.SessionId, "message", clientMsg)
	}

	handler := s.dispatcher.GetHandler(msgId)
	if handler == nil {
		logger.ErrorWith("未找到消息处理器", "type", msgId)
		return
	}

	context := &message.MessageContext{
		SessionId: msgHeader.SessionId,
		PlayerId:  msgHeader.PlayerId,
	}
	// 处理消息
	result, err := handler(context, &msgPayload)
	var respMessage *protodef.PServerMsg = nil
	var notificationMessage *protodef.PServerMsg = nil

	if gameError := errors.AsGameError(err); gameError != nil {
		if gameError.Category == errors.CategoryGameplay {
			logger.InfoWith("游戏逻辑错误", "errorCode", gameError.Code, "errorMessage", gameError.ClientMsg)
		} else {
			logger.ErrorWith("服务器错误", "errorCategory", gameError.Category, "errorCode", gameError.Code, "errorMessage", gameError.ClientMsg)
		}
		respMessage = message.CreateErrorMessage(msgHeader, gameError)
	} else {
		respMessage = message.CreateRespMessage(msgHeader, result.Resp)
		notificationMessage = message.CreateNotificationMessage(msgHeader, result.Notify)
	}

	responseData, err := serialize.Serialize(respMessage)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return
	}
	// TODO 对消息进行加密
	session, err := s.sessionMgr.GetSession(msg.SessionId)
	if err != nil {
		logger.ErrorWith("获取会话失败", "error", err)
		return
	}
	err = session.SendMessage(responseData)
	if err != nil {
		logger.ErrorWith("发送消息失败", "error", err)
		return
	}
	if isSecretMessage(msgId) {
		logger.InfoWith("发送消息成功", "msgId", msgId, "sessionId", msg.SessionId)
	} else {
		logger.InfoWith("发送消息成功", "msgId", msgId, "sessionId", msg.SessionId, "message", respMessage)
	}

	// 发送通知
	// 包装为 GameServerMessage
	if notificationMessage != nil {
		notificationData, err := serialize.Serialize(notificationMessage)
		if err != nil {
			logger.ErrorWith("序列化响应失败", "error", err)
			return
		}
		// TODO 对消息进行加密
		allSessions := s.sessionMgr.GetAllSessions()
		for _, session := range allSessions {
			err = session.SendMessage(notificationData)
			if err != nil {
				logger.ErrorWith("发送消息失败", "error", err)
				continue
			}
		}
		if isSecretMessage(msgId) {
			logger.InfoWith("发送通知结束", "msgId", msgId, "sessionId", msg.SessionId)
		} else {
			logger.InfoWith("发送通知结束", "msgId", msgId, "sessionId", msg.SessionId, "message", notificationMessage)
		}
	}
}

func isSecretMessage(msgId protodef.PMsgId) bool {
	if msgId == protodef.PMsgId_PMSG_ID_REGISTER ||
		msgId == protodef.PMsgId_PMSG_ID_LOGIN ||
		msgId == protodef.PMsgId_PMSG_ID_HANDSHAKE {
		return true
	}
	return false
}
