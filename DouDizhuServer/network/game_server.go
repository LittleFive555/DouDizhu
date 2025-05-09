package network

import (
	"DouDizhuServer/gameplay"
	"DouDizhuServer/network/handler"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/tcp"
	"reflect"
)

type GameServer struct {
	server *tcp.TCPServer
}

func NewGameServer(addr string) *GameServer {
	messageHandler := handler.NewProtoHandler()
	messageHandler.RegisterHandler(reflect.TypeOf(protodef.GameMsgReqPacket_ChatMsg{}), gameplay.HandleChatMessage)
	return &GameServer{
		server: tcp.NewTCPServer(addr, messageHandler, tcp.NewLengthPrefixConnIO()),
	}
}

func (s *GameServer) Start() error {
	return s.server.Start()
}

func (s *GameServer) Stop() error {
	return s.server.Stop()
}
