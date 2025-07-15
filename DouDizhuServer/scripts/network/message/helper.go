package message

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/network/protodef"
)

func CreateErrorPayload(gameError *errordef.GameError) *protodef.PError {
	var pError protodef.PError
	if gameError.Category == errordef.CategoryGameplay {
		pError = protodef.PError{
			Type:      protodef.PError_TYPE_BUSINESS,
			ErrorCode: string(gameError.Code),
			Message:   gameError.ClientMsg,
		}
	} else {
		pError = protodef.PError{
			Type:      protodef.PError_TYPE_SERVER_ERROR,
			ErrorCode: string(gameError.Code),
			Message:   gameError.ClientMsg,
		}
	}
	return &pError
}
