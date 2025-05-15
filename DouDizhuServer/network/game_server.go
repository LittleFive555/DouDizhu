package network

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/serialize"
	"DouDizhuServer/network/tcp"
	"errors"
	"reflect"
)

var Server *GameServer

// GetServer 返回游戏服务器实例
func GetServer() *GameServer {
	return Server
}

type GameServer struct {
	server   *tcp.TCPServer
	handlers map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)
}

func NewGameServer(addr string) *GameServer {
	server := tcp.NewTCPServer(addr, tcp.NewLengthPrefixConnIO())
	gameServer := &GameServer{
		server:   server,
		handlers: make(map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)),
	}
	server.SetMessageConsumer(gameServer.handleMessage)
	return gameServer
}

func (s *GameServer) Start() error {
	return s.server.Start()
}

func (s *GameServer) Stop() error {
	return s.server.Stop()
}

// RegisterHandler 注册消息处理器
func (s *GameServer) RegisterHandler(msgType reflect.Type, handler func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)) {
	s.handlers[msgType] = handler
}

func (s *GameServer) SendNotification(notification *protodef.PGameNotificationPacket) error {
	notificationMessage := &protodef.PGameServerMessage{
		Content: &protodef.PGameServerMessage_Notification{
			Notification: notification,
		},
	}
	logger.InfoWith("发送通知", "notification", notificationMessage)
	notificationData, err := serialize.Serialize(notificationMessage)
	if err != nil {
		logger.ErrorWith("序列化通知失败", "error", err)
		return err
	}
	return s.server.NotifyAll(notificationData)
}

// HandleMessage 实现Handler接口
func (s *GameServer) handleMessage(data []byte) ([]byte, error) {
	reqPacket, err := serialize.Deserialize(data)
	if err != nil {
		logger.ErrorWith("解析消息失败", "error", err)
		return nil, err
	}
	// 查找对应的消息处理器
	content := reqPacket.GetContent()
	if content == nil {
		logger.ErrorWith("消息内容为空")
		return nil, errors.New("消息内容为空")
	}

	// 获取消息的实际类型
	msgType := reflect.TypeOf(content).Elem()
	handler, exists := s.handlers[msgType]
	if !exists {
		logger.ErrorWith("未找到消息处理器", "type", msgType.String())
		return nil, errors.New("未找到消息处理器")
	}

	// 处理消息
	respPacket, err := handler(reqPacket)
	if err != nil {
		logger.ErrorWith("处理消息失败", "error", err)
		return nil, err
	}

	// 包装为 GameServerMessage
	respPacket.Header.MessageId = reqPacket.Header.MessageId
	serverMessage := &protodef.PGameServerMessage{
		Content: &protodef.PGameServerMessage_Response{
			Response: respPacket,
		},
	}

	logger.InfoWith("处理消息成功", "message", serverMessage)
	responseData, err := serialize.Serialize(serverMessage)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return nil, err
	}

	return responseData, nil
}
