package handler

import (
	"DouDizhuServer/errordef"
	"DouDizhuServer/gameplay/player"
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
		Content: chatMsg,
	}
	result := &message.HandleResult{
		Resp:        nil,
		NofityMsgId: protodef.PMsgId_PMSG_ID_CHAT_MSG,
		NotifyGroup: player.NewAllPlayerNotificationGroup(), // TODO 后续要主动设置group
		Notify:      notification,
	}
	return result, nil
}
