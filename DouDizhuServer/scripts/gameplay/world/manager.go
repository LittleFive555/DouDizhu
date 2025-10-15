package world

import (
	"DouDizhuServer/scripts/logger"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/protodef"
)

type WorldManager struct {
	worlds     map[string]*World
	characters map[string]*Character
	dispatcher message.INotificationDispatcher
}

func NewWorldManager(dispatcher message.INotificationDispatcher) *WorldManager {
	manager := &WorldManager{
		worlds:     make(map[string]*World),
		characters: make(map[string]*Character),
		dispatcher: dispatcher,
	}
	manager.AddWorld("default") // TODO 用配置
	return manager
}

func (wm *WorldManager) AddWorld(worldId string) {
	if _, exists := wm.worlds[worldId]; exists {
		logger.ErrorWith("worldId已存在", "worldId", worldId)
		return
	}
	world := newWorld(worldId, 20, 3)
	wm.worlds[worldId] = world
	logger.InfoWith("创建世界成功", "worldId", worldId)
}

func (wm *WorldManager) GetWorldFullState(worldId string) *protodef.PWorldState {
	world := wm.worlds[worldId]
	return world.GetFullWorldState()
}

func (wm *WorldManager) CharacterEnterWorld(characterId string, worldId string) {
	var character *Character
	var exists bool
	if character, exists = wm.characters[characterId]; !exists {
		character = NewCharacter(characterId)
		wm.characters[characterId] = character
	}
	// TODO verify worldId
	if character.worldId != "" && character.worldId != worldId {
		oldWorld := wm.worlds[character.worldId]
		oldWorld.OnCharacterLeave(character)
	}
	character.worldId = worldId
	world := wm.worlds[character.worldId]
	world.OnCharacterEnter(character)
}

func (wm *WorldManager) CharacterInput(characterId string, input *protodef.PCharacterMove) {
	character := wm.characters[characterId]
	character.EnqueueInput(input)
}

func (wm *WorldManager) RemoveCharacter(characterId string) {
	character := wm.characters[characterId]
	delete(wm.characters, characterId)

	if character.worldId != "" {
		world := wm.worlds[character.worldId]
		world.OnCharacterLeave(character)
	}
}

func (wm *WorldManager) RunLoop() {
	for _, world := range wm.worlds {
		go world.RunLoop(wm.dispatcher)
	}
}

func (wm *WorldManager) Stop() {
	for _, world := range wm.worlds {
		world.Stop()
	}
}
