package texas

import "ants/base"
import "math"

//底池
type ChipPool struct {
	m_size    int
	m_point   int
	m_maxsize int
	pot_list  []*PotPond
}

func newPots(sz int) *ChipPool {
	this := new(ChipPool)
	this.init(sz)
	return this
}

func (this *ChipPool) init(size int) {
	this.m_size = size
	this.pot_list = make([]*PotPond, size)
	for i := 0; i < size; i++ {
		this.pot_list[i] = new_only_pot(i)
	}
}

//func
func (this *ChipPool) begin() {
	this.m_point = 0
	this.m_maxsize = 0
}

func (this *ChipPool) reset() {
	this.m_point = 0
	this.m_maxsize = 0
	for i := 0; i < this.m_size; i++ {
		this.pot_list[i].reset()
	}
}

func (this *ChipPool) maxSize(size int) {
	this.m_maxsize = size
}

func (this *ChipPool) PotSize() int {
	return this.m_point + 1
}

func (this *ChipPool) TurnPot(seats []*Seat) {
	var actions []*GameAction
	fold_money := 0
	for i := range seats {
		act, ok := seats[i].checkAction()
		if ok {
			if act.isFold() {
				fold_money += act.round_chip
			} else {
				actions = append(actions, act)
			}
		}
	}
	current_size := len(actions)
	//最小的在前面
	base.Sort(actions, func(i, j int) bool {
		return actions[i].init_chip < actions[j].init_chip
	})
	//设置每一个池上限
	point := 0
	top_seat := actions[0]
	init_chip := top_seat.init_chip
	top_seat.pot_point = point
	totals_money := int64(top_seat.round_chip)
	this.setPotLimit(point, int64(init_chip*current_size))
	for i := 1; i < current_size; i++ {
		act := actions[i]
		totals_money += int64(act.round_chip)
		val := int64((act.init_chip - init_chip) * (current_size - i))
		if val > 0 {
			init_chip = act.init_chip
			point++
			this.setPotLimit(point, val)
		}
		act.pot_point = point
	}
	//弃牌的放在之前的池
	this.addFoldMoney(fold_money)
	//是否重新计算
	//this.begin()
	this.maxSize(point)
	//添加筹码
	this.pushMoney(totals_money)
	//打印
	this.Trace()
}

//入池
func (this *ChipPool) pushMoney(totals_money int64) {
	if this.m_point >= this.m_maxsize {
		this.pot_list[this.m_point].push_direct(totals_money)
	} else {
		sub_money := this.pot_list[this.m_point].push_chip(totals_money)
		if sub_money > 0 {
			this.m_point++
			this.pushMoney(sub_money)
		}
	}
}

//弃牌的进入当前池(并非一定主池)
func (this *ChipPool) addFoldMoney(val int) {
	this.pot_list[this.m_point].push_fold(int64(val))
}

//获取池金币
func (this *ChipPool) potMoney(idx int, size int) int64 {
	return int64(math.Ceil(float64(this.pot_list[idx].totals_money / int64(size))))
}

//设置池上限
func (this *ChipPool) setPotLimit(p int, val int64) {
	this.pot_list[p].set_upper_limit(val)
}

func (this *ChipPool) Trace() {
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
		println("主池", this.pot_index, ",数额=", this.totals_money, ",上限=", this.upper_limit)
	} else {
		println("边池", this.pot_index, ",数额=", this.totals_money, ",上限=", this.upper_limit)
	}
}
