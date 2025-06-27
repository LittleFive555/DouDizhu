package main

import (
	"DouDizhuServer/database"
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/gameplay/room"
	"DouDizhuServer/logger"
	"DouDizhuServer/network"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.PanicWith("环境加载失败", "error", err)
	}
	// 初始化日志
	if err := logger.InitLoggerWithConfig("config/log.yaml"); err != nil {
		logger.PanicWith("日志系统初始化失败", "error", err)
	}
	defer logger.Sync()
	logger.Info("日志系统初始化成功")
	database.ConnectDB()

	// 创建并启动TCP服务器
	network.Server = network.NewGameServer()

	// 注册消息处理函数
	network.Server.RegisterHandlers()

	player.Manager = player.NewPlayerManager()
	room.Manager = room.NewRoomManager()

	if err := network.Server.Start(":8080"); err != nil {
		logger.PanicWith("服务器启动失败", "error", err)
	}

	defer network.Server.Stop()
}
