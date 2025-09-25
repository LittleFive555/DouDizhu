package message

import (
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

type Notification struct {
	NotifyMsgId protodef.PMsgId
	Payload     proto.Message
	NotifyGroup INotificationGroup
}
