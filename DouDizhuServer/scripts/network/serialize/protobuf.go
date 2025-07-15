package serialize

import (
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

// ProtoHandler 实现基于protobuf的消息处理器

func Deserialize(data []byte) (*protodef.PClientMsg, error) {
	request := &protodef.PClientMsg{}
	if err := proto.Unmarshal(data, request); err != nil {
		return nil, err
	}
	return request, nil
}

func Serialize(response *protodef.PServerMsg) ([]byte, error) {
	return proto.Marshal(response)
}

func SerializePayload(payload proto.Message) ([]byte, error) {
	return proto.Marshal(payload)
}

func DeserializePayload(msgId protodef.PMsgId, data []byte) (proto.Message, error) {
	payload := GetMessage(msgId)
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
