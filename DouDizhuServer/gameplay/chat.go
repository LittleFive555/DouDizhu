package gameplay

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network"
	"DouDizhuServer/network/protodef"
)

func HandleChatMessage(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	logger.InfoWith("收到聊天消息", "content", req.GetChatMsg().GetContent())

	notification := &protodef.PGameNotificationPacket{
		Content: &protodef.PGameNotificationPacket_ChatMsg{
			ChatMsg: &protodef.PChatMsgNotification{
				From:    req.Header.Player,
				Content: req.GetChatMsg().GetContent(),
			},
		},
	}
	network.GetServer().SendNotification(notification)
	return &protodef.PGameMsgRespPacket{
		Header: &protodef.PGameMsgHeader{},
		Content: &protodef.PGameMsgRespPacket_CommonResponse{
			CommonResponse: &protodef.PCommonResponse{},
		},
	}, nil
}
