package session

import (
	"DouDizhuServer/cypher"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/message"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/hkdf"
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
	secureKey []byte

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

func (s *PlayerSession) GenerateSecureKey(clientPublicKeyStr string, salt []byte, info []byte) (string, error) {
	logger.InfoWith("客户端公钥", "clientPublicKeyStr", clientPublicKeyStr, "salt", salt, "info", info)
	serverPrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}
	clientPublicKey, err := stringToPublicKey(clientPublicKeyStr)
	if err != nil {
		return "", err
	}
	serverSharedKey, err := serverPrivateKey.ECDH(clientPublicKey)
	if err != nil {
		return "", err
	}
	logger.InfoWith("生成共享密钥", "serverSharedKey", serverSharedKey)

	serverPublicKey := serverPrivateKey.PublicKey()
	serverPublicKeyStr := publicKeyToString(serverPublicKey)
	logger.InfoWith("服务端公钥", "serverPublicKeyStr", serverPublicKeyStr)

	secureKey := deriveSecureKey(serverSharedKey, salt, info)
	logger.InfoWith("生成密钥", "secureKey", secureKey)
	s.secureKey = secureKey
	return serverPublicKeyStr, nil
}

func (s *PlayerSession) EncryptPayload(data []byte) ([]byte, []byte, error) {
	return cypher.AesEncryptCBC(data, s.secureKey)
}

func (s *PlayerSession) DecryptPayload(cyphertext []byte, iv []byte) ([]byte, error) {
	return cypher.AesDecryptCBC(cyphertext, iv, s.secureKey)
}

func (s *PlayerSession) IsSecureKeyValid() bool {
	return s.secureKey != nil
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

func stringToPublicKey(str string) (*ecdh.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	publicKey, err := ecdh.P256().NewPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func publicKeyToString(publicKey *ecdh.PublicKey) string {
	return base64.StdEncoding.EncodeToString(publicKey.Bytes())
}

func deriveSecureKey(sharedKey []byte, salt []byte, info []byte) []byte {
	hkdf := hkdf.New(sha256.New, sharedKey, salt, info)

	derivedKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, derivedKey); err != nil {
		return nil
	}
	return derivedKey
}
