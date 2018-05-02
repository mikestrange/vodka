package chat

type GameUser struct {
	Player *GameHeader
	//info(其他信息)
	UserName string
}

func (this *GameUser) SerID() int {
	return this.Player.GateID
}
