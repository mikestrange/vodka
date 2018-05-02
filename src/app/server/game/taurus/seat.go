package taurus

const (
	//
	GAME_SEAT_WAIT       = 0
	GAME_SEAT_PLAYIN     = 1
	GAME_SEAT_ROB_BANKER = 2 //抢庄
	GAME_SEAT_BET_CHIP   = 3 //下倍
	GAME_SEAT_COMMIT     = 4 //提交
	//
	STATE_STAND   = 0
	STATE_SIT     = 1
	STATE_HOSTING = 2
	STATE_OFFLINE = 3 //下局离开游戏
)

//不支持直接退出游戏
type Seat struct {
	seat_id      int
	uid          int
	gift         int
	name         string
	seat_money   int     //坐下金币
	rob_multiple int     //抢庄的倍数
	bet_multiple int     //下注的倍数
	game_status  int     //0等待游戏，1正在游戏
	seat_status  int     //0站起，1坐下，2托管，3离线
	auto_buy     bool    //是否自动买入
	auto_money   int     //自动买入时的金币
	cards        []int16 //5张牌
	card_type    int     //牌类型
	result       int     //输赢和
	result_money int     //金币
	player       *Player
}

func newSeat(idx int) *Seat {
	this := new(Seat)
	this.init(idx)
	return this
}

func (this *Seat) init(idx int) {
	this.cards = make([]int16, 0, 5)
	this.seat_id = idx
}

func (this *Seat) sit(uid int, money int, auto_buy bool, player *Player) {
	this.uid = uid
	this.seat_money = money
	this.auto_money = money
	this.auto_buy = auto_buy
	this.seat_status = STATE_SIT
	this.player = player
}

//多次坐下刷新
func (this *Seat) update(money int, auto_buy bool) {
	this.auto_buy = auto_buy
	this.auto_money = money
}

func (this *Seat) stand() {
	this.uid = 0
	this.player = nil
	this.game_status = GAME_SEAT_WAIT
	this.seat_status = STATE_STAND
}

func (this *Seat) over_stand() {
	this.seat_status = STATE_OFFLINE
}

func (this *Seat) issit() bool {
	return this.seat_status > STATE_STAND
}

//game state
func (this *Seat) isplayer() bool {
	return this.issit() && this.game_status > GAME_SEAT_WAIT
}

func (this *Seat) begin() {
	this.game_status = GAME_SEAT_PLAYIN
	this.rob_multiple = 0
	this.bet_multiple = 0
	this.card_type = 0
	this.result = -1
	this.result_money = 0
	this.clear_cards()
}

//抢庄
func (this *Seat) rob_banker(num int) {
	this.rob_multiple = num
	this.game_status = GAME_SEAT_ROB_BANKER
}

func (this *Seat) bet_chip(min int) {
	if this.bet_multiple < min {
		this.bet_multiple = min
	}
	this.game_status = GAME_SEAT_BET_CHIP
}

func (this *Seat) commit() {
	this.game_status = GAME_SEAT_COMMIT
}

func (this *Seat) iscommit() bool {
	return this.game_status == GAME_SEAT_COMMIT
}

func (this *Seat) over() {
	this.game_status = GAME_SEAT_WAIT
}

func (this *Seat) stand_check(chip int) bool {
	if this.seat_status == STATE_OFFLINE {
		return true
	}
	//钱不够
	if chip > this.seat_money {
		return true
	}
	return false
}

func (this *Seat) pushCard(num int16) {
	this.cards = append(this.cards, num)
}

func (this *Seat) clear_cards() {
	if len(this.cards) > 0 {
		this.cards = make([]int16, 0, 5)
	}
}

func (this *Seat) card_flush() {
	this.card_type = getNiuNiuType(this.cards)
}

func (this *Seat) card_str() string {
	return TypeStr(this.card_type) // + ">" + ToStr(this.cards...)
}

//牌型倍数
func (this *Seat) card_multiple() int {
	if this.card_type > TYPE_BEINIU && this.card_type < TYPE_NIUNIU {
		return 2
	}
	switch this.card_type {
	case TYPE_NIUNIU:
		return 3
	case TYPE_YINNIU:
		return 4
	case TYPE_JINNIU:
		return 5
	case TYPE_BOMB:
		return 6
	}
	return 1
}

//0 赢 1 输 2 和 (炸弹需要特许判断)
func (this *Seat) check_result(seat *Seat) int {
	if this.card_type > seat.card_type {
		return WIN
	}
	if this.card_type < seat.card_type {
		return LOSE
	}
	//炸弹(一副牌不会存在相同的情况, 所以判断中间那张牌就可以了)
	if this.card_type == TYPE_BOMB && seat.card_type == TYPE_BOMB {
		return CompareCard(this.cards[2], seat.cards[2])
	}
	atype := DRAW
	for i := 0; i < 5; i++ {
		ret := CompareCard(this.cards[i], seat.cards[i])
		if ret == DRAW {
			continue
		}
		atype = ret
		break
	}
	return atype
}

func (this *Seat) result_set(val int) {
	this.result = val
}

//减钱
func (this *Seat) sub_money(val int) int {
	if val > this.seat_money {
		val = this.seat_money
		this.seat_money = 0
	} else {
		this.seat_money = this.seat_money - val
	}
	this.result_money -= val
	return val
}

//加钱
func (this *Seat) add_money(val int) {
	this.seat_money = this.seat_money + val
	this.result_money += val
}
