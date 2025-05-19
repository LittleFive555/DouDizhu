package database

import (
	"database/sql"
	"fmt"
)

type Account struct {
	Account        string
	HashedPassword string
	PlayerId       string
}

func AddPlayer(db *sql.DB, account *Account) error {
	_, err := db.Exec("INSERT INTO accounts(player_account, player_password_hash, player_id) VALUES (?, ?, ?)",
		account.Account, account.HashedPassword, account.PlayerId)
	if err != nil {
		return err
	}
	return nil
}

func GetPlayer(db *sql.DB, accountStr string) (Account, error) {
	var account Account
	row := db.QueryRow("SELECT * FROM accounts WHERE player_account = ?", accountStr)
	if err := row.Scan(&account.Account, &account.HashedPassword, &account.PlayerId); err != nil {
		if err == sql.ErrNoRows {
			return account, fmt.Errorf("GetPlayer %s: no such account", accountStr)
		}
		return account, fmt.Errorf("GetPlayer %s: %v", accountStr, err)
	}
	return account, nil
}
