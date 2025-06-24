package message

import (
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

type HandleResult struct {
	Resp        proto.Message
	NofityMsgId protodef.PMsgId
	Notify      proto.Message
}
