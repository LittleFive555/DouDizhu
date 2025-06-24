package chat

import (
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
	"errors"

	"google.golang.org/protobuf/proto"
)

func HandleChatMessage(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PChatMsgRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	chatMsg := reqMsg.GetContent()
	player := player.Manager.GetPlayer(context.PlayerId)
	logger.InfoWith("收到聊天消息", "content", chatMsg)

	notification := &protodef.PChatMsgNotification{
		From:    player.ToProto(),
		Content: chatMsg,
	}
	result := &message.HandleResult{
		Resp:   nil,
		Notify: notification,
	}
	return result, nil
}
