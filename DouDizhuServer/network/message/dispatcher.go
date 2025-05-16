package message

import (
	"DouDizhuServer/network/protodef"
	"reflect"
)

var Dispatcher *MessageDispatcher

type MessageDispatcher struct {
	messageQueue chan *Message // 消息队列
	workers      []*Worker     // 工作协程池
	handlers     map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)
}

// 初始化消息分发器
func NewMessageDispatcher(workerCount int) *MessageDispatcher {
	md := &MessageDispatcher{
		messageQueue: make(chan *Message, 10000), // 带缓冲的队列
		handlers:     make(map[reflect.Type]func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)),
	}

	// 初始化工作协程
	for i := 0; i < workerCount; i++ {
		worker := NewWorker(md.messageQueue, md.handlers)
		md.workers = append(md.workers, worker)
		go worker.Run()
	}

	return md
}

// RegisterHandler 注册消息处理器
func (md *MessageDispatcher) RegisterHandler(msgType reflect.Type, handler func(*protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error)) {
	md.handlers[msgType] = handler
}

// 接收消息
func (md *MessageDispatcher) PostMessage(msg *Message) {
	md.messageQueue <- msg
}
