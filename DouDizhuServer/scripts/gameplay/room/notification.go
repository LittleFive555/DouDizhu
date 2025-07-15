package room

import "DouDizhuServer/scripts/gameplay/player"

type RoomNotificationGroup struct {
	sessionIds []string
}

func NewRoomNotificationGroup(roomId uint32) *RoomNotificationGroup {
	room, err := Manager.GetRoom(roomId)
	if err != nil {
		return nil
	}
	playerIds := room.GetPlayers()
	sessionIds := make([]string, 0)
	for _, playerId := range playerIds {
		sessionIds = append(sessionIds, player.Manager.GetPlayer(playerId).GetSessionId())
	}

	return &RoomNotificationGroup{
		sessionIds: sessionIds,
	}
}

func NewRoomNotificationGroupExcept(roomId uint32, exceptPlayerId string) *RoomNotificationGroup {
	room, err := Manager.GetRoom(roomId)
	if err != nil {
		return nil
	}
	playerIds := room.GetPlayers()
	sessionIds := make([]string, 0)
	for _, playerId := range playerIds {
		if playerId != exceptPlayerId {
			sessionIds = append(sessionIds, player.Manager.GetPlayer(playerId).GetSessionId())
		}
	}
	return &RoomNotificationGroup{
		sessionIds: sessionIds,
	}
}

func (g *RoomNotificationGroup) GetTargetSessionIds() []string {
	return g.sessionIds
}
