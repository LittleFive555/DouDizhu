package network

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/handler"
	"fmt"
	"net"
)

// TCPServer TCP服务器结构体
type TCPServer struct {
	addr           string
	listener       net.Listener
	messageHandler handler.Handler
}

// NewTCPServer 创建一个新的TCP服务器实例
func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		addr:           addr,
		messageHandler: &handler.DefaultHandler{},
	}
}

// SetMessageHandler 设置消息处理器
func (s *TCPServer) SetMessageHandler(handler handler.Handler) {
	s.messageHandler = handler
}

// Start 启动TCP服务器
func (s *TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("启动服务器失败: %v", err)
	}
	s.listener = ln

	logger.InfoWith("TCP服务器启动成功", "addr", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.ErrorWith("接受连接失败", "error", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// Stop 停止TCP服务器
func (s *TCPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// handleConnection 处理单个连接
func (s *TCPServer) handleConnection(conn net.Conn) {
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

		// 使用消息处理器处理数据
		if err := s.messageHandler.HandleMessage(buf[:n], remoteAddr); err != nil {
			logger.ErrorWith("处理消息失败", "remote_addr", remoteAddr, "error", err)
			continue
		}
	}
}
