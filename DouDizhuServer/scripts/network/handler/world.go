package handler

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/gameplay/playground"
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

func HandleEnterWorld(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	_, ok := (*req).(*protodef.PEnterWorldRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	if playground.Playground != nil { // 在playground模式下，直接进入playground世界
		characterId := playground.Playground.World.AddCharacter()
		worldState := playground.Playground.World.GetFullWrldState()
		return &message.HandleResult{
			Resp: &protodef.PEnterWorldResponse{
				CharacterId: characterId,
				WorldState:  worldState,
			},
		}, nil
	} else { // 在游戏模式下
		// TODO
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
}

func HandleLeaveWorld(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PLeaveWorldRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	if playground.Playground != nil { // 在playground模式下，直接离开playground世界
		playground.Playground.World.RemoveCharacter(reqMsg.CharacterId)
	} else { // 在游戏模式下

	}
	return &message.HandleResult{}, nil
}

func HandleCharacterMove(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PCharacterMove)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	logger.InfoWith("收到角色移动消息", "move", reqMsg)
	if playground.Playground != nil { // 在playground模式下，直接移动playground世界
		playground.Playground.World.MoveCharacter(reqMsg)
	} else { // 在游戏模式下
		// TODO
	}

	return &message.HandleResult{}, nil
}
