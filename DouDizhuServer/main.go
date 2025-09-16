package main

import (
	"DouDizhuServer/scripts/database"
	"DouDizhuServer/scripts/gameplay/player"
	"DouDizhuServer/scripts/gameplay/playground"
	"DouDizhuServer/scripts/gameplay/room"
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network"

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

	// startGameServer()
	startPlayground()
}

func startGameServer() {
	var server *network.GameServer
	// 创建并启动TCP服务器
	server = network.NewGameServer()

	// 注册消息处理函数
	server.RegisterHandlers()

	player.Manager = player.NewPlayerManager()
	room.Manager = room.NewRoomManager()

	if err := server.Start(":8080"); err != nil {
		logger.PanicWith("服务器启动失败", "error", err)
	}

	defer server.Stop()
}

func startPlayground() {
	var server *network.GameServer
	server = network.NewGameServer()
	server.RegisterPlaygroundHandlers()

	playground.Playground = playground.NewRoomPlayground()

	if err := server.Start(":9090"); err != nil {
		logger.PanicWith("服务器启动失败", "error", err)
	}

	defer server.Stop()
}
