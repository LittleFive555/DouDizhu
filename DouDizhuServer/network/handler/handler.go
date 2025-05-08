package handler

// MessageHandler 消息处理器接口
type Handler interface {
	// HandleMessage 处理接收到的消息
	HandleMessage(data []byte) ([]byte, error)
}
