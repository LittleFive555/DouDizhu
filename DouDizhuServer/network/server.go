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
	"reflect"
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

func (s *GameServer) RegisterHandler(msgType reflect.Type, handler func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, *protodef.PGameNotificationPacket, error)) {
	s.dispatcher.RegisterHandler(msgType, handler)
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
	reqPacket, err := serialize.Deserialize(msg.Data)
	if err != nil {
		logger.ErrorWith("解析消息失败", "error", err)
		return
	}

	// 查找对应的消息处理器
	content := reqPacket.GetContent()
	if content == nil {
		logger.ErrorWith("消息内容为空")
		return
	}
	msgType := reflect.TypeOf(content).Elem()

	if isSecretMessage(reqPacket) {
		logger.InfoWith("收到消息", "类型", msgType.String(), "sessionId", msg.SessionId)
	} else {
		logger.InfoWith("收到消息", "类型", msgType.String(), "sessionId", msg.SessionId, "message", reqPacket)
	}

	handler := s.dispatcher.GetHandler(msgType)
	if handler == nil {
		logger.ErrorWith("未找到消息处理器", "type", msgType.String())
		return
	}

	// 处理消息
	respPacket, notificationPacket, err := handler(reqPacket)

	if gameError := errors.AsGameError(err); gameError != nil {
		if gameError.Category == errors.CategoryGameplay {
			logger.ErrorWith("游戏逻辑错误", "errorCode", gameError.Code, "errorMessage", gameError.ClientMsg)
		} else {
			logger.ErrorWith("服务器错误", "errorCategory", gameError.Category, "errorCode", gameError.Code, "errorMessage", gameError.ClientMsg)
		}
		respPacket = message.CreateErrorPacket(reqPacket.Header, gameError)
	} else {
		respPacket.Header.MessageId = reqPacket.Header.MessageId
	}

	// 响应
	// 包装为 GameServerMessage
	serverRespMessage := message.CreateServerMessage()
	serverRespMessage.Content = &protodef.PGameServerMessage_Response{
		Response: respPacket,
	}
	responseData, err := serialize.Serialize(serverRespMessage)
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
	logger.InfoWith("发送消息成功", "message", serverRespMessage)

	// 发送通知
	// 包装为 GameServerMessage
	if notificationPacket != nil {
		serverNotificationMessage := message.CreateServerMessage()
		serverNotificationMessage.Content = &protodef.PGameServerMessage_Notification{
			Notification: notificationPacket,
		}
		notificationData, err := serialize.Serialize(serverNotificationMessage)
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
		logger.InfoWith("发送通知结束", "message", serverNotificationMessage)
	}
}

func isSecretMessage(msg *protodef.PGameClientMessage) bool {
	if msg.Content == nil {
		return false
	}
	contentType := reflect.TypeOf(msg.Content)
	if contentType == reflect.TypeOf(&protodef.PGameClientMessage_RegisterReq{}) ||
		contentType == reflect.TypeOf(&protodef.PGameClientMessage_LoginReq{}) {
		return true
	}
	return false
}
