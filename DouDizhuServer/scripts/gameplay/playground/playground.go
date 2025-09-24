package playground

import "DouDizhuServer/scripts/gameplay/room"

var Playground *RoomPlayground

type RoomPlayground struct {
	World *room.RoomWorld
}

func NewRoomPlayground() *RoomPlayground {
	return &RoomPlayground{
		World: room.NewRoomWorld("playground-111999555", 10),
	}
}

func (p *RoomPlayground) Start() {
	go p.World.RunLoop()
}
