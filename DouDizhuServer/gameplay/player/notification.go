package player

type AllPlayerNotificationGroup struct {
}

func NewAllPlayerNotificationGroup() *AllPlayerNotificationGroup {
	return &AllPlayerNotificationGroup{}
}

func (g *AllPlayerNotificationGroup) GetTargetSessionIds() []string {
	sessionIds := make([]string, 0)
	for _, player := range Manager.players {
		sessionIds = append(sessionIds, player.GetSessionId())
	}
	return sessionIds
}
