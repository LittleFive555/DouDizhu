package database

import (
	"DouDizhuServer/logger"
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

var DBInstance *sql.DB

func ConnectDB() *sql.DB {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "doudizhu_db"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		logger.Fatal(pingErr)
	}

	logger.Info("数据库连接成功")
	return db
}
