package room

type RoomNotificationGroup struct {
	roomId uint32
}

func NewRoomNotificationGroup(roomId uint32) *RoomNotificationGroup {
	return &RoomNotificationGroup{
		roomId: roomId,
	}
}

func (g *RoomNotificationGroup) GetTargetPlayerIds() []string {
	room, err := Manager.GetRoom(g.roomId)
	if err != nil {
		return nil
	}
	return room.GetPlayers()
}
