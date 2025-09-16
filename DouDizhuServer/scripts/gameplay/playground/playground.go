package playground

import "DouDizhuServer/scripts/gameplay/room"

var Playground *RoomPlayground

type RoomPlayground struct {
	World *room.RoomWorld
}

func NewRoomPlayground() *RoomPlayground {
	return &RoomPlayground{
		World: room.NewRoomWorld(),
	}
}
