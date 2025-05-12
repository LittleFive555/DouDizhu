package serialize

import (
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

// ProtoHandler 实现基于protobuf的消息处理器

func Deserialize(data []byte) (*protodef.GameClientMessage, error) {
	request := &protodef.GameClientMessage{}
	if err := proto.Unmarshal(data, request); err != nil {
		return nil, err
	}
	return request, nil
}

func Serialize(response *protodef.GameServerMessage) ([]byte, error) {
	responseData, err := proto.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
