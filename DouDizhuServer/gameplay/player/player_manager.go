package player

import "github.com/google/uuid"

var Manager *PlayerManager

type PlayerManager struct {
	players map[string]*Player
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{players: make(map[string]*Player)}
}

func (m *PlayerManager) CreatePlayer(account string, password string) (*Player, error) {
	// TODO Save to database
	return NewPlayer(uuid.New().String(), "NewPlayer"), nil
}

func (m *PlayerManager) Login(account string, password string) (*Player, error) {
	// TODO Check password
	return NewPlayer(uuid.New().String(), "NewPlayer"), nil
}

func (m *PlayerManager) GetPlayer(playerId string) *Player {
	return m.players[playerId]
}

func (m *PlayerManager) RemovePlayer(playerId string) {
	delete(m.players, playerId)
}
