package message

import (
	"DouDizhuServer/network/protodef"
	"reflect"
)

type MessageDispatcher struct {
	messageQueue chan *Message // 消息队列
	workers      []*Worker     // 工作协程池
	handlers     map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)
	handler      func(message *Message)
}

// 初始化消息分发器
func NewMessageDispatcher(workerCount int, handler func(message *Message)) *MessageDispatcher {
	md := &MessageDispatcher{
		messageQueue: make(chan *Message, 10000), // 带缓冲的队列
		handlers:     make(map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)),
		handler:      handler,
	}

	// 初始化工作协程
	for i := 0; i < workerCount; i++ {
		worker := NewWorker(md.messageQueue, md.handler)
		md.workers = append(md.workers, worker)
		go worker.Run()
	}

	return md
}

// RegisterHandler 注册消息处理器
func (md *MessageDispatcher) RegisterHandler(msgType reflect.Type, handler func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)) {
	md.handlers[msgType] = handler
}

func (md *MessageDispatcher) GetHandler(msgType reflect.Type) func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	return md.handlers[msgType]
}

// 接收消息
func (md *MessageDispatcher) EnqueueMessage(msg *Message) {
	md.messageQueue <- msg
}
