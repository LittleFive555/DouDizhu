package message

import (
	"DouDizhuServer/errors"
	"DouDizhuServer/network/protodef"
)

func CreateRespPacket(requestHeader *protodef.PGameMsgHeader) *protodef.PGameMsgRespPacket {
	return &protodef.PGameMsgRespPacket{
		Header: requestHeader,
	}
}

func CreateEmptyResponsePacket(requestHeader *protodef.PGameMsgHeader) *protodef.PGameMsgRespPacket {
	return &protodef.PGameMsgRespPacket{
		Header:  requestHeader,
		Content: &protodef.PGameMsgRespPacket_EmptyResponse{},
	}
}

func CreateErrorPacket(requestHeader *protodef.PGameMsgHeader, gameError *errors.GameError) *protodef.PGameMsgRespPacket {
	if gameError.Category == errors.CategoryGameplay {
		return &protodef.PGameMsgRespPacket{
			Header: requestHeader,
			Content: &protodef.PGameMsgRespPacket_Error{
				Error: &protodef.PError{
					Type:      protodef.PError_TYPE_BUSINESS,
					ErrorCode: string(gameError.Code),
					Message:   gameError.ClientMsg,
				},
			},
		}
	} else {
		return &protodef.PGameMsgRespPacket{
			Header: requestHeader,
			Content: &protodef.PGameMsgRespPacket_Error{
				Error: &protodef.PError{
					Type:      protodef.PError_TYPE_SERVER_ERROR,
					ErrorCode: string(gameError.Code),
					Message:   gameError.ClientMsg,
				},
			},
		}
	}
}
