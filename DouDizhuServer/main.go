package main

import (
	"DouDizhuServer/database"
	"DouDizhuServer/gameplay/chat"
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/logger"
	"DouDizhuServer/network"
	"DouDizhuServer/network/protodef"

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

	network.Server.RegisterHandler(protodef.PMsgId_PMSG_ID_HANDSHAKE, network.HandleHandshake)
	network.Server.RegisterHandler(protodef.PMsgId_PMSG_ID_CHAT_MSG, chat.HandleChatMessage)
	network.Server.RegisterHandler(protodef.PMsgId_PMSG_ID_REGISTER, player.HandleRegister)
	network.Server.RegisterHandler(protodef.PMsgId_PMSG_ID_LOGIN, player.HandleLogin)

	player.Manager = player.NewPlayerManager()

	if err := network.Server.Start(":8080"); err != nil {
		logger.PanicWith("服务器启动失败", "error", err)
	}

	defer network.Server.Stop()
}
