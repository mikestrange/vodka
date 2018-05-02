package texas

const (
	//座位
	STATE_STAND   = 0
	STATE_SIT     = 1
	STATE_HOSTING = 2
	STATE_OFFLINE = 3 //下局离开游戏
)

//0 - max
type Seat struct {
	uid        int
	seat_id    int
	seat_state int
	init_chip  int //每回合带入的金币
	auto_buy   bool
	auto_money int
	new_sit    bool
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
	this.seat_state = STATE_STAND
}

func (this *Seat) update(money int, auto bool) {
	this.auto_money = money
	this.auto_buy = auto
}

func (this *Seat) sit(uid int, money int, auto bool) {
	this.uid = uid
	this.init_chip = money
	this.auto_money = money
	this.auto_buy = auto
	this.seat_state = STATE_SIT
	this.new_sit = true
}

func (this *Seat) stand() {
	this.uid = 0
	this.seat_state = STATE_STAND
}

//钱不够站起
func (this *Seat) stand_check(chip int) bool {
	if chip >= this.init_chip {
		return true
	}
	return false
}

func (this *Seat) issit() bool {
	return this.seat_state > STATE_STAND
}

//新坐下的需要交大盲
func (this *Seat) checkNewsit() bool {
	if this.new_sit {
		this.new_sit = false
		return true
	}
	return false
}

//游戏相关
func (this *Seat) checkAction() (*GameAction, bool) {
	if this.game_action == nil {
		return nil, false
	}
	return this.game_action, true
}

//直接
func (this *Seat) Action() *GameAction {
	return this.game_action
}

func (this *Seat) isPlaying() bool {
	return this.game_action != nil
}

//开始
func (this *Seat) beginGame() bool {
	if this.issit() {
		if this.game_action == nil {
			this.game_action = newAction(this)
			return true
		}
	}
	return false
}

func (this *Seat) endGame() {
	if this.game_action != nil {
		if this.new_sit {
			//离开后重新坐下了的人
		} else {
			//没有离开
			this.init_chip = this.game_action.result()
		}
		this.game_action = nil
	}
}
