package ecs

type Entity struct {
	id         int32
	components map[ComponentTypeId]*Component
}

func NewEntity(id int32) *Entity {
	return &Entity{
		id:         id,
		components: make(map[ComponentTypeId]*Component),
	}
}

func (e *Entity) GetId() int32 {
	return e.id
}

func (e *Entity) HasComponent(id ComponentTypeId) bool {
	_, exists := e.components[id]
	return exists
}

func (e *Entity) GetComponent(id ComponentTypeId) *Component {
	component, exists := e.components[id]
	if !exists {
		return nil
	}
	return component
}

func (e *Entity) Is(archetype *Archetype) bool {
	for _, k := range archetype.GetComponentsType() {
		if _, exists := e.components[k]; !exists {
			return false
		}
	}

	return true
}
