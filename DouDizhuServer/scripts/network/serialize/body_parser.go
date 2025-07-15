package serialize

import (
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

var messageRegistry = map[protodef.PMsgId]func() proto.Message{
	protodef.PMsgId_PMSG_ID_REGISTER:  func() proto.Message { return &protodef.PRegisterRequest{} },
	protodef.PMsgId_PMSG_ID_LOGIN:     func() proto.Message { return &protodef.PLoginRequest{} },
	protodef.PMsgId_PMSG_ID_HANDSHAKE: func() proto.Message { return &protodef.PHandshakeRequest{} },

	// Room
	protodef.PMsgId_PMSG_ID_CREATE_ROOM:   func() proto.Message { return &protodef.PCreateRoomRequest{} },
	protodef.PMsgId_PMSG_ID_GET_ROOM_LIST: func() proto.Message { return &protodef.PGetRoomListRequest{} },
	protodef.PMsgId_PMSG_ID_ENTER_ROOM:    func() proto.Message { return &protodef.PEnterRoomRequest{} },
	protodef.PMsgId_PMSG_ID_LEAVE_ROOM:    func() proto.Message { return &protodef.PLeaveRoomRequest{} },

	// Chat
	protodef.PMsgId_PMSG_ID_CHAT_MSG: func() proto.Message { return &protodef.PChatMsgRequest{} },
}

func GetMessage(msgId protodef.PMsgId) proto.Message {
	return messageRegistry[msgId]()
}
