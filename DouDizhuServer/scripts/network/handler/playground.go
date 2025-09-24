package handler

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

func HandleCharacterMove(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PCharacterMove)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	logger.InfoWith("收到角色移动消息", "move", reqMsg)
	// world := playground.Playground.World

	return &message.HandleResult{}, nil
}
