package network

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/protodef"
	"DouDizhuServer/network/serialize"
	"DouDizhuServer/network/session"
	"DouDizhuServer/network/tcp"
)

var Server *GameServer

// GetServer 返回游戏服务器实例
func GetServer() *GameServer {
	return Server
}

type GameServer struct {
	server  *tcp.TCPServer
	session *session.SessionManager
}

func NewGameServer(addr string) *GameServer {
	server := tcp.NewTCPServer(addr)
	gameServer := &GameServer{
		server:  server,
		session: session.NewSessionManager(),
	}
	return gameServer
}

func (s *GameServer) Start() error {
	return s.server.Start()
}

func (s *GameServer) Stop() error {
	return s.server.Stop()
}

func (s *GameServer) SendNotificationToPlayer(playerId string, notification *protodef.PGameNotificationPacket) {
	playerSession := s.session.GetPlayerSession(playerId)
	if playerSession == nil {
		logger.ErrorWith("玩家不存在", "playerId", playerId)
		return
	}

	notificationMessage := &protodef.PGameServerMessage{
		Content: &protodef.PGameServerMessage_Notification{
			Notification: notification,
		},
	}
	logger.InfoWith("发送通知", "notification", notificationMessage)
	notificationData, err := serialize.Serialize(notificationMessage)
	if err != nil {
		logger.ErrorWith("序列化通知失败", "error", err)
		return
	}

	s.server.SendMessage(playerSession.SessionId, notificationData)
}
