package main

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network"
)

func main() {
	// 初始化日志
	if err := logger.InitLoggerWithConfig("config/log.yaml"); err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("日志系统初始化成功")

	// 创建并启动TCP服务器
	server := network.NewTCPServer(":8080")
	if err := server.Start(); err != nil {
		logger.ErrorWith("服务器启动失败", "error", err)
		panic(err)
	}
}
