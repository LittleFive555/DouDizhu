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
	err := room.SetOwner(ownerPlayer.GetPlayerId())
	if err != nil {
		return nil, err
	}
	ownerPlayer.EnterRoom(room.GetId())
	return &message.HandleResult{
		Resp: &protodef.PCreateRoomResponse{
			Room: translator.RoomToProto(room, player.Manager),
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
		rooms = append(rooms, translator.RoomToProto(r, player.Manager))
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
	err = targetRoom.AddPlayer(requestingPlayer.GetPlayerId())
	if err != nil {
		return nil, err
	}
	requestingPlayer.EnterRoom(targetRoom.GetId())

	return &message.HandleResult{
		Resp: &protodef.PEnterRoomResponse{
			Room: translator.RoomToProto(targetRoom, player.Manager),
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
	currentPlayer := player.Manager.GetPlayer(context.PlayerId)
	if currentPlayer == nil {
		return nil, errordef.NewGameplayError(errordef.CodePlayerOffline)
	}
	playerId := currentPlayer.GetPlayerId()
	if !currentPlayer.IsInRoom() {
		return nil, errordef.NewGameplayError(errordef.CodePlayerNotInRoom)
	}
	currentRoom, err := room.Manager.GetRoom(currentPlayer.GetRoomId())
	if err != nil {
		return nil, err
	}

	notifyGroup := &room.RoomNotifyGroup{
		RoomId:         currentRoom.GetId(),
		ExceptPlayerId: playerId,
	}
	if currentRoom.IsOwnedBy(playerId) { // 如果房间是房主，则直接解散整个房间
		players := currentRoom.GetPlayers()
		for _, playerId := range players {
			player.Manager.GetPlayer(playerId).LeaveRoom()
		}
		room.Manager.RemoveRoom(currentRoom.GetId())

		return &message.HandleResult{
			Notify:      &protodef.PRoomDisbandedNotification{},
			NotifyMsgId: protodef.PMsgId_PMSG_ID_ROOM_DISBANDED,
			NotifyGroup: notifyGroup,
		}, nil
	} else {
		currentRoom.RemovePlayer(playerId)
		currentPlayer.LeaveRoom()
		// 如果房间不是房主，则给房间其他人发通知
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

func tryGetPlayerRoom(player *player.Player) (*room.Room, error) {
	if player == nil {
		return nil, errordef.NewGameplayError(errordef.CodePlayerOffline)
	}
	if !player.IsInRoom() {
		return nil, errordef.NewGameplayError(errordef.CodePlayerNotInRoom)
	}
	room, err := room.Manager.GetRoom(player.GetRoomId())
	if err != nil {
		return nil, err
	}
	return room, nil
}
