package message

import (
	"google.golang.org/protobuf/proto"
)

type HandleResult struct {
	Resp   proto.Message
	Notify proto.Message
}
