package message

import (
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

type INotificationGroup interface {
	GetTargetSessionIds() []string
}

type HandleResult struct {
	Resp        proto.Message
	NotifyMsgId protodef.PMsgId
	NotifyGroup INotificationGroup
	Notify      proto.Message
}
