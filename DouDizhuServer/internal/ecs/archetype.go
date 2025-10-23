package ecs

import (
	"encoding/json"
	"sort"
)

type Archetype struct {
	id         string
	components []ComponentTypeId

	caches map[int32]*Entity
}

func CreateArchetypeId(types ...ComponentTypeId) string {
	sort.Sort(ComponentTypeIdSlice(types))
	id, _ := json.Marshal(types)
	return string(id)
}

func NewArchetype(types ...ComponentTypeId) *Archetype {
	if types == nil {
		return nil
	}
	return &Archetype{
		id:         CreateArchetypeId(types...),
		components: types,
		caches:     make(map[int32]*Entity),
	}
}

func (a *Archetype) GetId() string {
	return a.id
}

func (a *Archetype) GetComponentsType() []ComponentTypeId {
	return a.components
}
