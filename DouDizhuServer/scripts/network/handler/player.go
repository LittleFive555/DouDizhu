package handler

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/gameplay/player"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"
	"DouDizhuServer/scripts/network/translator"

	"google.golang.org/protobuf/proto"
)

func HandleRegister(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PRegisterRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	account := reqMsg.GetAccount()
	password := reqMsg.GetPassword()
	err := player.Manager.Register(account, password)
	if err != nil {
		return nil, err
	}
	return &message.HandleResult{}, nil
}

func HandleLogin(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PLoginRequest)
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeInvalidRequest)
	}
	account := reqMsg.GetAccount()
	password := reqMsg.GetPassword()

	player, err := player.Manager.Login(account, password, context.SessionId)
	if err != nil {
		return nil, err
	}

	result := &message.HandleResult{
		Resp: &protodef.PLoginResponse{
			Info: translator.PlayerToProto(player),
		},
	}
	return result, nil
}
