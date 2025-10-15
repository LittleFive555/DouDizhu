package handler

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/gameplay/player"
	"DouDizhuServer/scripts/gameplay/room"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"
	"DouDizhuServer/scripts/network/translator"

	"google.golang.org/protobuf/proto"
)

func HandleCreateRoom(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PCreateRoomRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	ownerPlayer := player.Manager.GetPlayer(context.PlayerId)
	if ownerPlayer == nil {
		return nil, errordef.NewGameplayError(errordef.CodePlayerOffline)
	}
	if !ownerPlayer.IsInLobby() {
		return nil, errordef.NewGameplayError(errordef.CodePlayerNotInLobby)
	}
	room := room.Manager.CreateRoom(reqMsg.GetRoomName(), context.Dispatcher)
	err := room.AddPlayer(ownerPlayer, true)
	if err != nil {
		return nil, err
	}
	return &message.HandleResult{
		Resp: &protodef.PCreateRoomResponse{
			Room: translator.RoomToProto(room, player.Manager, true),
		},
	}, nil
}

func HandleGetRoomList(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	_, ok := (*req).(*protodef.PGetRoomListRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	roomList := room.Manager.GetRoomList()
	rooms := make([]*protodef.PRoom, 0)
	for _, r := range roomList {
		rooms = append(rooms, translator.RoomToProto(r, player.Manager, false))
	}
	return &message.HandleResult{
		Resp: &protodef.PGetRoomListResponse{
			Rooms: rooms,
		},
	}, nil
}

func HandleEnterRoom(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PEnterRoomRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	requestingPlayer := player.Manager.GetPlayer(context.PlayerId)
	if requestingPlayer == nil {
		return nil, errordef.NewGameplayError(errordef.CodePlayerOffline)
	}
	if !requestingPlayer.IsInLobby() {
		return nil, errordef.NewGameplayError(errordef.CodePlayerNotInLobby)
	}
	targetRoom, err := room.Manager.GetRoom(reqMsg.GetRoomId())
	if err != nil {
		return nil, err
	}
	err = targetRoom.AddPlayer(requestingPlayer, false)
	if err != nil {
		return nil, err
	}

	return &message.HandleResult{
		Resp: &protodef.PEnterRoomResponse{
			Room: translator.RoomToProto(targetRoom, player.Manager, true),
		},
		NotifyMsgId: protodef.PMsgId_PMSG_ID_ROOM_CHANGED,
		Notify: &protodef.PRoomChangedNotification{
			Room: &protodef.PRoom{
				Players: translator.RoomPlayersToProto(targetRoom, player.Manager),
			},
		},
		NotifyGroup: &room.RoomNotifyGroup{
			RoomId:         targetRoom.GetId(),
			ExceptPlayerId: requestingPlayer.GetPlayerId(),
		},
	}, nil
}

func HandleLeaveRoom(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	_, ok := (*req).(*protodef.PLeaveRoomRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	playerId := context.PlayerId

	currentPlayer := player.Manager.GetPlayer(playerId)
	if currentPlayer == nil {
		return nil, errordef.NewGameplayError(errordef.CodePlayerOffline)
	}
	if !currentPlayer.IsInRoom() {
		return nil, errordef.NewGameplayError(errordef.CodePlayerNotInRoom)
	}
	currentRoom, err := tryGetPlayerRoom(playerId)
	if err != nil {
		return nil, err
	}

	notifyGroup := &room.RoomNotifyGroup{
		RoomId:         currentRoom.GetId(),
		ExceptPlayerId: playerId,
	}
	if currentRoom.IsOwnedBy(playerId) { // 如果房间是房主，则直接解散整个房间
		room.Manager.RemoveRoom(currentRoom.GetId())

		return &message.HandleResult{
			Notify:      &protodef.PRoomDisbandedNotification{},
			NotifyMsgId: protodef.PMsgId_PMSG_ID_ROOM_DISBANDED,
			NotifyGroup: notifyGroup,
		}, nil
	} else {
		currentRoom.RemovePlayer(playerId)
		// 如果房间不是房主，则给房间其他人发通知
		// TODO 需要把离开世界的玩家状态也通知到其他客户端
		return &message.HandleResult{
			Notify: &protodef.PRoomChangedNotification{
				Room: &protodef.PRoom{
					Players: translator.RoomPlayersToProto(currentRoom, player.Manager),
				},
			},
			NotifyMsgId: protodef.PMsgId_PMSG_ID_ROOM_CHANGED,
			NotifyGroup: notifyGroup,
		}, nil
	}
}

func tryGetPlayerRoom(playerId string) (*room.Room, error) {
	player := player.Manager.GetPlayer(playerId)
	if player == nil {
		return nil, errordef.NewGameplayError(errordef.CodePlayerOffline)
	}
	if !player.IsInRoom() {
		return nil, errordef.NewGameplayError(errordef.CodePlayerNotInRoom)
	}
	return room.Manager.GetRoom(player.GetRoomId())
}
