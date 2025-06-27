package player

type AllPlayerNotificationGroup struct {
}

func NewAllPlayerNotificationGroup() *AllPlayerNotificationGroup {
	return &AllPlayerNotificationGroup{}
}

func (g *AllPlayerNotificationGroup) GetTargetPlayerIds() []string {
	playerIds := make([]string, 0)
	for _, player := range Manager.players {
		playerIds = append(playerIds, player.GetPlayerId())
	}
	return playerIds
}
