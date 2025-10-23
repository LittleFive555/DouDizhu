package ecs

type ISystem interface {
	IsSystem()
	TypeId() SystemTypeId
}

type SystemTypeId int32
