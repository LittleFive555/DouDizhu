package message

import (
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

type MessageRegister struct {
	handlers map[protodef.PMsgId]func(*MessageContext, *proto.Message) (*HandleResult, error)
}

func NewMessageRegister() *MessageRegister {
	return &MessageRegister{
		handlers: make(map[protodef.PMsgId]func(*MessageContext, *proto.Message) (*HandleResult, error)),
	}
}

// RegisterHandler 注册消息处理器
func (md *MessageRegister) RegisterHandler(msgId protodef.PMsgId, handler func(*MessageContext, *proto.Message) (*HandleResult, error)) {
	md.handlers[msgId] = handler
}

func (md *MessageRegister) GetHandler(msgId protodef.PMsgId) func(*MessageContext, *proto.Message) (*HandleResult, error) {
	return md.handlers[msgId]
}
