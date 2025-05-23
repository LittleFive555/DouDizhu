package player

import (
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
)

func HandleRegister(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, *protodef.PGameNotificationPacket, error) {
	account := req.GetRegisterReq().GetAccount()
	password := req.GetRegisterReq().GetPassword()
	err := Manager.Register(account, password)
	if err != nil {
		return nil, nil, err
	}
	return message.CreateEmptyRespPacket(req.Header), nil, nil
}

func HandleLogin(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, *protodef.PGameNotificationPacket, error) {
	account := req.GetLoginReq().GetAccount()
	password := req.GetLoginReq().GetPassword()

	player, err := Manager.Login(account, password)
	if err != nil {
		return nil, nil, err
	}

	loginResp := &protodef.PLoginResponse{
		PlayerId: player.PlayerId,
	}
	respPacket := message.CreateRespPacket(req.Header)
	respPacket.Content = &protodef.PGameMsgRespPacket_LoginResp{
		LoginResp: loginResp,
	}
	return respPacket, nil, nil
}
