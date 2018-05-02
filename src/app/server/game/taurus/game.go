package taurus

import (
	"ants/base"
	"ants/core"
	"ants/gcode"
	"ants/glog"
	"app/command"
)

//牛牛逻辑
type GameLogic struct {
	GameSend
	cards *CardDealer
	//管理
	clock core.ITimer
}

func NewGame(tid int, data interface{}) *GameLogic {
	this := new(GameLogic)
	this.init(tid, data)
	return this
}

func (this *GameLogic) init(tid int, data interface{}) {
	//init
	this.game_state = GAME_WAIT
	this.table_id = tid
	this.seat_count = 4
	this.banker_time = 1 * 1000
	this.chip_time = 1 * 1000
	this.commit_time = 1 * 1000
	this.over_time = 1 * 1000
	this.min_player = 2
	this.base_chip = 25
	this.min_chip = 5 //最小5倍
	this.max_chip = 25
	this.max_look = 200
	//
	this.cards = NewDealer()
	this.BaseGame.init()
}

func (this *GameLogic) OnReady() {
	//	this.clock = this.Worker().NewClock()
	//	this.SetMessage(this)
	//等待
	this.clock.SetTimeout(func(data interface{}) {
		//处理计时器
		switch data {
		case GAME_START:
			this.rob_end()
		case GAME_CHIP:
			this.chip_over()
		case GAME_COMMIT:
			this.commit_over()
		case GAME_STOP:
			this.game_state = GAME_WAIT
			this.start()
		}
	})
}

func (this *GameLogic) OnDie() {
	this.clock.Stop()
}

func (this *GameLogic) OnEvent(pack gcode.ISocketPacket) {
	switch pack.Cmd() {
	case command.CLIENT_GAME_DROPS:
		this.leave(pack.ReadInt()) //掉线提醒离开
	case command.CLIENT_GAME_RECONNENT:
		this.reconnect(pack.ReadInt(), pack.ReadInt(), pack.ReadUInt64())
	case command.CLIENT_GAME_ENTER:
		this.enter(pack.ReadInt(), pack.ReadInt(), pack.ReadUInt64())
	case command.CLIENT_GAME_LEAVE:
		this.leave(pack.ReadInt())
	case command.CLIENT_GAME_SIT:
		this.sit(pack.ReadInt(), pack.ReadByte(), pack.ReadInt(), pack.ReadBool())
	case command.CLIENT_GAME_STAND:
		this.stand(pack.ReadInt())
	case command.CLIENT_NIUNIU_BANKER:
		this.rob_banker(pack.ReadInt(), pack.ReadByte())
	case command.CLIENT_NIUNIU_BET:
		this.chip_action(pack.ReadInt(), pack.ReadByte())
	case command.CLIENT_NIUNIU_COMMIT:
		this.commit(pack.ReadInt())
	default:

	}
}

func (this *GameLogic) check_start() int {
	num := 0
	this.each_sits(func(seat *Seat) {
		if seat.stand_check(this.base_chip) {
			this.check_stand_over(seat.uid, seat)
		} else {
			num++
		}
	})
	return num
}

//1,开始游戏(至少2人才能开始比赛)
func (this *GameLogic) start() {
	if !this.check_state(GAME_WAIT) {
		return
	}
	player_num := this.check_start()
	//开始抢庄
	if player_num < this.min_player {
		glog.Debug("无法开始游戏 人数太少 %d", player_num)
		return
	}
	glog.Debug("[游戏开始] 房间号 %d 场次 %d 参加人数 %d", this.table_id, this.table_num, player_num)
	this.SetState(GAME_START)
	this.table_num++
	this.banker_multiple = 0
	this.banker_idx = 0
	this.each_sits(func(seat *Seat) {
		seat.begin()
	})
	//洗牌
	this.cards.Random()
	//每人发4张牌
	this.dealer_cards(0, 4)
	//计时器
	this.clock.Onec(this.banker_time, GAME_START)
	//广播开始
	this.broadcast_start()
	//广播自己的牌
	this.each_actions(func(seat *Seat) {
		this.send_self_cards(seat)
	})
}

//on user 抢庄 倍数(1-4)倍
func (this *GameLogic) rob_banker(uid int, chip int) {
	if seat, ok := this.getSeatByUid(uid); ok && seat.isplayer() {
		seat.rob_banker(chip)
		this.broadcast_rob_banker(seat.seat_id, chip)
	}
}

//抢庄结束，确定庄
func (this *GameLogic) rob_end() {
	this.SetState(GAME_CHIP)
	//确定庄
	var bankers []*Seat
	this.each_actions(func(seat *Seat) {
		if seat.bet_multiple > this.banker_multiple {
			this.banker_multiple = seat.bet_multiple
			bankers = append([]*Seat{}, seat)
		} else if seat.bet_multiple == this.banker_multiple {
			bankers = append(bankers, seat)
		}
	})
	//随机一个庄
	rand_idx := base.Random(len(bankers))
	this.banker_idx = bankers[rand_idx].seat_id
	//不能为0
	if this.banker_multiple == 0 {
		this.banker_multiple = 1
	}
	glog.Debug("[抢庄结束] 庄家 %d 倍数 %d", this.banker_idx, this.banker_multiple)
	//开始下注
	this.start_chip()
}

//2，开始下注
func (this *GameLogic) start_chip() {
	this.clock.Onec(this.chip_time, GAME_CHIP)
	glog.Debug("2, 开始下注")
	this.broadcast_bet_start()
}

//on user
func (this *GameLogic) chip_action(uid int, num int) {
	if seat, ok := this.getSeatByUid(uid); ok && seat.isplayer() {
		if seat.seat_id != this.banker_idx {
			seat.bet_chip(num)
			this.broadcast_user_bet(seat.seat_id, num)
		}
	}
}

func (this *GameLogic) chip_over() {
	this.each_actions(func(seat *Seat) {
		if seat.seat_id != this.banker_idx {
			seat.bet_chip(this.min_chip)
		}
	})
	//发牌
	this.dealer_cards(4, 5)
	//广播自己的牌5长
	this.each_actions(func(seat *Seat) {
		this.send_self_cards(seat)
	})
	//最后确认
	this.start_commit()
}

//给玩游戏的人发牌
func (this *GameLogic) dealer_cards(bpos int, epos int) {
	for i := bpos; i < epos; i++ {
		this.each_actions(func(seat *Seat) {
			seat.pushCard(this.cards.Pop())
		})
	}
	//test
	this.each_actions(func(seat *Seat) {
		glog.Debug("座位 %d 手牌 %s", seat.seat_id, ToStr(seat.cards...))
	})
}

//3,开始计时提交
func (this *GameLogic) start_commit() {
	this.clock.Onec(this.commit_time, GAME_COMMIT)
	glog.Debug("3, 开始算牌")
	this.broadcast_commit_start()
}

//on user
func (this *GameLogic) commit(uid int) {
	if seat, ok := this.getSeatByUid(uid); ok && seat.isplayer() {
		seat.commit()
		this.broadcast_user_commit(uid)
	}
	//提交结束
	if this.is_all_commit() {
		this.commit_over()
	}
}

func (this *GameLogic) is_all_commit() bool {
	ok := true
	this.each_actions(func(seat *Seat) {
		if ok && !seat.iscommit() {
			ok = false
		}
	})
	return ok
}

//4,提交后结束游戏
func (this *GameLogic) commit_over() {
	this.over()
}

func (this *GameLogic) over() {
	if !this.check_state_set(GAME_STOP) {
		return
	}
	this.over_result()
	this.over_clear()
	glog.Debug("[游戏结束]")
	this.broadcast_over()
}

//结算 >>输的先给庄，庄再最高分配
func (this *GameLogic) over_result() {
	banker, ok := this.get_seat(this.banker_idx)
	if !ok {
		panic("系统错误，无法找到庄家")
	}
	//牌型结算
	banker.card_flush()
	//一个个结算
	var players []*Seat
	this.each_actions(func(seat *Seat) {
		if seat.seat_id != this.banker_idx {
			seat.card_flush()
			players = append(players, seat)
		}
	})
	//牌型最大的排名
	base.Sort(players, func(i, j int) bool {
		return players[i].check_result(players[j]) == WIN
	})
	glog.Debug("4, 游戏结算: 庄家 %d 倍数 %d 基础分 %d 牌型 %s",
		this.banker_idx, this.banker_multiple, this.base_chip, banker.card_str())
	//输的先给庄(赢钱只能赢带入一倍的钱)
	for i := range players {
		seat := players[i]
		ret := banker.check_result(seat)
		if ret == WIN {
			money := this.base_chip * this.banker_multiple * banker.card_multiple() * seat.bet_multiple
			money = seat.sub_money(money)
			seat.result_set(LOSE)   //闲输
			banker.add_money(money) //庄赢钱
			glog.Debug("庄家 %d 庄赢 %d 闲倍数 %d 庄牌型 %s 闲牌型 %s 座位钱 %d",
				this.banker_idx, money, seat.bet_multiple, banker.card_str(), seat.card_str(), banker.seat_money)
		} else if ret == DRAW {
			seat.result_set(DRAW)
		}
	}
	//赢钱的排序
	for i := range players {
		seat := players[i]
		if seat.check_result(banker) == WIN {
			money := this.base_chip * seat.bet_multiple * this.banker_multiple * seat.card_multiple()
			money = banker.sub_money(money)
			seat.result_set(WIN)
			seat.add_money(money) //闲赢钱
			glog.Debug("闲家 %d 闲赢 %d 倍数 %d 庄牌型 %s 闲牌型 %s 座位钱 %d",
				seat.seat_id, money, seat.bet_multiple, banker.card_str(), seat.card_str(), seat.seat_money)
		}
	}
	//广播最后的结果
	this.broadcast_result(players)
}

//清理
func (this *GameLogic) over_clear() {
	this.each_actions(func(seat *Seat) {
		seat.over()
	})
	//按照人数来结束
	this.clock.Onec(this.over_time, GAME_STOP)
}

//action
func (this *GameLogic) reconnect(uid int, gate int, session uint64) {
	if player, ok := this.get_player(uid); ok {
		//这个时候返回桌子信息
		player.update(gate, session)
		this.send_reconnect(player)
		this.do_enter(player)
	}
}

func (this *GameLogic) enter(uid int, gate int, session uint64) {
	if player, ok := this.get_player(uid); ok {
		player.update(gate, session)
	} else {
		if this.length() >= this.max_look {
			glog.Debug("over full table max = %d", this.max_look)
		} else {
			player := newPlayer(uid, gate, session)
			this.setPlayer(uid, player)
			this.do_enter(player)
		}
	}
}

//一定返回
func (this *GameLogic) leave(uid int) {
	player, ok := this.get_player(uid)
	if ok {
		ok_exit := true
		if seat, ok1 := this.getSeatByUid(uid); ok1 {
			if seat.isplayer() {
				seat.over_stand()
				ok_exit = false
			} else {
				this.stand(uid)
			}
		}
		if ok_exit {
			this.do_remove_user(uid)
		} else {
			player.over_leave()
		}
	}
}

func (this *GameLogic) check_leave(uid int) {
	if player, ok := this.get_player(uid); ok {
		//掉线了就删除
		if player.isLeave() {
			this.do_remove_user(uid)
		}
	}
}

func (this *GameLogic) sit(uid int, seat_id int, money int, autobuy bool) {
	player, ok := this.get_player(uid)
	if !ok {
		return
	}
	if this.min_money > money {
		glog.Debug("sit err money is lose %d", uid)
		return
	}
	if seat, ok1 := this.getSeatByUid(uid); ok1 {
		seat.update(money, autobuy)
	} else {
		if seat, ok2 := this.auto_sit(seat_id, uid, money, autobuy); ok2 {
			seat.sit(uid, money, autobuy, player)
			this.do_sit(seat, player)
		}
	}
}

func (this *GameLogic) auto_sit(seat_id int, uid int, money int, autobuy bool) (*Seat, bool) {
	for i := 0; i < this.seat_count; i++ {
		idx := (seat_id + i) % this.seat_count
		if seat, ok := this.get_seat(idx); ok {
			if seat.issit() {
				continue
			}
			return seat, true
		}
	}
	return nil, false
}

func (this *GameLogic) stand(uid int) {
	if seat, ok := this.getSeatByUid(uid); ok {
		if seat.isplayer() {
			seat.over_stand()
		} else {
			this.do_stand(seat)
		}
	}
}

//最后通告
func (this *GameLogic) check_stand_over(uid int, seat *Seat) {
	this.do_stand(seat)
	this.check_leave(uid)
}

//最终离开
func (this *GameLogic) do_remove_user(uid int) {
	delete(this.users, uid)
	this.broadcast_leave(uid)
}

func (this *GameLogic) do_stand(seat *Seat) {
	this.broadcast_stand(seat.seat_id)
	seat.stand()
}

func (this *GameLogic) do_enter(player *Player) {
	this.send_table_info(player)
	//通知其他
	this.broadcast_enter(player.uid, player.name)
	glog.Debug("enter player uid=%d %d", player.uid, len(this.users))
}

func (this *GameLogic) do_sit(seat *Seat, player *Player) {
	this.broadcast_sit(seat.seat_id, player.uid, seat.seat_money, player.name, player.gift)
	//
	glog.Debug("sit player uid=%d seat=%d", player.uid, seat.seat_id)
	//可以开始游戏
	this.start()
}

func _init() {
	//t := NewGame(101, nil)
	//	t.enter(1, "test1", 1)
	//	t.sit(1, 1000, 1, false)

	//	t.enter(2, "test2", 2)
	//	t.sit(2, 1000, 2, false)
}
