package message

import (
	"DouDizhuServer/errors"
	"DouDizhuServer/logger"
	"DouDizhuServer/network/protodef"

	"google.golang.org/protobuf/proto"
)

func CreateRespMessage(requestHeader *protodef.PMsgHeader, respPacket proto.Message) *protodef.PServerMsg {
	serialized, err := proto.Marshal(respPacket)
	if err != nil {
		logger.ErrorWith("序列化响应失败", "error", err)
		return nil
	}
	return &protodef.PServerMsg{
		Header: &protodef.PMsgHeader{
			UniqueId:  requestHeader.UniqueId,
			MsgId:     requestHeader.MsgId,
			SessionId: requestHeader.SessionId,
			PlayerId:  requestHeader.PlayerId,
			Timestamp: requestHeader.Timestamp,
		},
		MsgType: protodef.PServerMsgType_PSERVER_MSG_TYPE_RESPONSE,
		Payload: serialized,
	}
}

func CreateNotificationMessage(requestHeader *protodef.PMsgHeader, notification proto.Message) *protodef.PServerMsg {
	if notification == nil {
		return nil
	}
	serialized, err := proto.Marshal(notification)
	if err != nil {
		logger.ErrorWith("序列化通知响应失败", "error", err)
		return nil
	}
	return &protodef.PServerMsg{
		Header: &protodef.PMsgHeader{
			UniqueId:  requestHeader.UniqueId,
			MsgId:     requestHeader.MsgId,
			SessionId: requestHeader.SessionId,
			PlayerId:  requestHeader.PlayerId,
			Timestamp: requestHeader.Timestamp,
		},
		MsgType: protodef.PServerMsgType_PSERVER_MSG_TYPE_NOTIFICATION,
		Payload: serialized,
	}
}

func CreateErrorMessage(requestHeader *protodef.PMsgHeader, gameError *errors.GameError) *protodef.PServerMsg {
	var pError *protodef.PError
	if gameError.Category == errors.CategoryGameplay {
		pError = &protodef.PError{
			Type:      protodef.PError_TYPE_BUSINESS,
			ErrorCode: string(gameError.Code),
			Message:   gameError.ClientMsg,
		}
	} else {
		pError = &protodef.PError{
			Type:      protodef.PError_TYPE_SERVER_ERROR,
			ErrorCode: string(gameError.Code),
			Message:   gameError.ClientMsg,
		}
	}
	serialized, err := proto.Marshal(pError)
	if err != nil {
		logger.ErrorWith("序列化错误响应失败", "error", err)
		return nil
	}
	return &protodef.PServerMsg{
		Header: &protodef.PMsgHeader{
			UniqueId:  requestHeader.UniqueId,
			MsgId:     requestHeader.MsgId,
			SessionId: requestHeader.SessionId,
			PlayerId:  requestHeader.PlayerId,
			Timestamp: requestHeader.Timestamp,
		},
		MsgType: protodef.PServerMsgType_PSERVER_MSG_TYPE_ERROR,
		Payload: serialized,
	}
}
