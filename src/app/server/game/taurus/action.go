package taurus

//游戏行为
type GameAction struct {
}

func newAction(seat *Seat) *GameAction {
	this := new(GameAction)
	this.init(seat)
	return this
}

func (this *GameAction) init(seat *Seat) {

}
