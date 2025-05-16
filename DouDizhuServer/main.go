package main

import (
	"DouDizhuServer/gameplay"
	"DouDizhuServer/gameplay/login"
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/logger"
	"DouDizhuServer/network"
	"DouDizhuServer/network/protodef"
	"reflect"
)

func main() {
	// 初始化日志
	if err := logger.InitLoggerWithConfig("config/log.yaml"); err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("日志系统初始化成功")

	// 创建并启动TCP服务器
	network.Server = network.NewGameServer()

	network.Server.RegisterHandler(reflect.TypeOf(protodef.PGameClientMessage_ChatMsg{}), gameplay.HandleChatMessage)
	network.Server.RegisterHandler(reflect.TypeOf(protodef.PGameClientMessage_RegisterReq{}), login.HandleRegist)
	network.Server.RegisterHandler(reflect.TypeOf(protodef.PGameClientMessage_LoginReq{}), login.HandleLogin)

	player.Manager = player.NewPlayerManager()

	if err := network.Server.Start(":8080"); err != nil {
		logger.ErrorWith("服务器启动失败", "error", err)
		panic(err)
	}

	defer network.Server.Stop()
}
