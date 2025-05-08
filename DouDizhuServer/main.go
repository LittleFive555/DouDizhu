package main

import (
	"DouDizhuServer/logger"
	"fmt"
	"net"
	"unicode/utf8"
)

func main() {
	// 初始化日志
	if err := logger.InitLoggerWithConfig("config/log.yaml"); err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("日志系统初始化成功")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	logger.InfoWith("新连接", "remote_addr", remoteAddr)

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				logger.InfoWith("客户端正常断开连接", "remote_addr", remoteAddr)
			} else {
				logger.ErrorWith("连接异常断开", "remote_addr", remoteAddr, "error", err)
			}
			break
		}

		// 将字节转换为UTF-8字符串
		utf8Str := string(buf[:n])
		if !utf8.ValidString(utf8Str) {
			logger.ErrorWith("收到非UTF-8编码数据", "remote_addr", remoteAddr, "data", fmt.Sprintf("%x", buf[:n]))
			continue
		}

		// 输出接收到的数据详情
		logger.InfoWith("收到数据", "remote_addr", remoteAddr, "data", utf8Str)
	}
}
