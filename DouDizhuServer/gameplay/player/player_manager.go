package player

import (
	"DouDizhuServer/database"
	"encoding/base64"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var Manager *PlayerManager

type PlayerManager struct {
	players map[string]*Player
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{players: make(map[string]*Player)}
}

func (m *PlayerManager) CreatePlayer(accountStr string, password string) (*Player, error) {
	// 生成哈希
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPaswordStr := base64.StdEncoding.EncodeToString(hashedPasswordBytes)

	playerId := uuid.New().String()

	account := database.Account{
		Account:        accountStr,
		HashedPassword: hashedPaswordStr,
		PlayerId:       playerId,
	}
	database.AddPlayer(database.DBInstance, &account)
	return NewPlayer(playerId, "NewPlayer"), nil
}

func (m *PlayerManager) Login(accountStr string, password string) (*Player, error) {
	account, err := database.GetPlayer(database.DBInstance, accountStr)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(account.HashedPassword)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return nil, err
	}
	return NewPlayer(account.PlayerId, "NewPlayer"), nil
}

func (m *PlayerManager) GetPlayer(playerId string) *Player {
	return m.players[playerId]
}

func (m *PlayerManager) RemovePlayer(playerId string) {
	delete(m.players, playerId)
}
