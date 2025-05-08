package network

import (
	"DouDizhuServer/network/gameplay"
	"DouDizhuServer/network/handler"
	"DouDizhuServer/network/protodef"
	"reflect"
)

type GameServer struct {
	server *TCPServer
}

func NewGameServer(addr string) *GameServer {
	messageHandler := handler.NewProtoHandler()
	messageHandler.RegisterHandler(reflect.TypeOf(protodef.GameMsgReqPacket_ChatMsg{}), gameplay.HandleChatMessage)
	return &GameServer{
		server: NewTCPServer(addr, messageHandler),
	}
}

func (s *GameServer) Start() error {
	return s.server.Start()
}

func (s *GameServer) Stop() error {
	return s.server.Stop()
}
