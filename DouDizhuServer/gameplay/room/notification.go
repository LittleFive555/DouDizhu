package room

import "DouDizhuServer/gameplay/player"

type RoomNotificationGroup struct {
	roomId uint32
}

func NewRoomNotificationGroup(roomId uint32) *RoomNotificationGroup {
	return &RoomNotificationGroup{
		roomId: roomId,
	}
}

func (g *RoomNotificationGroup) GetTargetSessionIds() []string {
	room, err := Manager.GetRoom(g.roomId)
	if err != nil {
		return nil
	}
	playerIds := room.GetPlayers()
	sessionIds := make([]string, 0)
	for _, playerId := range playerIds {
		sessionIds = append(sessionIds, player.Manager.GetPlayer(playerId).GetSessionId())
	}
	return sessionIds
}
