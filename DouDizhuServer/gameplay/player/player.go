package player

import "DouDizhuServer/network/protodef"

type Player struct {
	PlayerId string
	Nickname string
}

func NewPlayer(playerId string, nickname string) *Player {
	return &Player{
		PlayerId: playerId,
		Nickname: nickname,
	}
}

func (p *Player) ToProto() *protodef.PPlayer {
	return &protodef.PPlayer{
		PlayerId: p.PlayerId,
		Nickname: p.Nickname,
	}
}
