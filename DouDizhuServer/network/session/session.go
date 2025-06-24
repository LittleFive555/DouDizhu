package session

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type PlayerState int

const (
	PlayerState_Connecting PlayerState = iota
	PlayerState_Lobby
	PlayerState_InGame
	PlayerState_Disconnected
)

type PlayerSession struct {
	Id        string
	Conn      net.Conn
	IP        string
	State     PlayerState
	SharedKey []byte

	PlayerId string
}

func (s *PlayerSession) StartReading(handler func(message *message.Message)) {
	defer s.Conn.Close()

	remoteAddr := s.Conn.RemoteAddr().String()
	for {
		rawMessage, err := read(s.Conn)
		if err != nil {
			if err.Error() == "EOF" {
				logger.InfoWith("客户端正常断开连接", "remote_addr", remoteAddr)
			} else {
				logger.ErrorWith("连接异常断开", "remote_addr", remoteAddr, "error", err)
			}
			break
		}

		handler(&message.Message{
			SessionId: s.Id,
			Data:      rawMessage,
		})
	}
}

func (s *PlayerSession) SendMessage(data []byte) error {
	if s.Conn == nil {
		return fmt.Errorf("连接已关闭")
	}
	return write(s.Conn, data)
}

func (s *PlayerSession) Close() error {
	if s.Conn != nil {
		return s.Conn.Close()
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
