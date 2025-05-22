package player

import (
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
)

func HandleRegister(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	account := req.GetRegisterReq().GetAccount()
	password := req.GetRegisterReq().GetPassword()
	err := Manager.Register(account, password)
	if err != nil {
		return nil, err
	}
	respPacket := message.CreateRespPacket(req.Header)
	respPacket.Content = &protodef.PGameMsgRespPacket_EmptyResponse{}
	return respPacket, nil
}

func HandleLogin(req *protodef.PGameClientMessage) (*protodef.PGameMsgRespPacket, error) {
	account := req.GetLoginReq().GetAccount()
	password := req.GetLoginReq().GetPassword()

	player, err := Manager.Login(account, password)
	if err != nil {
		return nil, err
	}
	playerId := player.PlayerId

	respPacket := message.CreateRespPacket(req.Header)
	respPacket.Content = &protodef.PGameMsgRespPacket_LoginResp{
		LoginResp: &protodef.PLoginResponse{
			PlayerId: playerId,
		},
	}
	return respPacket, nil
}
