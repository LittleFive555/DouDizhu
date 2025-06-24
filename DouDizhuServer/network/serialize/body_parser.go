package serialize

import (
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

var messageRegistry = map[protodef.PMsgId]func() proto.Message{
	protodef.PMsgId_PMSG_ID_REGISTER:  func() proto.Message { return &protodef.PRegisterRequest{} },
	protodef.PMsgId_PMSG_ID_LOGIN:     func() proto.Message { return &protodef.PLoginRequest{} },
	protodef.PMsgId_PMSG_ID_HANDSHAKE: func() proto.Message { return &protodef.PHandshakeRequest{} },
	protodef.PMsgId_PMSG_ID_CHAT_MSG:  func() proto.Message { return &protodef.PChatMsgRequest{} },
}

func GetMessage(msgId protodef.PMsgId) proto.Message {
	return messageRegistry[msgId]()
}
