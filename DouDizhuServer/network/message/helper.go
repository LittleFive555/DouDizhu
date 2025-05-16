package message

import (
	"DouDizhuServer/network/protodef"
)

func CreateRespPacket(requestHeader *protodef.PGameMsgHeader) *protodef.PGameMsgRespPacket {
	return &protodef.PGameMsgRespPacket{
		Header: requestHeader,
	}
}

func CreateErrorPacket(requestHeader *protodef.PGameMsgHeader, errorCode protodef.PError_Code, errorMessage string) *protodef.PGameMsgRespPacket {
	return &protodef.PGameMsgRespPacket{
		Header: requestHeader,
		Content: &protodef.PGameMsgRespPacket_Error{
			Error: &protodef.PError{
				Code:    errorCode,
				Message: errorMessage,
			},
		},
	}
}
