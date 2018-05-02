package texas

import (
	"ants/base"
	"ants/core"
	"ants/glog"
)

//--加注差额应当至少等同于当前下注轮中之前最大的下注或加注差额。
type GameLogic struct {
	//游戏固定数据
	BaseGame
	//游戏数据
	dealer_idx     int
	smallblind_idx int
	bigblind_idx   int
	need_call      int
	need_chip      int
	//座位
	cards *RobotCard
	pots  *ChipPool
	//
	clock core.ITimer
}

func NewGame(tid int, data interface{}) *GameLogic {
	this := new(GameLogic)
	this.init(tid, data)
	return this
}

func (this *GameLogic) init(tid int, data interface{}) {
	//数据
	this.room_id = tid
	this.seat_count = 9
	this.min_buyin = 100
	this.max_buyin = 10000000
	this.bet_time = 1 * 100
	this.small_blind = 50
	this.big_blind = 100
	this.table_free = 50
	this.max_look = 200
	this.min_player_num = 2
	//cards
	this.cards = newCards()
	//jackpot
	this.pots = newPots(this.seat_count)
	//center cards
	this.BaseGame.init()
}

func (this *GameLogic) OnTimeOutHandle(data interface{}) {
	switch data {
	case TIME_START:
		this.SetState(GAME_STATE_WAIT)
		this.start()
	case TIME_ACTION:
		//超时默认操作:看牌
		this.UserAction(SEAT_CHIP_CHECK, 0)
	}
}

func (this *GameLogic) OnReady() {
	//	this.clock = this.Worker().NewClock()
	//	this.clock.SetHandle(this)
	//	this.SetMessage(this)
}

func (this *GameLogic) OnDie() {
	this.clock.Stop()
}

func (this *GameLogic) OnEvent(args ...interface{}) {

}

func (this *GameLogic) check_start() int {
	num := 0
	this.each_sits(func(seat *Seat) {
		if seat.stand_check(this.table_free) {
			seat.stand()
		} else {
			num++
		}
	})
	return num
}

//流程
func (this *GameLogic) start() {
	if !this.check_state(GAME_STATE_WAIT) {
		return
	}
	//检测人数
	sit_num := this.check_start()
	//人数不够进入准备阶段
	if sit_num < this.min_player_num {
		glog.Debug("人数不够开始比赛 %d", sit_num)
		return
	}
	this.SetState(GAME_STATE_START)
	this.table_times++
	this.pots.reset()
	//设置基本参数
	this.dealer_idx = this.get_next_siter(this.dealer_idx)
	this.smallblind_idx = this.get_next_siter(this.dealer_idx)
	this.bigblind_idx = this.get_next_siter(this.smallblind_idx)
	glog.Debug("［游戏开始］Table %d Times %d", this.room_id, this.table_times)
	glog.Debug("庄家 %d 小盲 %d 大盲 %d 参与人数 %d", this.dealer_idx, this.smallblind_idx, this.bigblind_idx, sit_num)
	//设置状态
	this.each_seats(func(seat *Seat) {
		seat.beginGame()
		if act, ok := seat.checkAction(); ok {
			//扣除服务费
			act.SubFree(this.table_free)
			//glog.Debug("玩家 %d 扣除服务费 %d 剩余金额 %d", act.seat_id, this.table_free, act.seat_money)
			//扣除大小盲
			if seat.checkNewsit() {
				act.BetChip(this.big_blind)
				glog.Debug("玩家 %d 扣除坐下 %d 剩余金额 %d", seat.seat_id, this.big_blind, act.seat_money)
			} else {
				if act.seat_id == this.smallblind_idx {
					act.BetChip(this.small_blind)
					glog.Debug("玩家 %d 扣除小盲 %d 剩余金额 %d", seat.seat_id, this.small_blind, act.seat_money)
				} else if act.seat_id == this.bigblind_idx {
					act.BetChip(this.big_blind)
					glog.Debug("玩家 %d 扣除大盲 %d 剩余金额 %d", seat.seat_id, this.big_blind, act.seat_money)
				} else {
					glog.Debug("玩家 %d 无下注 剩余金额 %d", seat.seat_id, act.seat_money)
				}
			}
		}
	})
	//洗pai
	this.cards.Random()
	//发手牌
	for i := 0; i < 2; i++ {
		this.each_pos(this.dealer_idx+1, func(seat *Seat) {
			if act, ok := seat.checkAction(); ok {
				act.pushCard(this.cards.Pop())
			}
		})
	}
	//设置开始位置
	this.reset_round(this.bigblind_idx)
}

//次回合小盲开始下注
func (this *GameLogic) reset_round(ridx int) {
	this.setNoRound()
	//无攻击位置
	this.setNoAttack()
	//可以看牌
	this.need_call, _ = this.get_max_round(this.dealer_idx)
	//最小加注金额(need_call+need_chip)
	this.need_chip = this.small_blind
	//获取小盲开始找下一位活着的用户
	this.check_next(ridx)
}

func (this *GameLogic) get_max_round(pos int) (int, int) {
	att_idx := -1
	money := 0
	for i := 0; i < this.seat_count; i++ {
		idx := (pos + i) % this.seat_count
		if act, ok := this.seats[idx].checkAction(); ok {
			if act.isNoFold() {
				if act.round_chip > money {
					money = act.round_chip
					att_idx = act.seat_id
				}
			}
		}
	}
	return money, att_idx
}

//获取的下一位置
func (this *GameLogic) check_next(pos int) {
	act, ok := this.find_player(pos)
	if !ok {
		this.turn_round()
		return
	}
	glog.Debug("act=%d prev=%d attact=%d bet=%d", act.seat_id, pos, this.attack_idx, this.chip_idx)
	if this.hasAttack() {
		if this.checkAttack(act.seat_id) {
			glog.Debug("回到攻击位置:%d", act.seat_id)
			this.turn_round()
			return
		}
	} else {
		if this.checkRoundset(pos) {
			glog.Debug("回到初始位置:%d", pos)
			this.turn_round()
			return
		}
	}
	//提前结束
	if this.check_over() {
		return
	}
	//下一个玩家没有放弃
	if act.isNoFold() {
		//否则自己可以行动的话
		if act.isAction() {
			this.remind_user_action(act)
			return
		} else {
			glog.Debug("玩家已经allin:%d", act.seat_id)
		}
	} else {
		glog.Debug("玩家已经弃牌:%d", act.seat_id)
	}
	this.check_next(act.seat_id)
}

func (this *GameLogic) check_over() bool {
	//除了他其他人都放弃了
	last, ok_fold, ok_allin := this.get_only_player()
	if ok_fold {
		glog.Debug("其他人都弃,提前结束:%d", last.seat_id)
		this.over_with_action(last)
	} else if ok_allin {
		//除了他其他人都allin
		glog.Debug("多人ALLIN 亮牌结束")
		this.turn_round()
	} else {
		return false
	}
	return true
}

//提醒用户操作
func (this *GameLogic) remind_user_action(act *GameAction) {
	this.setCurrent(act.seat_id)
	max_chip := this.get_max_chip(act.seat_id, act.Totals())
	this.clock.Onec(this.bet_time, TIME_ACTION)
	glog.Debug("提醒: 玩家 %d 操作 [call=%d chip=%d max=%d money=%d]",
		act.seat_id, this.need_call, this.need_call+this.need_chip, max_chip, act.seat_money)
}

func (this *GameLogic) UserAction(action int, chip int) {
	if !this.isplaying() {
		return
	}
	seat, ok := this.current_action()
	//glog.Debug("玩家操作 %d 剩余 %d", seat.seat_id, seat.round_chip)
	round_chip := seat.round_chip
	seat_money := seat.seat_money
	if !ok {
		panic("无人操作，系统错误！")
	}
	if action == SEAT_CHIP_FOLD {
		seat.setFold()
		glog.Debug("[弃牌] 玩家 %d 剩余 %d", seat.seat_id, seat_money)
	} else if action == SEAT_CHIP_CHECK {
		if this.need_call > round_chip {
			seat.setFold()
			glog.Debug("[弃牌:必须下注] 玩家 %d 剩余 %d", seat.seat_id, seat_money)
		} else {
			glog.Debug("[看牌] 玩家 %d 剩余 %d", seat.seat_id, seat_money)
		}
	} else if action == SEAT_CHIP_CALL {
		if this.need_call == 0 {
			//不能跟注，看牌不管
		} else {
			seat.BetChip(this.need_call)
			glog.Debug("[跟注 %d] 玩家 %d 剩余 %d", seat.seat_id, round_chip, seat_money)
		}
	} else if action == SEAT_CHIP_RAISE {
		this.raise_chip(chip, seat)
		glog.Debug("[加注 %d] 玩家 %d 剩余 %d", seat.seat_id, round_chip, seat_money)
	} else {
		glog.Debug("无法识别的操作:", seat.seat_id, action)
		return
	}
	this.check_next(seat.seat_id)
}

//加注
func (this *GameLogic) raise_chip(chip int, action *GameAction) {
	//min_fill := this.need_call + this.need_chip
	round_chip := action.round_chip
	total_chip := action.Totals()
	if chip > total_chip {
		chip = total_chip
	}
	action.BetChip(chip)
	if chip > this.need_call {
		//if chip < min_fill {加注全下}
		this.set_attack_seat(action.seat_id, round_chip)
	} else {
		//跟注全下
	}
}

//如果是加注那么重新设置攻击位置
func (this *GameLogic) set_attack_seat(seat_id int, chip int) {
	need_chip := chip - this.need_call
	if need_chip > this.need_chip {
		this.need_chip = need_chip
	}
	this.need_call = chip
	this.setAttack(seat_id)
}

//gets
func (this *GameLogic) get_max_chip(seat_id int, round int) int {
	chip := 0
	this.each_actions(func(action *GameAction) {
		if action.seat_id != seat_id && action.isNoFold() {
			money := action.Totals()
			if money > chip {
				chip = money
			}
		}
	})
	if chip > round {
		return round
	}
	return chip
}

//3个阶段
func (this *GameLogic) turn_round() {
	this.SetState(this.game_state + 1)
	//分池
	this.pots.TurnPot(this.seats)
	//置空下注
	this.each_actions(func(action *GameAction) {
		action.turnOver()
	})
	switch this.game_state {
	case GAME_START_FLOP:
		this.pushCard(this.cards.Pop())
		this.pushCard(this.cards.Pop())
		this.pushCard(this.cards.Pop())
		glog.Debug("［翻牌］" + ToStr(this.public_cards))
	case GAME_START_TURN:
		this.pushCard(this.cards.Pop())
		glog.Debug("［转牌］" + CardStr(this.public_cards[3]))
	case GAME_START_RIVER:
		this.pushCard(this.cards.Pop())
		glog.Debug("［河牌］" + CardStr(this.public_cards[4]))
	case GAME_START_STOP:
		glog.Debug("［公共牌］" + ToStr(this.public_cards))
		this.over()
		return
	default:
	}
	this.reset_round(this.dealer_idx)
}

func (this *GameLogic) over() {
	glog.Debug(">>开始结算")
	this.over_result()
	this.over_clear(this.pots.PotSize())
	glog.Debug("［游戏结束］")
}

//提前结束(也就是只存在一个人)
func (this *GameLogic) over_with_action(action *GameAction) {
	if !this.check_state_set(GAME_START_STOP) {
		return
	}
	glog.Debug("［游戏提前结束］pot=%d", this.pots.PotSize())
	this.pots.TurnPot(this.seats)
	for p := 0; p < this.pots.PotSize(); p++ {
		money := this.pots.potMoney(p, 1)
		action.Result(int(money))
		glog.Debug("直接结算: 玩家 %d 赢[%d:%d] 金额 %d", action.seat_id, p, money, action.seat_money)
	}
	this.over_clear(1)
}

func (this *GameLogic) over_clear(size int) {
	this.each_seats(func(seat *Seat) {
		seat.endGame()
	})
	this.clearCards()
	//等待开始 每个池2秒 外加1秒
	this.clock.Onec(2000*size+1000, TIME_START)
}

func (this *GameLogic) over_result() {
	var result_list []*GameAction
	//未弃牌的进入比牌
	this.each_actions(func(action *GameAction) {
		if action.isNoFold() {
			action.cardFlush(this.public_cards)
			result_list = append(result_list, action)
		}
	})
	//绝对牌型排序
	base.Sort(result_list, func(i, j int) bool {
		return result_list[i].compare(result_list[j]) == WIN
	})
	//test打印
	for i := range result_list {
		println(ToStr(result_list[i].big_cards), TypeStr(result_list[i].card_type))
	}
	//分池
	for p := 0; p < this.pots.PotSize(); p++ {
		list := this.get_pot_seats(result_list, p)
		money := this.pots.potMoney(p, len(list))
		for i := range list {
			act := list[i]
			act.Result(int(money))
			glog.Debug("结算> 玩家:%d 参与奖池:%d 赢得:[%d:%d] 金额:%d", act.uid, act.pot_point, p, money, act.seat_money)
		}
	}
}

func (this *GameLogic) get_pot_seats(list []*GameAction, pot int) []*GameAction {
	var big_list []*GameAction = nil
	var big_action *GameAction = nil
	for i := range list {
		act := list[i]
		if act.pot_point >= pot {
			if big_action == nil {
				big_action = act
				big_list = append([]*GameAction{}, act)
			} else {
				if big_action.card_type > act.card_type {
					break
				}
				ret := big_action.compare(act)
				if ret == DRAW {
					big_list = append(big_list, act)
				} else {
					break
				}
			}
		}
	}
	return big_list
}

//message handle

func (this *GameLogic) enter(uid int, name string, gift int) {
	if player, ok := this.get_player(uid); ok {
		player.update(name, gift)
	} else {
		if this.length() >= this.max_look {
			glog.Debug("over full table max = %d", this.max_look)
		} else {
			this.setPlayer(uid, newPlayer(uid, name, gift))
		}
	}
}

func (this *GameLogic) leave(uid int) {
	if player, ok := this.delUser(uid); ok {
		this.stand(uid)
		glog.Debug("玩家 uid=%d 离开游戏", player.uid)
	}
}

//不分座位id
func (this *GameLogic) sit(uid int, money int, seat_id int, autobuy bool) {
	if !this.check_buyin(money) {
		//带入钱出现问题
		return
	}
	if this.has(uid) {
		if seat, ok := this.get_seat(seat_id); ok {
			if seat.issit() {
				if seat.uid == uid {
					seat.update(money, autobuy)
				} else {
					//其他人占用了
				}
			} else {
				seat.sit(uid, money, autobuy)
			}
		} else {
			this.auto_sit(uid, money, autobuy)
		}
		//坐满即开
		this.start()
	}
}

func (this *GameLogic) auto_sit(uid int, money int, autobuy bool) bool {
	seat, ok := this.find_with(func(val *Seat) bool {
		return !val.issit()
	})
	//找到了无人坐下的
	if ok {
		seat.sit(uid, money, autobuy)
	}
	return ok
}

func (this *GameLogic) stand(uid int) {
	if seat, ok := this.getSeatByUid(uid); ok {
		if seat.issit() {
			seat.stand()
			this.leave_action(seat)
		}
	}
}

//玩家操作
func (this *GameLogic) chip_action(uid int, action int, chip int) {
	if seat, ok := this.getSeatByUid(uid); ok {
		if this.isCurrent(seat.seat_id) {
			this.UserAction(action, chip)
		}
	}
}

//判断离开游戏
func (this *GameLogic) leave_action(seat *Seat) {
	if act, ok := seat.checkAction(); ok {
		if this.isCurrent(act.seat_id) {
			//正在操作，那么就弃牌
			this.UserAction(SEAT_CHIP_FOLD, 0)
		} else {
			act.setFold()
		}
	}
}

func _init() {
	base.Sleep(10)
	m := NewGame(1, nil)

	for i := 1; i < 6; i++ {
		m.enter(i, "test1", i)
		m.sit(i, i*100, i, false)
	}
}
