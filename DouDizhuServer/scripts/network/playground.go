package network

import (
	"DouDizhuServer/scripts/network/handler"
	"DouDizhuServer/scripts/network/protodef"
)

func (s *GameServer) RegisterPlaygroundHandlers() {
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_CHARACTER_MOVE, handler.HandleCharacterMove)
}
