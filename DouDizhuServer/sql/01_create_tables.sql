USE doudizhu_db;

-- 账号表
CREATE TABLE IF NOT EXISTS accounts(
    player_account VARCHAR(64) NOT NULL PRIMARY KEY,
    player_password_hash TEXT NOT NULL,
    player_id TEXT NOT NULL
);