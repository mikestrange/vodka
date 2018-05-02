package landlord

import (
	"ants/base"
)

//--加注差额应当至少等同于当前下注轮中之前最大的下注或加注差额。
type GameLogic struct {
	BaseGame
	cards *CardDealer
}

func NewGame(tid int, data interface{}) *GameLogic {
	this := new(GameLogic)
	this.init(tid, data)
	return this
}

func (this *GameLogic) init(tid int, data interface{}) {
	this.table_id = tid
	this.min_player_num = 3
	this.seat_count = 3
	//
	this.cards = NewDealer()
	this.BaseGame.init()
}

func (this *GameLogic) start() {
	if this.check_state(GAME_WAIT) {
		return
	}
	player_num := this.get_sit_num()
	if player_num < this.min_player_num {
		return
	}
	this.SetState(GAME_BANKER)
	//洗牌确定发牌位置已经叫庄位置
	this.game_multiple = 1
	this.cards.Random()
	this.banker_idx = base.Random(this.seat_count) //0-2
	this.attack_idx = 0                            //
	this.call_num = 0
	//发手牌
	for i := 0; i < 15; i++ {
		for j := 0; j < this.seat_count; j++ {
			idx := (this.attack_idx + j) % this.seat_count
			seat := this.seats[idx]
			seat.Begin()
			seat.Action().pushCard(this.cards.Pop())
		}
	}
	//公共牌(已经确定)
	this.public_cards = []int16{}
	for i := 0; i < 3; i++ {
		this.public_cards = append(this.public_cards, this.cards.Pop())
	}
}

//叫庄
func (this *GameLogic) betBanker(call bool) {
	if this.call_num == 4 {
		return
	}
	seat := this.seats[this.banker_idx]
	this.call_num++
	//叫地主
	if call {
		seat.Action().Call()
		this.attack_idx = seat.seat_id
	}
	//移动到下一个
	this.banker_idx++
	if this.banker_idx == this.seat_count {
		this.banker_idx = 0
	}
	//回到自己
	if this.call_num == 4 {
		this.callOver()
	}
}

//加倍
func (this *GameLogic) doubleKill(seat *Seat) {
	seat.Action().Double()
}

//结束
func (this *GameLogic) callOver() {
	if this.attack_idx == 0 {
		//1,4条2或2王的接手
		//2,重新开始
		this.clear()
		this.SetState(GAME_WAIT)
		this.start()
	} else {
		this.banker_idx = this.attack_idx
		this.current_idx = 0
		this.startOper(this.attack_idx)
	}
}

//提醒用户开始操作
func (this *GameLogic) startOper(idx int) {
	if idx == this.attack_idx {
		//首发重新开始
	} else {
		if idx > this.seat_count {
			this.current_idx = 1
		} else {
			this.current_idx = idx
		}
		if seat, ok := this.get_seat(this.current_idx); ok {
			//提醒
			println("action = ", seat.seat_id)
		} else {
			panic("获取不到当前位置")
		}
	}
}

//类型，和出牌的位置 0表示放弃
func (this *GameLogic) outCards(seat *Seat, action int, ctype int, cards []int16) {
	//获取攻击位置
	attack, ok := this.get_seat(this.attack_idx)
	if !ok {
		panic("获取不到攻击位置")
	}
	if seat.seat_id == attack.seat_id { //攻击位置操作
		if attack.Action().outCards(ctype, cards) {
			if attack.Action().overShow() {
				this.over(seat)
			} else {
				this.startOper(this.current_idx + 1)
			}
		} else {
			//不合法的出法
		}
	} else { //跟随位置操作
		if action == 0 {
			//过
			this.startOper(this.current_idx + 1)
		} else if seat.Action().BigTo(attack.Action(), ctype, cards) {
			if seat.Action().overShow() {
				this.over(seat)
			} else {
				this.attack_idx = seat.seat_id
				this.startOper(this.current_idx + 1)
			}
		}
	}
}

func (this *GameLogic) over(seat *Seat) {

}

func (this *GameLogic) clear() {
	this.each_seats(func(seat *Seat) {
		seat.End()
	})
}

func init() {
	m := NewGame(1, nil)
	m.start()
	println("game start")
}
