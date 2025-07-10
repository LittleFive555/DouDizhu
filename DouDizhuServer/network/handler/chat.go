package handler

import (
	"DouDizhuServer/errordef"
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/gameplay/room"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/translator"

	"google.golang.org/protobuf/proto"
)

func HandleChatMessage(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PChatMsgRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	chatMsg := reqMsg.GetContent()
	fromPlayer := player.Manager.GetPlayer(context.PlayerId)
	logger.InfoWith("收到聊天消息", "content", chatMsg)

	notification := &protodef.PChatMsgNotification{
		From:    translator.PlayerToProto(fromPlayer),
		Channel: reqMsg.GetChannel(),
		Content: chatMsg,
	}

	// 聊天消息通知组
	var notificationGroup message.INotificationGroup
	channel := reqMsg.GetChannel()
	switch channel {
	case protodef.PChatChannel_PCHAT_CHANNEL_ALL:
		notificationGroup = player.NewAllPlayerNotificationGroup()
	case protodef.PChatChannel_PCHAT_CHANNEL_ROOM:
		notificationGroup = room.NewRoomNotificationGroup(fromPlayer.GetRoomId())
	}

	result := &message.HandleResult{
		Resp:        nil,
		NotifyMsgId: protodef.PMsgId_PMSG_ID_CHAT_MSG,
		NotifyGroup: notificationGroup,
		Notify:      notification,
	}
	return result, nil
}
