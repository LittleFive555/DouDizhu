package session

import (
	"DouDizhuServer/scripts/cypher"
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network/message"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/crypto/hkdf"
)

type PlayerSession struct {
	Id        string
	conn      net.Conn
	ip        string
	secureKey []byte

	closed    chan int
	writeChan chan []byte
	waitGroup sync.WaitGroup
}

func newPlayerSession(id string, conn net.Conn, ip string) *PlayerSession {
	return &PlayerSession{
		Id:   id,
		conn: conn,
		ip:   ip,

		closed:    make(chan int),
		writeChan: make(chan []byte, 100),
		waitGroup: sync.WaitGroup{},
	}
}

func (s *PlayerSession) SendMessage(data []byte) error {
	if s.conn == nil {
		return fmt.Errorf("连接已关闭")
	}
	s.writeChan <- data
	return nil
}

func (s *PlayerSession) GenerateSecureKey(clientPublicKeyStr string, salt []byte, info []byte) (string, error) {
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

	serverPublicKey := serverPrivateKey.PublicKey()
	serverPublicKeyStr := publicKeyToString(serverPublicKey)

	secureKey := deriveSecureKey(serverSharedKey, salt, info)
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

func (s *PlayerSession) start(messageHandler func(*message.Message)) {
	s.waitGroup.Add(2)
	go s.startReading(messageHandler)
	go s.startWriting()
	s.waitGroup.Wait()
}

func (s *PlayerSession) startReading(messageHandler func(*message.Message)) {
	defer s.waitGroup.Done()
	remoteAddr := s.conn.RemoteAddr().String()
	for {
		rawMessage, err := read(s.conn)
		if err != nil {
			if err.Error() == "EOF" {
				logger.InfoWith("客户端正常断开连接", "remote_addr", remoteAddr)
			} else {
				logger.ErrorWith("连接异常断开", "remote_addr", remoteAddr, "error", err)
			}
			s.close()
			return
		}

		messageHandler(&message.Message{
			SessionId: s.Id,
			Data:      rawMessage,
		})
	}
}

func (s *PlayerSession) startWriting() {
	defer s.waitGroup.Done()
	remoteAddr := s.conn.RemoteAddr().String()
	for {
		select {
		case <-s.closed:
			return
		case data := <-s.writeChan:
			err := write(s.conn, data)
			if err != nil {
				logger.ErrorWith("发送消息失败", "remote_addr", remoteAddr, "error", err)
				return
			}
		}
	}
}

func (s *PlayerSession) close() error {
	if s.conn != nil {
		s.closed <- 0
		conn := s.conn
		s.conn = nil
		return conn.Close()
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
