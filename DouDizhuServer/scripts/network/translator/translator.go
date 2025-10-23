package translator

import (
	"DouDizhuServer/scripts/gameplay/player"
	"DouDizhuServer/scripts/gameplay/room"
	"DouDizhuServer/scripts/network/protodef"
)

func PlayerToProto(p *player.Player) *protodef.PPlayer {
	return &protodef.PPlayer{
		Id:       p.GetPlayerId(),
		Nickname: p.GetNickname(),
		State:    p.GetState(),
		RoomId:   p.GetRoomId(),
	}
}

func RoomToProto(r *room.Room, playerManager *player.PlayerManager) *protodef.PRoom {
	protoRoom := &protodef.PRoom{
		Id:             r.GetId(),
		Name:           r.GetName(),
		OwnerId:        r.GetOwnerId(),
		Players:        RoomPlayersToProto(r, playerManager),
		MaxPlayerCount: r.GetMaxPlayerCount(),
		State:          r.GetState(),
	}
	return protoRoom
}

func RoomPlayersToProto(r *room.Room, playerManager *player.PlayerManager) []*protodef.PPlayer {
	players := make([]*protodef.PPlayer, 0)
	for _, playerId := range r.GetPlayers() {
		players = append(players, PlayerToProto(playerManager.GetPlayer(playerId)))
	}
	return players
}
