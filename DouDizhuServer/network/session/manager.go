package session

import (
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
)

type SessionManager struct {
	playerSessions map[string]*PlayerSession
	mutex          sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{playerSessions: make(map[string]*PlayerSession)}
}

func (s *SessionManager) CreatePlayerSession(conn net.Conn) (*PlayerSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sessionId := "ps-" + uuid.New().String()
	ip := conn.RemoteAddr().String()
	session := &PlayerSession{
		Id:       sessionId,
		Conn:     conn,
		IP:       ip,
		State:    PlayerState_Connecting,
		PlayerId: "",
	}

	s.playerSessions[sessionId] = session
	return session, nil
}

func (s *SessionManager) Authenticate(sessionId, playerId string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.playerSessions[sessionId]
	if !exists {
		return fmt.Errorf("sessionId不存在")
	}

	session.PlayerId = playerId
	session.State = PlayerState_Lobby

	return nil
}

func (s *SessionManager) GetSession(sessionId string) (*PlayerSession, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.playerSessions[sessionId]
	if !exists {
		return nil, fmt.Errorf("sessionId不存在")
	}

	return session, nil
}

func (s *SessionManager) GetAllSessions() []*PlayerSession {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	allSessions := make([]*PlayerSession, 0, len(s.playerSessions))
	for _, session := range s.playerSessions {
		allSessions = append(allSessions, session)
	}

	return allSessions
}

func (s *SessionManager) CloseSession(sessionId string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.playerSessions[sessionId]
	if !exists {
		return fmt.Errorf("sessionId不存在")
	}

	session.Conn.Close()
	delete(s.playerSessions, sessionId)

	return nil
}

func (s *SessionManager) Shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, session := range s.playerSessions {
		session.Conn.Close()
	}

	s.playerSessions = make(map[string]*PlayerSession)
}
