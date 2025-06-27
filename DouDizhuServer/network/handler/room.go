package handler

import (
	"DouDizhuServer/errordef"
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/gameplay/room"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/translator"

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
	room := room.Manager.CreateRoom(reqMsg.GetRoomName())
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
		Notify: &protodef.PRoomChangedNotification{
			Room: &protodef.PRoom{
				Players: translator.RoomPlayersToProto(targetRoom, player.Manager),
			},
		},
		NotifyGroup: room.NewRoomNotificationGroup(targetRoom.GetId()),
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
	err = currentRoom.RemovePlayer(playerId)
	if err != nil {
		return nil, err
	}
	currentPlayer.LeaveRoom()

	var notify proto.Message
	if currentRoom.IsOwnedBy(playerId) {
		// 如果房间是房主，则直接解散整个房间
		room.Manager.RemoveRoom(currentRoom.GetId())

		notify = &protodef.PRoomDisbandedNotification{}
	} else {
		// 如果房间不是房主，则给房间其他人发通知
		notify = &protodef.PRoomChangedNotification{
			Room: &protodef.PRoom{
				Players: translator.RoomPlayersToProto(currentRoom, player.Manager),
			},
		}
	}
	return &message.HandleResult{
		Notify:      notify,
		NotifyGroup: room.NewRoomNotificationGroup(currentRoom.GetId()),
	}, nil
}
