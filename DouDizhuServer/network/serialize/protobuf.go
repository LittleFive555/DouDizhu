package serialize

import (
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

// ProtoHandler 实现基于protobuf的消息处理器

func Deserialize(data []byte) (*protodef.PGameClientMessage, error) {
	request := &protodef.PGameClientMessage{}
	if err := proto.Unmarshal(data, request); err != nil {
		return nil, err
	}
	return request, nil
}

func Serialize(response *protodef.PGameServerMessage) ([]byte, error) {
	responseData, err := proto.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
