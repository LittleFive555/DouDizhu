package player

import "DouDizhuServer/scripts/network/protodef"

type Player struct {
	playerId string
	nickname string
	roomId   uint32
	state    protodef.PPlayerState

	sessionId string
}

func NewPlayer(playerId string, nickname string, sessionId string) *Player {
	return &Player{
		playerId:  playerId,
		nickname:  nickname,
		state:     protodef.PPlayerState_PPLAYER_STATE_LOBBY,
		sessionId: sessionId,
	}
}

func (p *Player) GetPlayerId() string {
	return p.playerId
}

func (p *Player) GetNickname() string {
	return p.nickname
}

func (p *Player) SetNickname(nickname string) {
	p.nickname = nickname
}

func (p *Player) GetRoomId() uint32 {
	return p.roomId
}

func (p *Player) GetState() protodef.PPlayerState {
	return p.state
}

func (p *Player) GetSessionId() string {
	return p.sessionId
}

func (p *Player) IsInLobby() bool {
	return p.state == protodef.PPlayerState_PPLAYER_STATE_LOBBY
}

func (p *Player) IsOnline() bool {
	return p.state != protodef.PPlayerState_PPLAYER_STATE_NONE
}

func (p *Player) IsInRoom() bool {
	return p.state == protodef.PPlayerState_PPLAYER_STATE_IN_ROOM
}

func (p *Player) IsInGame() bool {
	return p.state == protodef.PPlayerState_PPLAYER_STATE_IN_GAME
}

func (p *Player) Login() {
	p.state = protodef.PPlayerState_PPLAYER_STATE_LOBBY
}

func (p *Player) Logout() {
	p.state = protodef.PPlayerState_PPLAYER_STATE_NONE
}

func (p *Player) EnterRoom(roomId uint32) {
	p.roomId = roomId
	p.state = protodef.PPlayerState_PPLAYER_STATE_IN_ROOM
}

func (p *Player) LeaveRoom() {
	p.roomId = 0
	p.state = protodef.PPlayerState_PPLAYER_STATE_LOBBY
}

func (p *Player) StartGame() {
	p.state = protodef.PPlayerState_PPLAYER_STATE_IN_GAME
}
