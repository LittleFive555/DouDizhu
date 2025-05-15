package tcp

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// TCPServer TCP服务器结构体
type TCPServer struct {
	addr     string
	listener net.Listener
	conns    map[string]net.Conn
}

// NewTCPServer 创建一个新的TCP服务器实例
func NewTCPServer(addr string) *TCPServer {
	return &TCPServer{
		addr: addr,
	}
}

// Start 启动TCP服务器
func (s *TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("启动服务器失败: %v", err)
	}
	s.listener = ln
	s.conns = make(map[string]net.Conn)

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
	s.conns[remoteAddr] = conn
	logger.InfoWith("新连接", "remote_addr", remoteAddr)

	for {
		rawMessage, err := read(conn)
		if err != nil {
			if err.Error() == "EOF" {
				logger.InfoWith("客户端正常断开连接", "remote_addr", remoteAddr)
			} else {
				logger.ErrorWith("连接异常断开", "remote_addr", remoteAddr, "error", err)
			}
			delete(s.conns, remoteAddr)
			break
		}

		message.Dispatcher.PostMessage(&message.Message{
			SessionId: remoteAddr,
			Data:      rawMessage,
		})
	}
}

func (s *TCPServer) SendMessage(sessionId string, data []byte) error {
	conn, ok := s.conns[sessionId]
	if !ok {
		return fmt.Errorf("sessionId不存在")
	}
	return write(conn, data)
}

func (s *TCPServer) SendMessageToAll(data []byte) error {
	for _, conn := range s.conns {
		write(conn, data)
	}
	return nil
}

func read(conn net.Conn) ([]byte, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)
	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func write(conn net.Conn, data []byte) error {
	lengthBuf := make([]byte, 4)
	length := uint32(len(data))
	binary.BigEndian.PutUint32(lengthBuf, length)
	_, err := conn.Write(lengthBuf)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}
