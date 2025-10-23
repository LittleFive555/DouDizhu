package room

import (
	"DouDizhuServer/scripts/errordef"
	"DouDizhuServer/scripts/network/message"
	"sync"

	"github.com/google/uuid"
)

var Manager *RoomManager

type RoomManager struct {
	rooms      map[uint32]*Room
	listOpLock sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[uint32]*Room),
	}
}

func (m *RoomManager) CreateRoom(name string, dispatcher message.INotificationDispatcher) *Room {
	m.listOpLock.Lock()
	defer m.listOpLock.Unlock()

	roomId := uuid.New().ID()
	room := NewRoom(roomId, name, dispatcher)
	m.rooms[roomId] = room
	return room
}

func (m *RoomManager) GetRoom(roomId uint32) (*Room, error) {
	m.listOpLock.RLock()
	defer m.listOpLock.RUnlock()

	room, ok := m.rooms[roomId]
	if !ok {
		return nil, errordef.NewGameplayError(errordef.CodeRoomNotExists)
	}
	return room, nil
}

func (m *RoomManager) GetRoomList() []*Room {
	m.listOpLock.RLock()
	defer m.listOpLock.RUnlock()

	rooms := make([]*Room, 0, len(m.rooms))
	for _, room := range m.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func (m *RoomManager) RemoveRoom(roomId uint32) error {
	m.listOpLock.Lock()
	defer m.listOpLock.Unlock()

	var room *Room
	var ok bool
	if room, ok = m.rooms[roomId]; !ok {
		return errordef.NewGameplayError(errordef.CodeRoomNotExists)
	}
	room.Destroy()
	delete(m.rooms, roomId)
	return nil
}
