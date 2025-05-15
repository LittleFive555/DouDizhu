package session

type PlayerSession struct {
	PlayerId  string
	SessionId string
}

func NewPlayerSession(playerId string, sessionId string) *PlayerSession {
	return &PlayerSession{
		PlayerId:  playerId,
		SessionId: sessionId,
	}
}
