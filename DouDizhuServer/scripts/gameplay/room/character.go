package room

import (
	"DouDizhuServer/scripts/gameplay/component"
	"DouDizhuServer/scripts/network/protodef"

	"github.com/ungerik/go3d/vec3"
)

type RoomCharacter struct {
	id        string
	transform *component.Transform
	move      vec3.T

	inputBuffer []*protodef.PCharacterMove
	stateBuffer []*protodef.PCharacterState
}

func NewRoomCharacter(id string) *RoomCharacter {
	return &RoomCharacter{
		id:          id,
		transform:   component.NewTransform(),
		inputBuffer: make([]*protodef.PCharacterMove, 0),
		stateBuffer: make([]*protodef.PCharacterState, 0),
	}
}

func (c *RoomCharacter) GetId() string {
	return c.id
}

func (c *RoomCharacter) GetTransform() *component.Transform {
	return c.transform
}

func (c *RoomCharacter) GetMove() vec3.T {
	return c.move
}

func (c *RoomCharacter) SetMove(move vec3.T) {
	c.move = move
}

func (c *RoomCharacter) HasInput() bool {
	return len(c.inputBuffer) > 0
}

func (c *RoomCharacter) EnqueueInput(input *protodef.PCharacterMove) {
	c.inputBuffer = append(c.inputBuffer, input)
}

func (c *RoomCharacter) DequeueInput(untilTimestamp int64) []*protodef.PCharacterMove {
	if !c.HasInput() {
		return nil
	}
	inputs := make([]*protodef.PCharacterMove, 0)
	for _, input := range c.inputBuffer {
		if input.GetTimestamp() <= untilTimestamp {
			inputs = append(inputs, input)
		}
	}
	c.inputBuffer = c.inputBuffer[len(inputs):]
	return inputs
}

func (c *RoomCharacter) HasStateChange() bool {
	return len(c.stateBuffer) > 0
}

func (c *RoomCharacter) EnqueueState(state *protodef.PCharacterState) {
	c.stateBuffer = append(c.stateBuffer, state)
}

func (c *RoomCharacter) PopAllStateChange() []*protodef.PCharacterState {
	if !c.HasStateChange() {
		return nil
	}
	states := make([]*protodef.PCharacterState, len(c.stateBuffer))
	copy(states, c.stateBuffer)
	clear(c.stateBuffer)
	return states
}
