package landlord

const (
	//座位
	STATE_STAND   = 0
	STATE_SIT     = 1
	STATE_HOSTING = 2
	STATE_OFFLINE = 3 //下局离开游戏
)

type Seat struct {
	uid       int
	seat_id   int
	state     int
	init_chip int
	//游戏行动
	game_action *GameAction
}

func newSeat(idx int) *Seat {
	this := new(Seat)
	this.init(idx)
	return this
}

func (this *Seat) init(idx int) {
	this.seat_id = idx
	this.state = STATE_STAND
}

func (this *Seat) update() {

}

func (this *Seat) sit(uid int) {
	this.uid = uid
	this.state = STATE_SIT
}

func (this *Seat) stand() {
	this.uid = 0
	this.state = STATE_STAND
}

//钱不够站起
func (this *Seat) stand_check(chip int) bool {
	if chip >= this.init_chip {
		return true
	}
	return false
}

func (this *Seat) issit() bool {
	return this.state > STATE_STAND
}

//直接
func (this *Seat) Action() *GameAction {
	return this.game_action
}

func (this *Seat) Begin() {
	this.game_action = newAction(this)
}

func (this *Seat) End() {
	this.game_action = nil
}
