package tcp

import (
	"net"
)

type NetworkSession struct {
	SessionId string
	Conn      net.Conn
}

func NewNetworkSession(sessionId string, conn net.Conn) *NetworkSession {
	return &NetworkSession{
		SessionId: sessionId,
		Conn:      conn,
	}
}
