package database

import (
	"database/sql"
)

type Account struct {
	Account        string
	HashedPassword string
	PlayerId       string
}

func AddAccount(account *Account) error {
	_, err := GetDB().Exec("INSERT INTO accounts(player_account, player_password_hash, player_id) VALUES (?, ?, ?)",
		account.Account, account.HashedPassword, account.PlayerId)
	if err != nil {
		return err
	}
	return nil
}

func GetAccount(accountStr string) (Account, error) {
	var account Account
	row := GetDB().QueryRow("SELECT * FROM accounts WHERE player_account = ?", accountStr)
	if err := row.Scan(&account.Account, &account.HashedPassword, &account.PlayerId); err != nil {
		if err == sql.ErrNoRows {
			return account, nil
		}
		return account, err
	}
	return account, nil
}

func (account *Account) IsExists() bool {
	return account.Account != ""
}
