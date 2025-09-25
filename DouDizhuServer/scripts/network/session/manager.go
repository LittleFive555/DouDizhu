package session

import (
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network/message"
	"fmt"
	"net"
	"sync"
)

type SessionManager struct {
	playerSessions map[string]*PlayerSession
	mutex          sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{playerSessions: make(map[string]*PlayerSession)}
}

func (s *SessionManager) StartPlayerSession(sessionId string, conn net.Conn, receiveChan chan<- *message.Message) {
	s.mutex.Lock()

	ip := conn.RemoteAddr().String()
	session := newPlayerSession(sessionId, conn, ip)
	s.playerSessions[sessionId] = session

	s.mutex.Unlock()

	logger.InfoWith("创建会话成功，开始处理消息", "sessionId", session.Id)
	session.start(receiveChan)
}

func (s *SessionManager) GetSession(sessionId string) (*PlayerSession, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.playerSessions[sessionId]
	if !exists {
		logger.ErrorWith("sessionId不存在", "sessionId", sessionId)
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
		logger.ErrorWith("sessionId不存在", "sessionId", sessionId)
		return fmt.Errorf("sessionId不存在")
	}

	session.close()
	delete(s.playerSessions, sessionId)
	logger.InfoWith("关闭会话成功", "sessionId", sessionId)

	return nil
}

func (s *SessionManager) Shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, session := range s.playerSessions {
		session.close()
	}

	s.playerSessions = make(map[string]*PlayerSession)
}
