package database

type Account struct {
	Account        string
	HashedPassword string
	PlayerId       string
}

func AddPlayer(account *Account) error {
	_, err := GetDB().Exec("INSERT INTO accounts(player_account, player_password_hash, player_id) VALUES (?, ?, ?)",
		account.Account, account.HashedPassword, account.PlayerId)
	if err != nil {
		return err
	}
	return nil
}

func GetPlayer(accountStr string) (Account, error) {
	var account Account
	row := GetDB().QueryRow("SELECT * FROM accounts WHERE player_account = ?", accountStr)
	if err := row.Scan(&account.Account, &account.HashedPassword, &account.PlayerId); err != nil {
		return account, err
	}
	return account, nil
}
