package session

type SessionManager struct {
	playerSessions map[string]*PlayerSession
}

func NewSessionManager() *SessionManager {
	return &SessionManager{playerSessions: make(map[string]*PlayerSession)}
}

func (s *SessionManager) AddPlayerSession(session *PlayerSession) {
	s.playerSessions[session.PlayerId] = session
}

func (s *SessionManager) GetPlayerSession(playerId string) *PlayerSession {
	return s.playerSessions[playerId]
}
