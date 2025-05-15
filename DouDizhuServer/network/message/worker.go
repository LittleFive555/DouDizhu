package message

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/serialize"
	"errors"
	"reflect"
)

type Worker struct {
	queue    <-chan *Message
	handlers map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)
}

func NewWorker(queue <-chan *Message, handlers map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)) *Worker {
	return &Worker{
		queue:    queue,
		handlers: handlers,
	}
}

func (w *Worker) Run() {
	for msg := range w.queue {
		// 获取消息的实际类型
		w.handleMessage(msg.SessionId, msg.Data)
	}
}

// HandleMessage 实现Handler接口
func (s *Worker) handleMessage(sessionId string, data []byte) ([]byte, error) {
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
