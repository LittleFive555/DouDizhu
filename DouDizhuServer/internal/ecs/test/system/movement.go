package system

import (
	"DouDizhuServer/internal/ecs"
	"DouDizhuServer/internal/ecs/test/component"
)

type MovementSystem struct {
}

func (s *MovementSystem) IsSystem() {}

func (s *MovementSystem) TypeId() ecs.SystemTypeId {
	return SYSTEM_TYPE_MOVEMENT
}

func (s *MovementSystem) NeedComponents() []ecs.ComponentTypeId {
	return []ecs.ComponentTypeId{component.COMPONENT_TYPE_POSITION, component.COMPONENT_TYPE_VELOCITY}
}

func (s *MovementSystem) Execute(world *ecs.World) {
	entities := world.QueryEntities(component.COMPONENT_TYPE_POSITION, component.COMPONENT_TYPE_VELOCITY)
	for _, entity := range entities {
		position := entity.GetComponent(component.COMPONENT_TYPE_POSITION)
		velocity := entity.GetComponent(component.COMPONENT_TYPE_VELOCITY)
		position.Data.(*component.Position).Value.Add(velocity.Data.(*component.Velocity).Value)
	}
}
