package database

import (
	"DouDizhuServer/scripts/logger"
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

var instance *sql.DB

func ConnectDB() *sql.DB {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "doudizhu_db"

	var err error
	instance, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Fatal(err)
	}

	err = instance.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("数据库连接成功")
	return instance
}

func GetDB() *sql.DB {
	return instance
}
