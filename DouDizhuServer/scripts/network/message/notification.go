package message

import (
	"DouDizhuServer/scripts/network/protodef"

	"google.golang.org/protobuf/proto"
)

type INotifyGroup interface {
	IsNotifyGroup()
}

type Notification struct {
	NotifyMsgId protodef.PMsgId
	Payload     proto.Message
	NotifyGroup INotifyGroup
}

type AllConnectionNotificationGroup struct {
}

func (g *AllConnectionNotificationGroup) IsNotifyGroup() {
}
