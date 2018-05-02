package texas

//游戏状态
const (
	GAME_SEAT_WAIT   = 0
	GAME_SEAT_PLAYIN = 1
)

type GameAction struct {
	uid        int
	seat_id    int
	give_up    bool
	pot_point  int //参与池
	round_chip int //回合下注
	cards      []int16
	big_cards  []int16
	card_type  int
	seat_money int
	init_chip  int
}

func newAction(seat *Seat) *GameAction {
	this := new(GameAction)
	this.init(seat)
	return this
}

//game
func (this *GameAction) init(seat *Seat) {
	this.uid = seat.uid
	this.seat_id = seat.seat_id
	this.seat_money = seat.init_chip
	this.init_chip = seat.init_chip
	//
	this.round_chip = 0
	this.pot_point = 0
	this.give_up = false
	this.cards = make([]int16, 0, 2)
	this.big_cards = nil
	this.card_type = 0
}

//回合清理
func (this *GameAction) turnOver() {
	this.round_chip = 0
}

//能行动在游戏（没有弃牌和allin）
func (this *GameAction) isAction() bool {
	if this.isAllin() || this.isFold() {
		return false
	}
	return true
}

func (this *GameAction) isNoAction() bool {
	if this.isAllin() || this.isFold() {
		return true
	}
	return false
}

func (this *GameAction) Totals() int {
	return this.round_chip + this.seat_money
}

//allin的
func (this *GameAction) isAllin() bool {
	return this.seat_money == 0
}

func (this *GameAction) isNoAllin() bool {
	return this.seat_money > 0
}

//弃牌
func (this *GameAction) isFold() bool {
	return this.give_up
}

func (this *GameAction) isNoFold() bool {
	return !this.give_up
}

func (this *GameAction) setFold() {
	this.give_up = true
}

func (this *GameAction) checkSeat(idx int) bool {
	return this.seat_id == idx
}

//cards
func (this *GameAction) pushCard(num int16) {
	this.cards = append(this.cards, num)
}

func (this *GameAction) cardFlush(cards []int16) {
	list := make([]int16, 7)
	copy(list, this.cards)
	copy(list[2:], cards)
	this.card_type, this.big_cards = CardTypeOfTexas(list)
}

func (this *GameAction) compare(action *GameAction) int {
	if this.card_type > action.card_type {
		return WIN
	} else if action.card_type > this.card_type {
		return LOSE
	}
	ret := CompareCards(this.big_cards, action.big_cards, 5)
	if ret == 1 {
		return WIN
	} else if ret == -1 {
		return LOSE
	}
	return DRAW
}

func (this *GameAction) result() int {
	return this.seat_money
}

//返回真实下注的筹码
func (this *GameAction) BetChip(val int) int {
	//加注的值(扣)
	chip := val - this.round_chip
	//allin
	if chip >= this.seat_money {
		chip = this.seat_money
	}
	this.seat_money -= chip
	this.round_chip = this.round_chip + chip
	return this.round_chip
}

//结果赢得
func (this *GameAction) Result(val int) {
	this.seat_money += val
}

//服务费
func (this *GameAction) SubFree(val int) {
	if val >= this.seat_money {
		panic("系统错误, 服务费不足!")
		return
	}
	this.seat_money -= val
}
