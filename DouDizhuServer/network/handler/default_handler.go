package handler

import (
	"DouDizhuServer/logger"
	"fmt"
	"unicode/utf8"
)

// DefaultHandler 默认的消息处理器实现
type DefaultHandler struct{}

// HandleMessage 默认的消息处理实现
func (h *DefaultHandler) HandleMessage(data []byte) ([]byte, error) {
	// 将字节转换为UTF-8字符串
	utf8Str := string(data)
	if !utf8.ValidString(utf8Str) {
		logger.ErrorWith("收到非UTF-8编码数据", "data", fmt.Sprintf("%x", data))
		return nil, fmt.Errorf("invalid UTF-8 data")
	}

	// 输出接收到的数据详情
	logger.InfoWith("收到数据", "data", utf8Str)
	return []byte("Received"), nil
}
