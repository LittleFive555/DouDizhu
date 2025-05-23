package player

import (
	"DouDizhuServer/database"
	"DouDizhuServer/errors"
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

func (m *PlayerManager) Register(accountStr string, password string) error {
	// 验证账号格式，并且检查是否存在
	err := validateAccount(accountStr)
	if err != nil {
		return err
	}
	account, err := database.GetAccount(accountStr)
	if err != nil {
		return errors.NewDatabaseError(errors.CodeDBReadError, err)
	}
	if account.IsExists() {
		return errors.NewGameplayError(errors.CodeAccountExists)
	}

	// 验证密码
	err = validatePassword(password)
	if err != nil {
		return err
	}

	// 生成密码哈希
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewThirdPartyError(errors.CodeUnknown, err)
	}
	hashedPaswordStr := base64.StdEncoding.EncodeToString(hashedPasswordBytes)

	// 生成玩家ID
	playerId := uuid.New().String()

	account = database.Account{
		Account:        accountStr,
		HashedPassword: hashedPaswordStr,
		PlayerId:       playerId,
	}
	err = database.AddAccount(&account)
	if err != nil {
		return errors.NewDatabaseError(errors.CodeDBWriteError, err)
	}
	return nil
}

func (m *PlayerManager) Login(accountStr string, password string) (*Player, error) {
	// 判断账号是否存在
	account, err := database.GetAccount(accountStr)
	if err != nil {
		return nil, errors.NewDatabaseError(errors.CodeDBReadError, err)
	}
	if !account.IsExists() {
		return nil, errors.NewGameplayError(errors.CodeAccountNotExists)
	}
	// 验证密码
	hashedPassword, err := base64.StdEncoding.DecodeString(account.HashedPassword)
	if err != nil {
		return nil, errors.NewThirdPartyError(errors.CodeUnknown, err)
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return nil, errors.NewGameplayError(errors.CodePasswordWrong)
	}

	player := NewPlayer(account.PlayerId, "NewPlayer")
	m.players[account.PlayerId] = player
	return player, nil
}

func (m *PlayerManager) GetPlayer(playerId string) *Player {
	return m.players[playerId]
}

func (m *PlayerManager) RemovePlayer(playerId string) {
	delete(m.players, playerId)
}

func validateAccount(account string) error {
	// TODO
	return nil
}

func validatePassword(password string) error {
	// TODO
	return nil
}
