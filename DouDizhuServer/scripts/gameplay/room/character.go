package room

import "DouDizhuServer/scripts/network/protodef"

type RoomCharacter struct {
	id       string
	position *protodef.PVector3
}
