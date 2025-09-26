package room

type RoomNotifyGroup struct {
	RoomId         uint32
	ExceptPlayerId string
}

func (g *RoomNotifyGroup) IsNotifyGroup() {
}
