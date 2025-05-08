package handler

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/protodef"
	"errors"
	"reflect"

	"google.golang.org/protobuf/proto"
)

// ProtoHandler 实现基于protobuf的消息处理器
type ProtoHandler struct {
	handlers map[reflect.Type]func(*protodef.GameMsgReqPacket) (*protodef.GameMsgRespPacket, error)
}

// NewProtoHandler 创建一个新的ProtoHandler实例
func NewProtoHandler() *ProtoHandler {
	return &ProtoHandler{
		handlers: make(map[reflect.Type]func(*protodef.GameMsgReqPacket) (*protodef.GameMsgRespPacket, error)),
	}
}

// RegisterHandler 注册消息处理器
func (h *ProtoHandler) RegisterHandler(msgType reflect.Type, handler func(*protodef.GameMsgReqPacket) (*protodef.GameMsgRespPacket, error)) {
	h.handlers[msgType] = handler
}

// HandleMessage 实现Handler接口
func (h *ProtoHandler) HandleMessage(data []byte) ([]byte, error) {
	// 解析请求包
	reqPacket := &protodef.GameMsgReqPacket{}

	if err := proto.Unmarshal(data, reqPacket); err != nil {
		logger.ErrorWith("解析消息失败", "error", err)
		return nil, err
	}

	// 查找对应的消息处理器
	content := reqPacket.GetContent()
	if content == nil {
		logger.ErrorWith("消息内容为空")
		return nil, errors.New("消息内容为空")
	}

	// 获取消息的实际类型
	msgType := reflect.TypeOf(content).Elem()
	handler, exists := h.handlers[msgType]
	if !exists {
		logger.ErrorWith("未找到消息处理器", "type", msgType.String())
		return nil, errors.New("未找到消息处理器")
	}

	// 处理消息
	respPacket, err := handler(reqPacket)
	if err != nil {
		logger.ErrorWith("处理消息失败", "error", err)
		return nil, err
	}

	// 序列化响应包
	responseData, err := proto.Marshal(respPacket)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return nil, err
	}

	return responseData, nil
}
