package handler

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/gameplay/player"
	"DouDizhuServer/scripts/gameplay/room"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

// 进入或切换世界
func HandleEnterWorld(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	enterWorldReq, ok := (*req).(*protodef.PEnterWorldRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	// 先查找玩家所处的房间
	player := player.Manager.GetPlayer(context.PlayerId)
	if player == nil {
		return nil, errordef.NewGameplayError(errordef.CodePlayerOffline)
	}
	roomId := player.GetRoomId()
	if roomId == 0 {
		return nil, errordef.NewGameplayError(errordef.CodePlayerNotInRoom)
	}
	room, err := room.Manager.GetRoom(roomId)
	if err != nil {
		return nil, errordef.NewGameplayError(errordef.CodeRoomNotExists)
	}
	characterId, err := room.GetCharacterId(context.PlayerId)
	if err != nil {
		return nil, err
	}
	newWorldId := enterWorldReq.WorldId
	worldManager := room.GetWorldManager()
	worldManager.CharacterEnterWorld(characterId, newWorldId)
	return &message.HandleResult{
		Resp: &protodef.PEnterWorldResponse{
			WorldState: worldManager.GetWorldFullState(newWorldId),
		},
	}, nil
}

func HandleCharacterMove(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PCharacterMove)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	room, err := tryGetPlayerRoom(context.PlayerId)
	if err != nil {
		return nil, err
	}
	world := room.GetWorldManager()
	world.CharacterInput(reqMsg.CharacterId, reqMsg)

	return &message.HandleResult{}, nil
}
