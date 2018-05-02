package landlord

//每抢地主 *2 、每炸弹 *2 、明牌 *2-5 、春天*2。

type BaseGame struct {
	//通用动态
	banker_idx    int     //地主
	attack_idx    int     //当前玩家(包括抢庄)
	current_idx   int     //当前操作
	call_num      int     //叫的
	game_state    int     //游戏状态
	game_multiple int     //当前游戏倍数
	public_cards  []int16 //公共牌
	//房间数据
	table_id       int
	table_times    int //房间场次
	table_type     int //癞子，普通，快速等
	table_free     int //服务费
	seat_count     int //人数:默认3人
	oper_time      int //操作时间
	base_chip      int //基础分
	min_buyin      int //最小入座金币
	min_player_num int //最小开始人数(必须3人)
	max_look       int //观看人数限制
	//0000
	seats []*Seat
	users map[int]*Player
}

func (this *BaseGame) init() {
	this.seats = make([]*Seat, this.seat_count)
	for i := 0; i < this.seat_count; i++ {
		this.seats[i] = newSeat(i + 1)
	}
	this.users = make(map[int]*Player)
}

func (this *BaseGame) length() int {
	return len(this.users)
}

func (this *BaseGame) setPlayer(uid int, player *Player) {
	this.users[uid] = player
}

func (this *BaseGame) SetState(val int) {
	this.game_state = val
}

func (this *BaseGame) check_state(val int) bool {
	return this.game_state == val
}

func (this *BaseGame) check_state_set(val int) bool {
	if this.game_state == val {
		return false
	}
	this.game_state = val
	return true
}

func (this *BaseGame) get_seat(idx int) (*Seat, bool) {
	if idx > this.seat_count || idx < 1 {
		return nil, false
	}
	return this.seats[idx-1], true
}

func (this *BaseGame) each_seats(f func(*Seat)) {
	for i := range this.seats {
		f(this.seats[i])
	}
}

func (this *BaseGame) each_sits(f func(*Seat)) {
	for i := range this.seats {
		if this.seats[i].issit() {
			f(this.seats[i])
		}
	}
}

func (this *BaseGame) get_sit_num() int {
	num := 0
	this.each_sits(func(seat *Seat) {
		num++
	})
	return num
}

func (this *BaseGame) each_players(f func(*Player)) {
	for i := range this.users {
		f(this.users[i])
	}
}

func (this *BaseGame) get_player(uid int) (*Player, bool) {
	p, ok := this.users[uid]
	if ok {
		return p, true
	}
	return nil, false
}

func (this *BaseGame) getSeatByUid(uid int) (*Seat, bool) {
	for i := range this.seats {
		if this.seats[i].uid == uid {
			return this.seats[i], true
		}
	}
	return nil, false
}

func (this *BaseGame) has(uid int) bool {
	if _, ok := this.users[uid]; ok {
		return true
	}
	return false
}

func (this *BaseGame) delUser(uid int) (*Player, bool) {
	if player, ok := this.get_player(uid); ok {
		delete(this.users, uid)
		return player, true
	}
	return nil, false
}

func (this *BaseGame) isplaying() bool {
	return this.game_state > GAME_WAIT && this.game_state < GAME_STOP
}
