package player

import (
	"DouDizhuServer/database"
	"DouDizhuServer/network/protodef"
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

func (m *PlayerManager) Register(accountStr string, password string) (protodef.PRegisterResponse_Result, error) {
	result := validateAccount(accountStr)
	if result != protodef.PRegisterResponse_RESULT_SUCCESS {
		return result, nil
	}
	result = validatePassword(password)
	if result != protodef.PRegisterResponse_RESULT_SUCCESS {
		return result, nil
	}

	// 生成哈希
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return protodef.PRegisterResponse_RESULT_UNKNOWN, err
	}
	hashedPaswordStr := base64.StdEncoding.EncodeToString(hashedPasswordBytes)

	playerId := uuid.New().String()

	account := database.Account{
		Account:        accountStr,
		HashedPassword: hashedPaswordStr,
		PlayerId:       playerId,
	}
	err = database.AddPlayer(&account)
	if err != nil {
		return protodef.PRegisterResponse_RESULT_ACCOUNT_EXISTS, nil
	}
	return protodef.PRegisterResponse_RESULT_SUCCESS, nil
}

func (m *PlayerManager) Login(accountStr string, password string) (*Player, protodef.PLoginResponse_Result, error) {
	account, err := database.GetPlayer(accountStr)
	if err != nil {
		return nil, protodef.PLoginResponse_RESULT_ACCOUNT_NOT_EXISTS, nil
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(account.HashedPassword)
	if err != nil {
		return nil, protodef.PLoginResponse_RESULT_UNKNOWN, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return nil, protodef.PLoginResponse_RESULT_PASSWORD_WRONG, nil
	}
	player := NewPlayer(account.PlayerId, "NewPlayer")
	m.players[account.PlayerId] = player
	return player, protodef.PLoginResponse_RESULT_SUCCESS, nil
}

func (m *PlayerManager) GetPlayer(playerId string) *Player {
	return m.players[playerId]
}

func (m *PlayerManager) RemovePlayer(playerId string) {
	delete(m.players, playerId)
}

func validateAccount(account string) protodef.PRegisterResponse_Result {
	// TODO
	return protodef.PRegisterResponse_RESULT_SUCCESS
}

func validatePassword(password string) protodef.PRegisterResponse_Result {
	// TODO
	return protodef.PRegisterResponse_RESULT_SUCCESS
}
