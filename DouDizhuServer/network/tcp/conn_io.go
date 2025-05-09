package tcp

import (
	"encoding/binary"
	"io"
	"net"
)

type ConnIO interface {
	Read(conn net.Conn) ([]byte, error)
	Write(conn net.Conn, data []byte) error
}

type LengthPrefixConnIO struct {
	ConnIO
}

func NewLengthPrefixConnIO() *LengthPrefixConnIO {
	return &LengthPrefixConnIO{}
}

func (c *LengthPrefixConnIO) Read(conn net.Conn) ([]byte, error) {
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

func (c *LengthPrefixConnIO) Write(conn net.Conn, data []byte) error {
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
