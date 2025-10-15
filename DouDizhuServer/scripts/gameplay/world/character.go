package world

import (
	"DouDizhuServer/scripts/gameplay/component"
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network/protodef"
	"sort"
	"sync"

	"github.com/ungerik/go3d/vec3"
)

type Character struct {
	id        string
	transform *component.Transform
	speed     float32
	move      vec3.T

	worldId string

	inputBufferLock sync.RWMutex
	inputBuffer     []*protodef.PCharacterMove

	stateBufferLock sync.RWMutex
	stateBuffer     []*protodef.PCharacterState
}

func NewCharacter(id string) *Character {
	return &Character{
		id:          id,
		speed:       2,
		transform:   component.NewTransform(),
		inputBuffer: make([]*protodef.PCharacterMove, 0),
		stateBuffer: make([]*protodef.PCharacterState, 0),
	}
}

func (c *Character) GetId() string {
	return c.id
}

func (c *Character) GetTransform() *component.Transform {
	return c.transform
}

func (c *Character) GetSpeed() float32 {
	return c.speed
}

func (c *Character) GetMove() vec3.T {
	return c.move
}

func (c *Character) SetMove(move vec3.T) {
	c.move = move
}

func (c *Character) HasInput() bool {
	c.inputBufferLock.RLock()
	defer c.inputBufferLock.RUnlock()
	return len(c.inputBuffer) > 0
}

func (c *Character) EnqueueInput(input *protodef.PCharacterMove) {
	c.inputBufferLock.Lock()
	defer c.inputBufferLock.Unlock()
	c.inputBuffer = append(c.inputBuffer, input)
}

func (c *Character) DequeueInput(fromTimestamp int64, untilTimestamp int64) []*protodef.PCharacterMove {
	c.inputBufferLock.Lock()
	defer c.inputBufferLock.Unlock()
	if len(c.inputBuffer) == 0 {
		return nil
	}
	inputs := make([]*protodef.PCharacterMove, 0)
	sort.Sort(CharacterInputByTimestamp(c.inputBuffer))
	dequeueCount := 0
	for _, input := range c.inputBuffer {
		inputTimestamp := input.GetTimestamp()
		if inputTimestamp < fromTimestamp {
			logger.WarnFormat("Drop input: inputTimestamp %d is before fromTimestamp %d", inputTimestamp, fromTimestamp)
			dequeueCount++
			continue
		}
		if inputTimestamp <= untilTimestamp {
			dequeueCount++
			inputs = append(inputs, input)
		}
	}
	c.inputBuffer = c.inputBuffer[dequeueCount:]
	return inputs
}

func (c *Character) HasStateChange() bool {
	c.stateBufferLock.RLock()
	defer c.stateBufferLock.RUnlock()
	return len(c.stateBuffer) > 0
}

func (c *Character) EnqueueState(state *protodef.PCharacterState) {
	c.stateBufferLock.Lock()
	defer c.stateBufferLock.Unlock()
	c.stateBuffer = append(c.stateBuffer, state)
}

func (c *Character) PopAllStateChange() []*protodef.PCharacterState {
	c.stateBufferLock.Lock()
	defer c.stateBufferLock.Unlock()
	if len(c.stateBuffer) == 0 {
		return nil
	}

	states := make([]*protodef.PCharacterState, len(c.stateBuffer))
	copy(states, c.stateBuffer)

	clear(c.stateBuffer)
	c.stateBuffer = c.stateBuffer[:0]

	return states
}

func (c *Character) GetFullState() *protodef.PCharacterState {
	return &protodef.PCharacterState{
		Pos: Vec3ToProto(c.transform.Position),
	}
}

type CharacterInputByTimestamp []*protodef.PCharacterMove

func (c CharacterInputByTimestamp) Len() int {
	return len(c)
}

func (c CharacterInputByTimestamp) Less(i, j int) bool {
	return c[i].Timestamp < c[j].Timestamp
}

func (c CharacterInputByTimestamp) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
