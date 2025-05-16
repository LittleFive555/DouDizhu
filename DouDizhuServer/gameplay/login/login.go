package login

import (
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
)

func HandleRegist(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	account := req.GetRegisterReq().GetAccount()
	password := req.GetRegisterReq().GetPassword()

	respPacket := message.CreateRespPacket(req.Header)
	result := ValidateAccount(account)
	if result != protodef.PRegisterResponse_RESULT_SUCCESS {
		respPacket.Content = &protodef.PGameMsgRespPacket_RegisterResp{
			RegisterResp: &protodef.PRegisterResponse{
				Result: result,
			},
		}
		return respPacket, nil
	}
	result = ValidatePassword(password)
	if result != protodef.PRegisterResponse_RESULT_SUCCESS {
		respPacket.Content = &protodef.PGameMsgRespPacket_RegisterResp{
			RegisterResp: &protodef.PRegisterResponse{
				Result: result,
			},
		}
		return respPacket, nil
	}

	_, err := player.Manager.CreatePlayer(account, password)
	if err != nil {
		return message.CreateErrorPacket(req.Header, protodef.PError_CODE_DATABASE_WRITE_ERROR, err.Error()), nil
	}
	respPacket.Content = &protodef.PGameMsgRespPacket_RegisterResp{
		RegisterResp: &protodef.PRegisterResponse{
			Result: protodef.PRegisterResponse_RESULT_SUCCESS,
		},
	}
	return respPacket, nil
}

func HandleLogin(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	account := req.GetLoginReq().GetAccount()
	password := req.GetLoginReq().GetPassword()

	player, err := player.Manager.Login(account, password)
	if err != nil {
		return message.CreateErrorPacket(req.Header, protodef.PError_CODE_DATABASE_READ_ERROR, err.Error()), nil
	}

	return &protodef.PGameMsgRespPacket{
		Header: &protodef.PGameMsgHeader{
			Player: player.ToProto(),
		},
		Content: &protodef.PGameMsgRespPacket_LoginResp{
			LoginResp: &protodef.PLoginResponse{
				Result: protodef.PLoginResponse_RESULT_SUCCESS,
			},
		},
	}, nil
}

func ValidateAccount(account string) protodef.PRegisterResponse_Result {
	// TODO
	return protodef.PRegisterResponse_RESULT_SUCCESS
}

func ValidatePassword(password string) protodef.PRegisterResponse_Result {
	// TODO
	return protodef.PRegisterResponse_RESULT_SUCCESS
}
