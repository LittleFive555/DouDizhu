package ecs

type World struct {
	idCounter  int32
	entities   map[int32]*Entity
	archetypes map[string]*Archetype
	systems    map[SystemTypeId]ISystem

	// value为true表示entity被删除了，为false表示entity组件有修改
	dirtyEntities map[int32]bool
}

func NewWorld() *World {
	return &World{
		idCounter:  0,
		entities:   make(map[int32]*Entity),
		archetypes: make(map[string]*Archetype),
		systems:    make(map[SystemTypeId]ISystem),

		dirtyEntities: make(map[int32]bool),
	}
}

func (w *World) AddEntity(components ...*Component) int32 {
	id := w.idCounter
	w.idCounter++

	// 生成新实体
	w.entities[id] = NewEntity(id)

	// 添加组件
	for _, component := range components {
		w.AddComponent(id, component)
	}
	// 标记修改
	w.markEntityDirty(id)
	return id
}

func (w *World) RemoveEntity(id int32) {
	if _, exists := w.entities[id]; !exists {
		return
	}

	// 标记删除
	w.markEntityDeleted(id)

	// 移除实体
	delete(w.entities, id)
}

func (w *World) GetEntity(id int32) *Entity {
	if entity, exists := w.entities[id]; exists {
		return entity
	}
	return nil
}

func (w *World) AddComponent(entityId int32, component *Component) {
	entitiy := w.entities[entityId]
	if entitiy == nil {
		return
	}
	if entitiy.HasComponent(component.typeId) {
		return
	}
	w.markEntityDirty(entityId)
	entitiy.components[component.typeId] = component
}

func (w *World) RemoveComponent(entityId int32, componentTypeId ComponentTypeId) {
	entitiy := w.entities[entityId]
	if entitiy == nil {
		return
	}
	if !entitiy.HasComponent(componentTypeId) {
		return
	}
	w.markEntityDirty(entityId)
	delete(entitiy.components, componentTypeId)
}

func (w *World) QueryEntities(types ...ComponentTypeId) map[int32]*Entity {
	// 原型有变化的实体需要重新分类
	w.resolveDirtyEntities()

	// 查找是否有缓存，有缓存则直接返回
	id := CreateArchetypeId(types...)
	if archetype, exists := w.archetypes[id]; exists {
		return archetype.caches
	}

	// 否则生成新原型
	archetype := NewArchetype(types...)
	for entityId, entity := range w.entities {
		if entity.Is(archetype) {
			archetype.caches[entityId] = entity
		}
	}
	return archetype.caches
}

func (w *World) AddSystem(system ISystem) {
	if _, exists := w.systems[system.TypeId()]; exists {
		return
	}
	w.systems[system.TypeId()] = system
}

func (w *World) RemoveSystem(id SystemTypeId) {
	if _, exists := w.systems[id]; !exists {
		return
	}
	delete(w.systems, id)
}

func (w *World) markEntityDeleted(entityId int32) {
	w.dirtyEntities[entityId] = true
}

func (w *World) markEntityDirty(entityId int32) {
	if _, exists := w.dirtyEntities[entityId]; exists {
		return
	}
	w.dirtyEntities[entityId] = false
}

func (w *World) resolveDirtyEntities() {
	for _, archetype := range w.archetypes {
		for entityId, isDeleted := range w.dirtyEntities {
			if isDeleted {
				delete(archetype.caches, entityId)
			} else {
				if entity, exists := w.entities[entityId]; exists {
					if entity.Is(archetype) {
						archetype.caches[entityId] = entity
					} else {
						delete(archetype.caches, entityId)
					}
				} else {
					delete(archetype.caches, entityId)
				}
			}
		}
	}
	clear(w.dirtyEntities)
}
