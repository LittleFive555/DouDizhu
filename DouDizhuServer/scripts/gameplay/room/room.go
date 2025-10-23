package room

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/gameplay/player"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"
	"sync"
)

type Room struct {
	id      uint32
	name    string
	ownerId string

	players        map[string]*player.Player
	maxPlayerCount uint32

	state protodef.PRoomState
	lock  sync.RWMutex
}

func NewRoom(id uint32, name string, dispatcher message.INotificationDispatcher) *Room {
	return &Room{
		id:   id,
		name: name,

		players:        make(map[string]*player.Player),
		maxPlayerCount: 3,

		state: protodef.PRoomState_PROOM_STATE_WAITING,
	}
}

func (r *Room) GetId() uint32 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.id
}

func (r *Room) GetName() string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.name
}

func (r *Room) GetOwnerId() string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.ownerId
}

func (r *Room) GetState() protodef.PRoomState {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.state
}

func (r *Room) GetMaxPlayerCount() uint32 {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.maxPlayerCount
}

func (r *Room) IsOwnedBy(playerId string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.ownerId == playerId
}

func (r *Room) GetPlayers() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	playerIds := make([]string, 0, len(r.players))
	for playerId := range r.players {
		playerIds = append(playerIds, playerId)
	}
	return playerIds
}

func (r *Room) AddPlayer(player *player.Player, isOwner bool) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	playerId := player.GetPlayerId()
	if _, exists := r.players[playerId]; exists {
		return errordef.NewGameplayError(errordef.CodePlayerInRoom)
	}
	if len(r.players) >= int(r.maxPlayerCount) {
		return errordef.NewGameplayError(errordef.CodeRoomFull)
	}
	if isOwner {
		r.ownerId = playerId
	}
	player.OnEnterRoom(r.id)
	r.players[playerId] = player
	return nil
}

func (r *Room) RemovePlayer(playerId string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	// 从房间中移除玩家
	if player, exists := r.players[playerId]; exists {
		player.OnLeaveRoom()
		delete(r.players, playerId)
	}
}

func (r *Room) Destroy() {
	for _, player := range r.players {
		player.OnLeaveRoom()
	}
}
