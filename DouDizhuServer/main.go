package main

import (
	"DouDizhuServer/logger"
)

func main() {
	// 初始化日志
	if err := logger.InitLoggerWithConfig("config/log.yaml"); err != nil {
		panic(err)
	}
	logger.Info("日志系统初始化成功")
}
