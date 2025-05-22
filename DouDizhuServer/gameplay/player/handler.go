package player

import (
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
)

func HandleRegister(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	account := req.GetRegisterReq().GetAccount()
	password := req.GetRegisterReq().GetPassword()
	result, err := Manager.Register(account, password)
	if err != nil {
		return message.CreateErrorPacket(req.Header, protodef.PError_CODE_DATABASE_WRITE_ERROR, err.Error()), nil
	}
	respPacket := message.CreateRespPacket(req.Header)
	respPacket.Content = &protodef.PGameMsgRespPacket_RegisterResp{
		RegisterResp: &protodef.PRegisterResponse{
			Result: result,
		},
	}
	return respPacket, nil
}

func HandleLogin(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	account := req.GetLoginReq().GetAccount()
	password := req.GetLoginReq().GetPassword()

	player, result, err := Manager.Login(account, password)
	if err != nil {
		return message.CreateErrorPacket(req.Header, protodef.PError_CODE_DATABASE_READ_ERROR, err.Error()), nil
	}
	playerId := ""
	if player != nil {
		playerId = player.PlayerId
	}

	respPacket := message.CreateRespPacket(req.Header)
	respPacket.Content = &protodef.PGameMsgRespPacket_LoginResp{
		LoginResp: &protodef.PLoginResponse{
			Result:   result,
			PlayerId: playerId,
		},
	}
	return respPacket, nil
}
