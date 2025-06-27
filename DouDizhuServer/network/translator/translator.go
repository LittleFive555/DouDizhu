package translator

import (
	"DouDizhuServer/gameplay/player"
	"DouDizhuServer/gameplay/room"
	"DouDizhuServer/network/protodef"
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
	players := RoomPlayersToProto(r, playerManager)
	return &protodef.PRoom{
		Id:             r.GetId(),
		Name:           r.GetName(),
		Owner:          PlayerToProto(playerManager.GetPlayer(r.GetOwnerId())),
		Players:        players,
		MaxPlayerCount: r.GetMaxPlayerCount(),
		State:          r.GetState(),
	}
}

func RoomPlayersToProto(r *room.Room, playerManager *player.PlayerManager) []*protodef.PPlayer {
	players := make([]*protodef.PPlayer, 0)
	for _, playerId := range r.GetPlayers() {
		players = append(players, PlayerToProto(playerManager.GetPlayer(playerId)))
	}
	return players
}
