package room

import (
	"DouDizhuServer/scripts/network/protodef"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ungerik/go3d/vec3"
)

type RoomWorld struct {
	stop bool
	lock sync.RWMutex

	worldId             string
	frameTime           int64
	lastUpdateTimestamp int64
	characters          map[string]*RoomCharacter
}

func NewRoomWorld(worldId string, frameRate int) *RoomWorld {
	return &RoomWorld{
		worldId:             worldId,
		frameTime:           1000 / int64(frameRate),
		lastUpdateTimestamp: time.Now().UnixMilli(),
		characters:          make(map[string]*RoomCharacter),

		lock: sync.RWMutex{},
	}
}

func (world *RoomWorld) RunLoop() {
	world.stop = false
	for {
		if world.stop {
			return
		}
		now := time.Now().UnixMilli()
		delta := now - world.lastUpdateTimestamp
		for delta >= world.frameTime {
			delta -= world.frameTime
			lastTimestamp := world.lastUpdateTimestamp
			world.lastUpdateTimestamp += world.frameTime
			world.update(lastTimestamp, world.lastUpdateTimestamp)
		}
		time.Sleep(time.Millisecond)
	}
}

func (world *RoomWorld) Stop() {
	world.stop = true
}

func (world *RoomWorld) AddCharacter() string {
	world.lock.Lock()
	defer world.lock.Unlock()

	characterId := uuid.New().String()
	world.characters[characterId] = NewRoomCharacter(characterId)
	return characterId
}

func (world *RoomWorld) RemoveCharacter(characterId string) {
	world.lock.Lock()
	defer world.lock.Unlock()

	delete(world.characters, characterId)
}

func (world *RoomWorld) update(lastTimestamp int64, nowTimestamp int64) {
	world.lock.RLock()
	defer world.lock.RUnlock()

	// 依次处理每个角色的输入
	for _, character := range world.characters {
		// 如果有输入，则处理输入
		if character.HasInput() {
			inputs := character.DequeueInput(nowTimestamp)
			move := character.GetMove()
			// 如果有已存在的输入状态，则处理该输入直到新的输入事件
			if move.LengthSqr() > 0 {
				startTime := lastTimestamp
				endTime := inputs[0].GetTimestamp()
				duration := endTime - startTime

				durationSecond := float32(duration) / 1000
				move.Scale(durationSecond)

				transform := character.GetTransform()
				transform.Position.Add(&move)
				character.EnqueueState(&protodef.PCharacterState{
					Id:        character.GetId(),
					Timestamp: endTime,
					Pos:       Vec3ToProto(transform.Position),
				})
			}

			// 新收到的输入
			for i, input := range inputs {
				startTime := input.GetTimestamp()
				var endTime int64
				// 每个输入的持续时间为两次输入的间隔，如果没有后续输入，则结束时间为updateTimestamp
				if i+1 < len(inputs) {
					endTime = inputs[i+1].GetTimestamp()
				} else {
					endTime = nowTimestamp
				}
				duration := endTime - startTime

				transform := character.GetTransform()
				move := Vec3FromProto(input.GetMove())
				character.SetMove(move)
				if move.LengthSqr() > 0 {
					durationSecond := float32(duration) / 1000
					move.Scale(durationSecond)

					transform.Position.Add(&move)
					character.EnqueueState(&protodef.PCharacterState{
						Id:        character.GetId(),
						Timestamp: endTime,
						Pos:       Vec3ToProto(transform.Position),
					})
				}
			}
		} else {
			// 如果没有新的输入，则处理已存在的输入
			move := character.GetMove()
			if move.LengthSqr() > 0 {
				durationSecond := float32(world.frameTime) / 1000
				move.Scale(durationSecond)
				transform := character.GetTransform()
				transform.Position.Add(&move)
				character.EnqueueState(&protodef.PCharacterState{
					Id:        character.GetId(),
					Timestamp: nowTimestamp,
					Pos:       Vec3ToProto(transform.Position),
				})
			}
		}
	}
	// TODO 向客户端同步每个角色的状态
}

func Vec3ToProto(v vec3.T) *protodef.PVector3 {
	return &protodef.PVector3{
		X: v[0],
		Y: v[1],
		Z: v[2],
	}
}

func Vec3FromProto(v *protodef.PVector3) vec3.T {
	return vec3.T{
		v.GetX(),
		v.GetY(),
		v.GetZ(),
	}
}
