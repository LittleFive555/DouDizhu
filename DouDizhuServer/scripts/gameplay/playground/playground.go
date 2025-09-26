package playground

import (
	"DouDizhuServer/scripts/gameplay/room"
	"DouDizhuServer/scripts/network/message"
)

var Playground *RoomPlayground

type RoomPlayground struct {
	World *room.RoomWorld
}

func NewRoomPlayground() *RoomPlayground {
	return &RoomPlayground{
		World: room.NewRoomWorld("playground-111999555", 10),
	}
}

func (p *RoomPlayground) Start(dispatcher message.INotificationDispatcher) {
	go p.World.RunLoop(dispatcher)
}
