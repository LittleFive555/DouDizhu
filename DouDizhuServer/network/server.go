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

func (s *GameServer) RegisterHandler(msgType reflect.Type, handler func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)) {
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
	logger.InfoWith("收到消息", "sessionId", msg.SessionId, "message", reqPacket)

	// 查找对应的消息处理器
	content := reqPacket.GetContent()
	if content == nil {
		logger.ErrorWith("消息内容为空")
		return
	}
	msgType := reflect.TypeOf(content).Elem()
	handler := s.dispatcher.GetHandler(msgType)
	if handler == nil {
		logger.ErrorWith("未找到消息处理器", "type", msgType.String())
		return
	}

	// 处理消息
	respPacket, err := handler(reqPacket)

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
	// 包装为 GameServerMessage
	serverMessage := &protodef.PGameServerMessage{
		Content: &protodef.PGameServerMessage_Response{
			Response: respPacket,
		},
	}

	responseData, err := serialize.Serialize(serverMessage)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return
	}
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

	logger.InfoWith("发送消息成功", "message", serverMessage)
}
