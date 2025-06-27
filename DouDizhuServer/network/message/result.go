package message

import (
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

type NotificationGroup interface { // TODO 其实应该用session，暂时先这样，后续重构session部分
	GetTargetPlayerIds() []string
}

type HandleResult struct {
	Resp        proto.Message
	NofityMsgId protodef.PMsgId
	NotifyGroup NotificationGroup
	Notify      proto.Message
}
