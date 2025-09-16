package room

import (
	"DouDizhuServer/scripts/network/protodef"
	"sync"
	"time"
)

type RoomWorld struct {
	characters map[string]*RoomCharacter
	inputs     map[string]*CharacterInput
	stop       chan struct{}
	lock       sync.RWMutex
}

type CharacterInput struct {
	playerId  string
	move      *protodef.PCharacterMove
	timestamp int64
}

func NewRoomWorld() *RoomWorld {
	return &RoomWorld{
		characters: make(map[string]*RoomCharacter),
		inputs:     make(map[string]*CharacterInput),
		lock:       sync.RWMutex{},
	}
}

func (world *RoomWorld) RunLoop() {
	world.stop = make(chan struct{})
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-world.stop:
			ticker.Stop()
			return
		default:
		}
		select {
		case <-ticker.C:
			world.Update()
		default:
		}
	}
}

func (world *RoomWorld) Stop() {
	world.stop <- struct{}{}
}

func (world *RoomWorld) Update() {
	world.lock.RLock()
	defer world.lock.RUnlock()

}

func (world *RoomWorld) ChangeInput(input *CharacterInput) {

}

func (world *RoomWorld) AddCharacter(playerId string) {
	world.lock.Lock()
	defer world.lock.Unlock()
	world.characters[playerId] = &RoomCharacter{
		id: playerId,
		position: &protodef.PVector3{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}
}

func (world *RoomWorld) RemoveCharacter(playerId string) {
	world.lock.Lock()
	defer world.lock.Unlock()
	delete(world.characters, playerId)
}
