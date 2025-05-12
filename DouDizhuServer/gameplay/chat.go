package gameplay

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network"
	"DouDizhuServer/network/protodef"
)

func HandleChatMessage(req *protodef.GameClientMessage) (*protodef.GameMsgRespPacket, error) {
	logger.InfoWith("收到聊天消息", "content", req.GetChatMsg().GetContent())

	notification := &protodef.GameNotificationPacket{
		Content: &protodef.GameNotificationPacket_ChatMsg{
			ChatMsg: &protodef.ChatMsgNotification{
				From:    req.Header.Player,
				Content: req.GetChatMsg().GetContent(),
			},
		},
	}
	network.GetServer().SendNotification(notification)
	return &protodef.GameMsgRespPacket{
		Header: &protodef.GameMsgHeader{},
		Content: &protodef.GameMsgRespPacket_CommonResponse{
			CommonResponse: &protodef.CommonResponse{},
		},
	}, nil
}
