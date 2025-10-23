package network

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network/handler"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"
	"DouDizhuServer/scripts/network/serialize"
	"DouDizhuServer/scripts/network/session"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type GameServer struct {
	listener        net.Listener
	sessionMgr      *session.SessionManager
	messageRegister *message.MessageRegister

	messageQueue      chan *message.Message      // 消息队列
	notificationQueue chan *message.Notification // 通知队列
}

func NewGameServer() *GameServer {
	gameServer := &GameServer{}
	gameServer.sessionMgr = session.NewSessionManager()
	gameServer.messageRegister = message.NewMessageRegister()

	gameServer.messageQueue = make(chan *message.Message, 10000)
	gameServer.notificationQueue = make(chan *message.Notification, 10000)
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

	go s.handleMessages()
	go s.handleNotifications()

	s.handleConnections(ln)

	return nil
}

// Stop 停止TCP服务器
func (s *GameServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *GameServer) RegisterHandlers() {
	// 启动相关
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_HANDSHAKE, s.handleHandshake)

	// 账号相关
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_REGISTER, handler.HandleRegister)
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_LOGIN, handler.HandleLogin)

	// 房间相关
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_CREATE_ROOM, handler.HandleCreateRoom)
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_GET_ROOM_LIST, handler.HandleGetRoomList)
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_ENTER_ROOM, handler.HandleEnterRoom)
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_LEAVE_ROOM, handler.HandleLeaveRoom)

	// 聊天相关
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_CHAT_MSG, handler.HandleChatMessage)
}

func (s *GameServer) RegisterHandler(msgId protodef.PMsgId, handler func(*message.MessageContext, *proto.Message) (*message.HandleResult, error)) {
	s.messageRegister.RegisterHandler(msgId, handler)
}

func (s *GameServer) EnqueueNotification(notification *message.Notification) {
	s.notificationQueue <- notification
}

func (s *GameServer) handleConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.ErrorWith("接受连接失败", "error", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection 处理单个连接
func (s *GameServer) handleConnection(conn net.Conn) {
	sessionId := "ps-" + uuid.New().String()
	s.sessionMgr.StartPlayerSession(sessionId, conn, s)
	defer s.sessionMgr.CloseSession(sessionId)
}

func (s *GameServer) EnqueueMessage(msg *message.Message) {
	s.messageQueue <- msg
}

func (s *GameServer) handleMessages() {
	for msg := range s.messageQueue {
		s.handleMessage(msg)
	}
}

// HandleMessage 实现Handler接口
func (s *GameServer) handleMessage(msg *message.Message) {
	sessionId := msg.SessionId
	clientMsg, err := serialize.Deserialize(msg.Data)
	if err != nil { // 这里反序列化只能够丢弃处理，因为没法给客户端返回正确的消息头，或者可以用通知？
		logger.ErrorWith("解析消息失败", "error", err)
		return
	}

	session, err := s.sessionMgr.GetSession(sessionId)
	if err != nil { // 这里获取不到session也无法处理，因为无法发送消息
		logger.ErrorWith("获取会话失败", "error", err)
		return
	}

	// 处理请求
	result, enableEncryption, err := s.handleRequest(session, clientMsg)

	respMessage := createResponseMsg(clientMsg.GetHeader())
	var respPayload proto.Message
	var notifyPayload proto.Message

	// 处理错误
	var gameError *errordef.GameError
	if gameError = errordef.AsGameError(err); gameError != nil {
		if gameError.Category == errordef.CategoryGameplay {
			logger.ErrorWith("游戏逻辑错误", "errorCode", gameError.Code, "errorMessage", gameError.ClientMsg)
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
		respPayloadBytes, iv, err = session.EncryptPayload(respPayloadBytes)
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
	msgId := respMessage.Header.MsgId
	if isSensitiveMessage(msgId) {
		logger.InfoWith("发送消息成功", "msgId", msgId, "sessionId", msg.SessionId)
	} else {
		logger.InfoWith("发送消息成功", "msgId", msgId, "sessionId", msg.SessionId, "payload", respPayload)
	}

	// 发送通知
	// 包装为 GameServerMessage
	if gameError == nil && notifyPayload != nil && result.NotifyGroup != nil {
		s.notify(notifyPayload, result.NotifyGroup, result.NotifyMsgId)
	}
}

func (s *GameServer) handleNotifications() {
	for notification := range s.notificationQueue {
		s.notify(notification.Payload, notification.NotifyGroup, notification.NotifyMsgId)
	}
}

func (s *GameServer) notify(notifyPayload proto.Message, notifyGroup message.INotifyGroup, notifyMsgId protodef.PMsgId) {
	if isSensitiveMessage(notifyMsgId) {
		logger.InfoWith("发送通知", "msgId", notifyMsgId)
	} else {
		logger.InfoWith("发送通知", "msgId", notifyMsgId, "payload", notifyPayload)
	}

	notifyPayloadBytes, err := serialize.SerializePayload(notifyPayload)
	if err != nil {
		logger.ErrorWith("序列化通知失败", "error", err)
		return
	}
	targetSessionIds := getNotifySessionIds(s.sessionMgr, notifyGroup)
	for _, targetSessionId := range targetSessionIds {
		go s.notifySession(targetSessionId, notifyPayloadBytes, notifyMsgId)
	}
}

func (s *GameServer) notifySession(targetSessionId string, notifyPayloadBytes []byte, notifyMsgId protodef.PMsgId) {
	targetSession, err := s.sessionMgr.GetSession(targetSessionId)
	if err != nil {
		logger.ErrorWith("获取会话失败", "error", err)
		return
	}
	var notificationPayloadBytes []byte
	var iv []byte
	if targetSession.IsSecureKeyValid() {
		notificationPayloadBytes, iv, err = targetSession.EncryptPayload(notifyPayloadBytes)
		if err != nil {
			logger.ErrorWith("加密通知失败", "error", err)
			return
		}
	} else {
		notificationPayloadBytes = notifyPayloadBytes
	}
	notificationMessage := createNotificationMsg(targetSession, notifyMsgId)
	notificationMessage.Payload = notificationPayloadBytes
	notificationMessage.Header.Iv = iv
	notificationData, err := serialize.Serialize(notificationMessage)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return
	}
	err = targetSession.SendMessage(notificationData)
	if err != nil {
		logger.ErrorWith("发送消息失败", "error", err)
		return
	}
}

func (s *GameServer) handleRequest(session *session.PlayerSession, clientMsg *protodef.PClientMsg) (result *message.HandleResult, enableEncryption bool, err error) {
	sessionId := session.Id
	// 消息体解密并反序列化
	enableEncryption = session.IsSecureKeyValid()
	msgHeader := clientMsg.GetHeader()
	msgId := msgHeader.GetMsgId()
	var reqPayloadBytes []byte
	if enableEncryption {
		reqPayloadBytes, err = session.DecryptPayload(clientMsg.GetPayload(), msgHeader.Iv)
		if err != nil {
			logger.ErrorWith("解密消息失败", "error", err)
			return nil, enableEncryption, err
		}
	} else {
		reqPayloadBytes = clientMsg.GetPayload()
	}
	reqPayload, err := serialize.DeserializePayload(msgId, reqPayloadBytes)
	if err != nil {
		logger.ErrorWith("反序列化消息失败", "error", err)
		return nil, enableEncryption, err
	}

	if isSensitiveMessage(msgId) {
		logger.InfoWith("收到消息", "类型", msgId, "sessionId", sessionId)
	} else {
		logger.InfoWith("收到消息", "类型", msgId, "sessionId", sessionId, "payload", reqPayload)
	}

	// 处理消息
	handler := s.messageRegister.GetHandler(msgId)
	if handler == nil {
		logger.ErrorWith("未找到消息处理器", "type", msgId)
		return nil, enableEncryption, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	context := &message.MessageContext{
		SessionId: sessionId,
		PlayerId:  msgHeader.PlayerId,
		Timestamp: msgHeader.Timestamp,

		Dispatcher: s,
	}
	result, err = handler(context, &reqPayload)
	return result, enableEncryption, err
}

func (s *GameServer) handleHandshake(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PHandshakeRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}

	session, err := s.sessionMgr.GetSession(context.SessionId)
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
			Timestamp: time.Now().UnixMilli(),

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
		},
		MsgType: protodef.PServerMsgType_PSERVER_MSG_TYPE_NOTIFICATION,
	}
}
