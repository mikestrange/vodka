package texas

import "sort"
import "math"

type PotPondList []*PotPond

//底池
type Jackpot struct {
	m_size    int
	m_point   int
	m_maxsize int
	pot_list  PotPondList
}

func (this *Jackpot) InitJackpot(size int) {
	this.m_size = size
	this.pot_list = make(PotPondList, size)
	for i := 0; i < size; i++ {
		this.pot_list[i] = new_only_pot(i)
	}
}

func (this *Jackpot) SetBeginPot() {
	this.m_point = 0
	this.m_maxsize = 0
}

func (this *Jackpot) ResetPot() {
	this.m_point = 0
	this.m_maxsize = 0
	for i := 0; i < this.m_size; i++ {
		this.pot_list[i].reset()
	}
}

//上限池
func (this *Jackpot) SetMaxSize(size int) {
	this.m_maxsize = size
}

//弃牌的进入当前池(并非一定主池)
func (this *Jackpot) PushFoldMoney(val int32) {
	pot := this.pot_list[this.m_point]
	pot.push_fold(int64(val))
}

func (this *Jackpot) PotSize() int {
	return this.m_point + 1
}

func (this *Jackpot) TurnPot(list SeatPlayerList) {
	current_seats := make(SeatPlayerList, 0)
	var current_size int32 = 0 //当前玩家人数
	for i := 0; i < this.m_size; i++ {
		if list[i].IsPlayer() {
			current_size++
			current_seats = append(current_seats, list[i])
		}
	}
	if current_size == 0 {
		return
	}
	sort.Sort(current_seats)
	//设置每一个池上限
	var point int = 0
	top_seat := current_seats[0]
	var begin int32 = top_seat.m_init_chip
	top_seat.SetPot(point)
	var totals_money int64 = int64(top_seat.m_round_chip)
	this.set_pot_limit(point, int64(begin*current_size))
	for i := 1; i < int(current_size); i++ {
		seat := current_seats[i]
		totals_money += int64(seat.m_round_chip)
		val := int64((seat.m_init_chip - begin) * (current_size - int32(i)))
		if val > 0 {
			begin = seat.m_init_chip
			//下一个池上限
			point++
			this.set_pot_limit(point, val)
		}
		seat.SetPot(point)
	}
	this.SetBeginPot()
	this.SetMaxSize(point)
	this.set_current_jackpot(totals_money)
	//
	this.TraceJackPot()
}

//参加人数，所有人的初始筹码，当前位置，剩余筹码
func (this *Jackpot) set_current_jackpot(totals_money int64) {
	//如果进入了最后一个池，那么就写入本池，假设多筹码玩家多下了，那么应该要返回给他
	if this.m_point >= this.m_maxsize {
		this.pot_list[this.m_point].push_direct(totals_money)
	} else {
		sub_money := this.pot_list[this.m_point].push_chip(totals_money)
		//进入下一个池(sub_money==0也不会进入下一个池)
		if sub_money > 0 {
			this.m_point++
			this.set_current_jackpot(sub_money)
		}
	}
}

//获取一份
func (this *Jackpot) GetPotMoney(idx int, size int) int64 {
	return int64(math.Ceil(float64(this.pot_list[idx].totals_money / int64(size))))
}

func (this *Jackpot) set_pot_limit(p int, val int64) {
	this.pot_list[p].set_upper_limit(val)
}

func (this *Jackpot) TraceJackPot() {
	for i := 0; i < this.PotSize(); i++ {
		this.pot_list[i].trace()
	}
}

//每一个池的数据
type PotPond struct {
	pot_index     int
	current_money int64
	fold_money    int64
	upper_limit   int64
	totals_money  int64
}

func new_only_pot(idx int) *PotPond {
	this := new(PotPond)
	this.pot_index = idx
	return this
}

func (this *PotPond) set_upper_limit(val int64) {
	this.upper_limit = val
}

func (this *PotPond) push_direct(val int64) {
	this.current_money += val
	this.totals_money += val
}

func (this *PotPond) push_chip(val int64) int64 {
	if this.Overfull() {
		return val
	}
	sub_money := this.upper_limit - this.current_money
	var chip int64 = 0
	if val > sub_money {
		chip = val - sub_money
	} else {
		sub_money = val
	}
	this.current_money += int64(sub_money)
	this.totals_money += int64(sub_money)
	return chip
}

func (this *PotPond) push_fold(val int64) {
	this.fold_money += val
	this.totals_money += val
}

func (this *PotPond) Overfull() bool {
	return this.current_money >= this.upper_limit
}

func (this *PotPond) reset() {
	this.current_money = 0
	this.fold_money = 0
	this.upper_limit = 0
	this.totals_money = 0
}

func (this *PotPond) trace() {
	if this.pot_index == 0 {
		println("主池", this.pot_index, ":", this.totals_money, ",上限=", this.upper_limit)
	} else {
		println("边池", this.pot_index, ":", this.totals_money, ",上限=", this.upper_limit)
	}
}
