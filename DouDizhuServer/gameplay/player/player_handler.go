package player

import (
	"DouDizhuServer/network/message"
	"DouDizhuServer/network/protodef"
	"errors"

	"google.golang.org/protobuf/proto"
)

func HandleRegister(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PRegisterRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	account := reqMsg.GetAccount()
	password := reqMsg.GetPassword()
	err := Manager.Register(account, password)
	if err != nil {
		return nil, err
	}
	return &message.HandleResult{}, nil
}

func HandleLogin(context *message.MessageContext, req *proto.Message) (*message.HandleResult, error) {
	reqMsg, ok := (*req).(*protodef.PLoginRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}
	account := reqMsg.GetAccount()
	password := reqMsg.GetPassword()

	player, err := Manager.Login(account, password)
	if err != nil {
		return nil, err
	}

	result := &message.HandleResult{
		Resp: &protodef.PLoginResponse{
			PlayerId: player.PlayerId,
		},
	}
	return result, nil
}
