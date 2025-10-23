package test

import (
	"DouDizhuServer/internal/ecs"
	"DouDizhuServer/internal/ecs/test/component"
	"DouDizhuServer/internal/ecs/test/system"

	"github.com/ungerik/go3d/vec3"
)

func RunMoveSystem(runTimes int, position *vec3.T, velocity *vec3.T) *vec3.T {
	world := ecs.NewWorld()
	entityId := world.AddEntity(
		ecs.NewComponent(component.COMPONENT_TYPE_POSITION, &component.Position{Value: position}),
		ecs.NewComponent(component.COMPONENT_TYPE_VELOCITY, &component.Velocity{Value: velocity}),
	)
	movementSystem := system.MovementSystem{}
	for range runTimes {
		movementSystem.Execute(world)
	}

	return world.GetEntity(entityId).GetComponent(component.COMPONENT_TYPE_POSITION).Data.(*component.Position).Value
}
