package chat

import (
	"DouDizhuServer/errordef"
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

func HandleChatMessage(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PChatMsgRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	chatMsg := reqMsg.GetContent()
	player := player.Manager.GetPlayer(context.PlayerId)
	logger.InfoWith("收到聊天消息", "content", chatMsg)

	notification := &protodef.PChatMsgNotification{
		From:    player.ToProto(),
		Content: chatMsg,
	}
	result := &message.HandleResult{
		Resp:        nil,
		NofityMsgId: protodef.PMsgId_PMSG_ID_CHAT_MSG,
		Notify:      notification,
	}
	return result, nil
}
