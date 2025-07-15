package room

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/network/protodef"
	"sync"
)

type Room struct {
	id             uint32
	name           string
	ownerId        string
	playerIds      []string
	maxPlayerCount uint32
	state          protodef.PRoomState
	lock           sync.RWMutex
}

func NewRoom(id uint32, name string) *Room {
	return &Room{
		id:             id,
		name:           name,
		playerIds:      []string{},
		maxPlayerCount: 3,
		state:          protodef.PRoomState_PROOM_STATE_WAITING,
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
	return r.playerIds
}

func (r *Room) SetOwner(playerId string) error {
	if err := r.AddPlayer(playerId); err != nil {
		return err
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.ownerId = playerId
	return nil
}

func (r *Room) AddPlayer(playerId string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if playerId == r.ownerId {
		return errordef.NewGameplayError(errordef.CodePlayerInRoom)
	}
	if len(r.playerIds) >= int(r.maxPlayerCount) {
		return errordef.NewGameplayError(errordef.CodeRoomFull)
	}
	r.playerIds = append(r.playerIds, playerId)
	return nil
}

func (r *Room) RemovePlayer(playerId string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i, p := range r.playerIds {
		if p == playerId {
			r.playerIds = append(r.playerIds[:i], r.playerIds[i+1:]...)
			return nil
		}
	}
	return errordef.NewGameplayError(errordef.CodePlayerNotInRoom)
}
