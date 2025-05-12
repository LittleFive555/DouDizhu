package gameplay

import (
	"DouDizhuServer/logger"
	"DouDizhuServer/network/protodef"
)

func HandleChatMessage(req *protodef.GameMsgReqPacket) (*protodef.GameMsgRespPacket, error) {
	logger.InfoWith("收到聊天消息", "content", req.GetChatMsg().GetContent())

	return &protodef.GameMsgRespPacket{
		Header: &protodef.GameMsgHeader{},
		Content: &protodef.GameMsgRespPacket_CommonResponse{
			CommonResponse: &protodef.CommonResponse{},
		},
	}, nil
}
