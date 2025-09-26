package network

import (
	"DouDizhuServer/scripts/network/handler"
	"DouDizhuServer/scripts/network/protodef"
)

func (s *GameServer) RegisterPlaygroundHandlers() {
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_ENTER_WORLD, handler.HandleEnterWorld)
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_LEAVE_WORLD, handler.HandleLeaveWorld)
	s.RegisterHandler(protodef.PMsgId_PMSG_ID_CHARACTER_MOVE, handler.HandleCharacterMove)
}
