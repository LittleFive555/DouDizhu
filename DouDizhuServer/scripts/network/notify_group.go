package network

import (
	"DouDizhuServer/scripts/gameplay/player"
	"DouDizhuServer/scripts/gameplay/room"
	"DouDizhuServer/scripts/network/message"
	"DouDizhuServer/scripts/network/session"
)

func getNotifySessionIds(sessionMgr *session.SessionManager, notifyGroup message.INotifyGroup) []string {
	switch notifyGroup := notifyGroup.(type) {
	case *room.RoomNotifyGroup:
		return getRoomNotifySessionIds(notifyGroup)
	case *player.AllPlayerNotificationGroup:
		return getAllPlayerTargetSessionIds()
	case *message.AllConnectionNotificationGroup:
		return getAllConnectionTargetSessionIds(sessionMgr)
	}
	return nil
}

func getAllConnectionTargetSessionIds(sessionMgr *session.SessionManager) []string {
	sessionIds := make([]string, 0)
	for _, session := range sessionMgr.GetAllSessions() {
		sessionIds = append(sessionIds, session.Id)
	}
	return sessionIds
}

func getAllPlayerTargetSessionIds() []string {
	sessionIds := make([]string, 0)
	for _, player := range player.Manager.GetAllPlayers() {
		sessionIds = append(sessionIds, player.GetSessionId())
	}
	return sessionIds
}

func getRoomNotifySessionIds(roomNotifyGroup *room.RoomNotifyGroup) []string {
	room, err := room.Manager.GetRoom(roomNotifyGroup.RoomId)
	if err != nil {
		return nil
	}
	playerIds := room.GetPlayers()
	sessionIds := make([]string, 0)
	for _, playerId := range playerIds {
		if playerId == roomNotifyGroup.ExceptPlayerId {
			continue
		}
		sessionIds = append(sessionIds, player.Manager.GetPlayer(playerId).GetSessionId())
	}

	return sessionIds
}
