package main

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network"
	"DouDizhuServer/network/protodef"
)

func main() {
	// 初始化日志
	if err := logger.InitLoggerWithConfig("config/log.yaml"); err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("日志系统初始化成功")

	// 创建并启动TCP服务器
	server := network.NewGameServer(":8080")
	if err := server.Start(); err != nil {
		logger.ErrorWith("服务器启动失败", "error", err)
		panic(err)
	}
}

func handleChatMessage(req *protodef.GameMsgReqPacket, remoteAddr string) (*protodef.GameMsgRespPacket, error) {
	logger.InfoWith("收到聊天消息", "remote_addr", remoteAddr, "content", req.GetChatMsg().GetContent())

	return &protodef.GameMsgRespPacket{
		Header: &protodef.GameMsgHeader{},
		Content: &protodef.GameMsgRespPacket_ChatMsg{
			ChatMsg: &protodef.ChatMsgResponse{},
		},
	}, nil
}
