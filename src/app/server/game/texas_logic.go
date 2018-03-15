package game

import (
	"ants/gnet"
	"ants/gsys"
	"ants/gutil"
	"app/command"
	"app/server/game/texas"
	"sort"
)

//--加注差额应当至少等同于当前下注轮中之前最大的下注或加注差额。
type TexasLogic struct {
	texas.Jackpot
	texas.RobotCard
	//房间数据
	room_id     int
	room_type   int
	seat_count  int
	cool_time   int
	small_blind int32
	big_blind   int32
	min_buyin   int32
	max_buyin   int32
	table_free  int32
	max_look    int
	//游戏数据
	m_game_state     int
	m_current_player int
	m_dealer_id      int
	m_small_seat_id  int
	m_big_seat_id    int
	m_chip_seat_id   int
	m_attack_seat_id int
	m_need_call      int32
	m_need_chip      int32
	m_public_cards   [7]int16
	//座位
	seat_list texas.SeatPlayerList
	players   map[int32]*texas.PlayerVo
	//
	m_timer gsys.ITimer
	m_chan  gsys.IAsynDispatcher
}

func NewTexasLogic() ITableLogic {
	this := new(TexasLogic)
	this.InitTexasLogic()
	return this
}

func (this *TexasLogic) InitTexasLogic() {
	this.InitRobotCard()
	this.m_chan = gsys.NewChannel()
	//同步事务
	go func() {
		this.m_chan.Loop(func(args interface{}) {
			this.OnNotice(args.([]interface{})...)
		})
	}()
	//
	this.m_timer = gsys.NewTimerWithChannel(this.m_chan)
	this.m_timer.SetHandle(func(data interface{}) {
		switch data.(int) {
		case texas.TIME_OVER_STATE:
			this.set_ready()
			this.StartGame()
		case texas.TIME_CHIP_STATE:
			this.outtime_chip()
		}
	})
}

func (this *TexasLogic) OnNotice(args ...interface{}) {
	pack := args[0].(gnet.ISocketPacket)
	header := args[1].(*GameHeader)
	switch pack.Cmd() {
	case command.CLIENT_ENTER_TEXAS_ROOM:
		this.OnEnterTable(texas.NewPlayerVo(header.UserID, header.GateID, header.SessionID, pack))
	case command.CLIENT_LEAVE_TEXAS_ROOM:
		this.OnLeaveTable(header.UserID)
	case command.CLIENT_TEXAS_SITDOWN:
		seat_id := args[2].(int8)
		seat_money := args[3].(int32)
		//auto_buy := args[4].(int8)
		this.OnSitdown(header.UserID, seat_id, seat_money)
	case command.CLIENT_TEXAS_STAND:
		this.OnStand(header.UserID)
	case command.CLIENT_TEXAS_CHIP:
		atype := args[2].(int8)
		chip := args[3].(int32)
		this.OnChipAction(header.UserID, atype, chip)
	}
}

func (this *TexasLogic) PushNotice(args ...interface{}) {
	this.m_chan.Push(args)
}

//通知所有
func (this *TexasLogic) notice_players(cmd int, data interface{}) {
	var list []*texas.PlayerVo
	for _, v := range this.players {
		list = append(list, v)
	}
	//gnet.NewPacketWithTopic(cmd, config.TOPIC_CLIENT, data)
	//	for _, v := range list {

	//	}
}

func (this *TexasLogic) get_seat(idx int) *texas.SeatPlayer {
	if idx >= this.seat_count || idx < 0 {
		return nil
	}
	return this.seat_list[idx]
}

func (this *TexasLogic) get_chip_seat() *texas.SeatPlayer {
	return this.get_seat(this.m_chip_seat_id)
}

func (this *TexasLogic) Type() int {
	return 1
}

func (this *TexasLogic) TableID() int {
	return this.room_id
}

func (this *TexasLogic) OnLaunch(tid int, data interface{}) {
	this.room_id = tid
	this.seat_count = 9
	this.min_buyin = 100
	this.max_buyin = 10000000
	this.cool_time = 3000
	this.small_blind = 50
	this.big_blind = 100
	this.table_free = 50
	this.m_dealer_id = 0
	this.max_look = 200
	//jackpot
	this.InitJackpot(this.seat_count)
	//seats
	this.seat_list = make(texas.SeatPlayerList, this.seat_count)
	for i := 0; i < this.seat_count; i++ {
		this.seat_list[i] = texas.NewTexasSeatPlayer(i)
	}
	//println("总人数啊:", len(this.seat_list))
	//user
	this.players = make(map[int32]*texas.PlayerVo)
	//设置为开始
	this.set_ready()
}

func (this *TexasLogic) OnFree() {
	unRegTable(this)
	this.m_timer.Stop()
	this.m_chan.Close()
}

//流程
func (this *TexasLogic) StartGame() {
	if this.m_game_state != texas.GAME_STATE_READY {
		println("[游戏尚未结束]")
		return
	}
	this.ResetPot()
	sit_num := this.check_table()
	this.m_timer.Stop()
	//人数不够进入准备阶段
	if sit_num < 2 {
		println("#Ready For Num:", sit_num)
		this.m_game_state = texas.GAME_STATE_READY
		return
	}
	this.m_game_state = texas.GAME_STATE_START
	println("#Game start player num = ", sit_num)
	//设置基本参数
	this.m_dealer_id = this.get_next_player(this.m_dealer_id)
	this.m_small_seat_id = this.get_next_player(this.m_dealer_id)
	this.m_big_seat_id = this.get_next_player(this.m_small_seat_id)
	//如果出现这个，那么系统有问题
	if this.m_dealer_id == -1 || this.m_small_seat_id == -1 || this.m_big_seat_id == -1 {
		panic("系统存在问题")
	}
	println("庄家=", this.m_dealer_id, " 小盲=", this.m_small_seat_id, " 大盲=", this.m_big_seat_id)
	//设置状态
	for _, seat := range this.seat_list {
		if seat.IsSit() {
			seat.Begin()
			//扣除服务费
			seat.SubTableFree(this.table_free)
			println(seat.SeatID(), "玩家 扣除服务费=", this.table_free, " 座位钱=", seat.SeatMoney())
			//大小盲
			if seat.CheckSeatID(this.m_small_seat_id) {
				seat.SubChip(this.small_blind)
				println("玩家小盲注:", seat.SeatID(), this.small_blind)
			} else if seat.CheckSeatID(this.m_big_seat_id) {
				seat.SubChip(this.big_blind)
				println("玩家大盲注:", seat.SeatID(), this.big_blind)
			}
		}
	}
	//手牌
	this.dealer_hand_cards()
	//下注参数
	this.reset_round(this.big_blind, this.m_big_seat_id)
	//开始下注
	this.ChipNextStart()
}

//是否可以开始
func (this *TexasLogic) check_table() int {
	//牌太少要洗牌
	if this.CardAvailable() < this.seat_count*2+5 {
		this.RandomCards()
	}
	//剔除筹码不足的玩家
	var sit_num int = 0
	for _, seat := range this.seat_list {
		if seat.IsStand() {
			continue
		}
		//费用不足(服务费＋小盲)
		if this.table_free+this.small_blind > seat.SeatMoney() {
			this.StandSeat(seat)
		} else {
			sit_num++
		}
	}
	return sit_num
}

//发2张牌(小盲开始发起)
func (this *TexasLogic) dealer_hand_cards() {
	for m := 0; m < 2; m++ {
		var index int = this.m_dealer_id
		for i := 0; i < this.seat_count; i++ {
			index++
			if index >= this.seat_count {
				index = index - this.seat_count
			}
			seat := this.seat_list[index]
			if seat.IsPlayer() {
				seat.SetCard(m, this.NextCard())
			}
		}
	}
}

//获取下一个坐下的玩家
func (this *TexasLogic) get_next_player(val int) int {
	pos := val + 1
	if pos >= this.seat_count {
		pos = 0
	}
	for i := 0; i < this.seat_count; i++ {
		idx := (pos + i) % this.seat_count
		if this.seat_list[idx].IsSit() {
			return idx
		}
	}
	return -1
}

//次回合小盲开始下注
func (this *TexasLogic) reset_round(call int32, begin_id int) {
	this.m_chip_seat_id = begin_id
	this.m_attack_seat_id = begin_id
	this.m_need_call = call
	//初始化的最小加注金额
	this.m_need_chip = this.small_blind
}

//3个回合
func (this *TexasLogic) turn_round() {
	this.m_game_state++
	println("#Turn Round:", this.m_game_state, "========================")
	this.reset_round(0, this.m_dealer_id)
	//push
	this.TurnPot(this.seat_list)
	//clean
	for _, seat := range this.seat_list {
		seat.Turn()
	}
	//state
	if this.m_game_state == texas.GAME_START_FLOP {
		for i := 0; i < 3; i++ {
			this.m_public_cards[i] = this.NextCard()
		}
		list := this.m_public_cards
		str := texas.PokerString(list[0]) + "," + texas.PokerString(list[1]) + "," + texas.PokerString(list[2])
		println("翻牌[", str, "]")
	} else if this.m_game_state == texas.GAME_START_TURN {
		this.m_public_cards[3] = this.NextCard()
		println("转牌[", texas.PokerString(this.m_public_cards[3]), "]")
	} else if this.m_game_state == texas.GAME_START_RIVER {
		this.m_public_cards[4] = this.NextCard()
		println("河牌[", texas.PokerString(this.m_public_cards[4]), "]")
	} else if this.m_game_state == texas.GAME_START_STOP {

	}
	//游戏结束
	if this.m_game_state == texas.GAME_START_STOP {
		this.over()
	} else {
		//开始下注
		this.ChipNextStart()
	}
}

//提前结束
func (this *TexasLogic) over_premature() {
	println("#Premature Over")
	this.TurnPot(this.seat_list)
	this.m_game_state = texas.GAME_START_STOP
	//分池
	for i := 0; i < this.PotSize(); i++ {
		money := this.GetPotMoney(i, 1)
		for _, seat := range this.seat_list {
			if seat.IsPlayer() {
				seat.Result(money)
				println("玩家:", seat.SeatID(), " 奖池ID:", i, " 金额:", money, " 总金额:", seat.SeatMoney())
			}
		}
	}
	//设置结束
	this.over_set(0)
}

func (this *TexasLogic) over() {
	println("#Game Over")
	this.m_game_state = texas.GAME_START_STOP
	this.result_handler()
	this.over_set(this.PotSize())
}

func (this *TexasLogic) over_set(size int) {
	for _, seat := range this.seat_list {
		seat.Over()
	}
	//等待开始 每个池2秒 外加1秒
	this.m_timer.Onec(2000*size+1000, texas.TIME_OVER_STATE)
}

func (this *TexasLogic) result_handler() {
	var result_list texas.SeatResultList
	for _, seat := range this.seat_list {
		if seat.IsPlayer() {
			this.m_public_cards[5] = seat.GetCard(0)
			this.m_public_cards[6] = seat.GetCard(1)
			//println("在线:", seat.SeatID(), PokerStrings(this.m_public_cards[0:7]))
			texas.CardTypeOfTexas(this.m_public_cards[0:7])
			result := texas.NewSeatResult(seat, this.m_public_cards[0:7])
			result_list = append(result_list, result)

		}
	}
	sort.Sort(result_list)
	for _, r := range result_list {
		r.Trace()
	}
	//分池
	for i := 0; i < this.PotSize(); i++ {
		list := this.get_max_result(result_list, i)
		money := this.GetPotMoney(i, len(list))
		for _, result := range list {
			result.Seat.Result(money)
			println("玩家:", result.Seat.SeatID(), " 参与奖池=", result.PotIdx, " 奖池ID:",
				i, " 金额:", money, " 总金额:", result.Seat.SeatMoney())
		}
	}
}

func (this *TexasLogic) get_max_result(list texas.SeatResultList, pot int) texas.SeatResultList {
	big_list := make(texas.SeatResultList, list.Len())
	var big_result *texas.SeatResult = nil
	var point int = 0
	//最大类型开始轮
	for _, result := range list {
		if result.PotIdx >= pot {
			if big_result == nil {
				big_result = result
				big_list[point] = result
				point++
			} else {
				if big_result.CardType > result.CardType {
					break
				}
				ret := texas.CompareCards(big_result.CardList, result.CardList)
				if ret == -1 {
					big_result = result
					point = 0
					big_list[point] = result
					point++
				} else if ret == 1 {
					//不大于
				} else {
					big_list[point] = result
					point++
				}
			}
		}
	}
	return big_list[0:point]
}

func (this *TexasLogic) set_ready() {
	this.m_game_state = texas.GAME_STATE_READY
}

func (this *TexasLogic) ChipNextStart() {
	if !this.is_playing() {
		println("下注失败，游戏还未开始")
		return
	}
	ret := this.check_end()
	if ret == 1 {
		this.turn_round()
	} else if ret == 2 {
		this.over_premature()
	} else {
		seat := this.next_chip_seat(this.m_attack_seat_id)
		if seat == nil {
			this.turn_round()
		} else {
			//通知下注
			max_chip := this.get_max_chip(seat.SeatID(), seat.RoundChip()+seat.SeatMoney())
			this.m_timer.Onec(this.cool_time, texas.TIME_CHIP_STATE)
			//
			println("开始下注:", this.m_chip_seat_id,
				"跟注=", this.m_need_call, ",最低加注=", this.m_need_call+this.m_need_chip,
				"最大下注=", max_chip, ",剩余=", seat.SeatMoney())
		}
	}
}

func (this *TexasLogic) next_chip_seat(attack_seat_id int) *texas.SeatPlayer {
	//当前下载玩家
	begin_id := this.m_chip_seat_id
	//移步下一个
	this.m_chip_seat_id = this.m_chip_seat_id + 1
	if this.m_chip_seat_id >= this.seat_count {
		this.m_chip_seat_id = 0
	}
	//查看其他人
	for i := 0; i < this.seat_count; i++ {
		idx := (this.m_chip_seat_id + i) % this.seat_count
		seat := this.seat_list[idx]
		//回到动作的玩家
		if seat.CheckSeatID(attack_seat_id) {
			if !seat.IsPlayer() || seat.RoundChip() == this.m_need_call {
				return nil
			}
		}
		if seat.IsActionPlayer() {
			this.m_chip_seat_id = idx
			break
		}
	}
	//无人操作
	if begin_id == this.m_chip_seat_id {
		return nil
	}
	return this.seat_list[this.m_chip_seat_id]
}

func (this *TexasLogic) is_playing() bool {
	return this.m_game_state > texas.GAME_STATE_READY && this.m_game_state < texas.GAME_START_STOP
}

//on_event
func (this *TexasLogic) OnChipAction(uid int32, atype int8, roundchip int32) {
	if seat, ok := this.get_seat_by_uid(uid); ok {
		this.ChipAction(int8(seat.SeatID()), atype, roundchip)
	}
}

func (this *TexasLogic) ChipAction(seat_id int8, atype int8, roundchip int32) {
	//是否游戏或者是否id一致
	if !this.is_playing() || int(seat_id) != this.m_chip_seat_id {
		println("游戏尚未开始或者还未轮到该玩家操作Warn:", seat_id)
		return
	}
	seat := this.get_chip_seat()
	if seat == nil {
		println("致命错误Error: 找不到当前操作的玩家 seat id=", seat_id)
		return
	}
	if atype == texas.SEAT_CHIP_RAISE {
		this.raise_chip(roundchip, seat)
		println(seat_id, "加注", seat.RoundChip(), " 剩余=", seat.SeatMoney())
	} else {
		if atype == texas.SEAT_CHIP_CALL {
			if this.m_need_call == 0 {
				//过牌
				println(seat_id, "过牌[无注跟]", " 剩余=", seat.SeatMoney())
			} else {
				seat.SubChip(this.m_need_call)
				println(seat_id, "跟注:", seat.RoundChip(), " 剩余=", seat.SeatMoney())
			}
		} else if atype == texas.SEAT_CHIP_FOLD {
			this.PushFoldMoney(seat.RoundChip())
			seat.SetFold()
			println(seat_id, "弃牌", " 剩余=", seat.SeatMoney())
		} else if atype == texas.SEAT_CHIP_CHECK {
			if this.m_need_call > seat.RoundChip() {
				this.PushFoldMoney(seat.RoundChip())
				seat.SetFold()
				println(seat_id, "弃牌[必须下注]", " 剩余=", seat.SeatMoney())
			} else {
				//过牌
				println(seat_id, "看牌", " 剩余=", seat.SeatMoney())
			}
		}
	}
	//println("玩家操作:", seat_id, ",类型=", atype, ",下注金额=", seat.RoundChip(), ",剩余=", seat.SeatMoney())
	this.ChipNextStart()
}

//加注
func (this *TexasLogic) raise_chip(roundchip int32, seat *texas.SeatPlayer) {
	total_chip := seat.RoundChip() + seat.SeatMoney()
	//总筹码不够加注（那么就全下）
	if total_chip < this.m_need_call+this.m_need_chip {
		seat.SubChip(total_chip)
		if total_chip > this.m_need_call {
			//加allin
			this.set_attack_seat(seat)
		} else if total_chip <= this.m_need_call {
			//跟注allin
		}
	} else {
		if roundchip >= this.m_need_call+this.m_need_chip {
			seat.SubChip(roundchip)
			this.set_attack_seat(seat)
		} else {
			//如果加注错误，那就默认最小
			seat.SubChip(this.m_need_call + this.m_need_chip)
			this.set_attack_seat(seat)
		}
	}
}

func (this *TexasLogic) set_attack_seat(seat *texas.SeatPlayer) {
	need_chip := seat.RoundChip() - this.m_need_call
	if need_chip > this.m_need_chip {
		this.m_need_chip = need_chip
	}
	this.m_need_call = seat.RoundChip()
	this.m_attack_seat_id = seat.SeatID()
}

//on action
func (this *TexasLogic) OnEnterTable(player *texas.PlayerVo) int {
	if this.m_current_player >= this.max_look {
		println("房间已满:", player.UserID, this.m_current_player, this.max_look)
		return 1
	}
	_, ok := this.players[player.UserID]
	if ok {
		println("该玩家在房间内:", player.UserID)
		return 2
	}
	println("进入房间:", player.UserID)
	this.m_current_player++
	this.players[player.UserID] = player
	//通知进入
	this.notice_players(command.CLIENT_ENTER_TEXAS_ROOM, gnet.NewByteArrayWithVals(player.UserID))
	return 0
}

func (this *TexasLogic) OnLeaveTable(uid int32) {
	player, ok := this.players[uid]
	if ok {
		this.m_current_player--
		delete(this.players, uid)
		//如果坐下了
		if seat, ok := this.get_seat_by_uid(uid); ok {
			this.stand_player(seat)
		}
		//通知离开
		this.notice_players(command.CLIENT_LEAVE_TEXAS_ROOM, gnet.NewByteArrayWithVals(player.UserID))
	}
}

func (this *TexasLogic) OnSitdown(uid int32, seat_id int8, seat_money int32) {
	//无座
	seat := this.get_seat(int(seat_id))
	if seat == nil {
		println("No seat: size=", this.seat_count, ",seat_id=", seat_id)
		return
	}
	data, ok := this.get_player(uid)
	if !ok {
		println("User is not in Table Err:", uid)
		return
	}
	//不在区间
	if this.min_buyin > seat_money || seat_money > this.max_buyin {
		println("Sitdown money is less Err: min=", this.min_buyin, ",max=", this.max_buyin, ", money=", seat_money)
		return
	}
	if seat.SitDown(seat_money, data) {
		this.notice_players(command.CLIENT_TEXAS_SITDOWN, gnet.NewByteArrayWithVals(uid, seat_id, seat_money))
		println("Sitdown OK: seat_id=", seat_id, ",seat_money=", seat_money, ",uid=", data.UserID)
		//开始游戏试试
		this.StartGame()
	} else {
		println("Sitdown Err is Sit")
	}
}

func (this *TexasLogic) OnStand(uid int32) {
	if seat, ok := this.get_seat_by_uid(uid); ok {
		this.stand_player(seat)
	}
}

func (this *TexasLogic) get_player(uid int32) (*texas.PlayerVo, bool) {
	data, ok := this.players[uid]
	return data, ok
}

//gets
func (this *TexasLogic) get_seat_by_uid(uid int32) (*texas.SeatPlayer, bool) {
	for _, seat := range this.seat_list {
		if seat.IsSit() && seat.Player.UserID == uid {
			return seat, true
		}
	}
	return nil, false
}

func (this *TexasLogic) get_max_chip(seat_id int, round int32) int32 {
	var big int32 = 0
	for _, seat := range this.seat_list {
		if !seat.CheckSeatID(seat_id) && seat.IsActionPlayer() {
			money := seat.RoundChip() + seat.SeatMoney()
			if money > big {
				big = money
			}
		}
	}
	if big > round {
		return round
	}
	return big
}

//主动站起
func (this *TexasLogic) stand_player(seat *texas.SeatPlayer) {
	this.StandSeat(seat)
	//是否游戏结束
	if this.is_playing() {
		this.PushFoldMoney(seat.RoundChip()) //将下注的钱直接注入池里面
		seat.SetFold()
		if seat.CheckSeatID(this.m_chip_seat_id) {
			this.ChipNextStart()
		} else {
			ret := this.check_end()
			if ret == 1 {
				this.turn_round()
			} else if ret == 2 {
				this.over_premature()
			}
		}
	}
}

func (this *TexasLogic) StandSeat(seat *texas.SeatPlayer) {
	this.notice_players(command.CLIENT_TEXAS_STAND, gnet.NewByteArrayWithVals(int8(seat.SeatID()), seat.Player.UserID))
	seat.Stand()
}

//0未结束，1轮次结束，2直接提前结束
//1,如果只有一个人，且其他玩家全部弃牌，直接结束
//2,如果只有一个人，有人allin，如果他need_call==0那么轮次结束
//3,如果所有人allin，那么轮次结束
func (this *TexasLogic) check_end() int {
	player_num, allin_num, seat := this.get_player_num()
	println("当前操作人数:", player_num, " 当前All-In人数:", allin_num)
	if player_num == 0 {
		//1人allin,其他人都弃牌
		if allin_num == 1 {
			return 2
		}
		//多人allin
		return 1
	} else if player_num == 1 {
		//1个人能操作,且其他人都弃牌
		if allin_num == 0 {
			return 2
		}
		//1个人能操作，并且这个人不需要跟注
		if seat.RoundChip() >= this.m_need_call {
			return 1
		}
	}
	return 0
}

//获取能操作的玩家和allin的玩家，返回一个能操作的玩家
func (this *TexasLogic) get_player_num() (int, int, *texas.SeatPlayer) {
	var player_num int = 0
	var allin_num int = 0
	var current_seat *texas.SeatPlayer = nil
	for _, seat := range this.seat_list {
		if seat.IsPlayer() {
			if seat.IsAllin() {
				allin_num++
			} else {
				player_num++
				current_seat = seat
			}
		}
	}
	return player_num, allin_num, current_seat
}

//超时默认操作
func (this *TexasLogic) outtime_chip() {
	this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_CHECK, 0)
}

func (this *TexasLogic) Test() {
	t := gutil.GetNano()
	for i := 0; i < 8; i++ {
		vo := new(texas.PlayerVo)
		vo.UserID = int32(1002 + i)
		vo.UserName = "test"
		//this.OnEnterTable(New)
		this.OnSitdown(vo.UserID, int8(i), 5000)
	}
	//游戏开始
	this.StartGame()
	for i := 0; i < 8; i++ {
		this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_CALL, 0)
	}
	//第二回合
	this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_CHECK, 0)
	this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_RAISE, this.big_blind)
	for i := 0; i < 7; i++ {
		this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_CALL, 0)
	}
	//第三回合
	this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_RAISE, 1000)
	//this.ChipAction(int8(this.m_chip_seat_id), SEAT_CHIP_FOLD, 0)
	for i := 0; i < 7; i++ {
		this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_CALL, 0)
	}
	//河牌
	this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_RAISE, 2000)
	for i := 0; i < 7; i++ {
		//跟
		this.ChipAction(int8(this.m_chip_seat_id), texas.SEAT_CHIP_FOLD, 0)
	}
	//this.ChipAction(int8(this.m_chip_seat_id),SEAT_CHIP_CALL,  0)
	println("耗时:", (gutil.GetNano()-t)/1000, "微秒")
}

//t := vat.GetWsTime()
//	//	for i := 0; i < 100; i++ {
//	//		p := texas.NewRobotCard()

//	//		list := make([]int16, 0)
//	//		for i := 0; i < 10; i++ {
//	//			list = append(list, p.NextCard())
//	//		}

//	//	}
//	list := []int16{0x0202, 0x0102, 0x0104, 0x0103, 0x0203, 0x010E, 0x0105}
//	ret, bigs := texas.CardTypeOfTexas(list)
//	println("耗时:", texas.POKER_TYPE_STR[ret], texas.TraceBigCards(bigs), vat.GetWsTime()-t, "微秒")
