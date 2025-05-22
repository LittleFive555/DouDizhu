package chat

import (
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
)

func HandleChatMessage(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, *protodef.PGameNotificationPacket, error) {
	chatMsg := req.GetChatMsg().GetContent()
	player := player.Manager.GetPlayer(req.Header.PlayerId)
	logger.InfoWith("收到聊天消息", "content", chatMsg)

	notification := message.CreateNotificationPacket(req.Header)
	notification.Content = &protodef.PGameNotificationPacket_ChatMsg{
		ChatMsg: &protodef.PChatMsgNotification{
			From:    player.ToProto(),
			Content: chatMsg,
		},
	}
	return message.CreateEmptyRespPacket(req.Header), notification, nil
}
