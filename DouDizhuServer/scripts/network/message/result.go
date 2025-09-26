package message

import (
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

type HandleResult struct {
	Resp        proto.Message
	NotifyMsgId protodef.PMsgId
	NotifyGroup INotifyGroup
	Notify      proto.Message
}
