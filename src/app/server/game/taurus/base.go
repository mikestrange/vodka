package taurus

const (
	GAME_WAIT   = 0
	GAME_START  = 1
	GAME_CHIP   = 2
	GAME_COMMIT = 3
	GAME_STOP   = 4
)

//固定参数
type BaseGame struct {
	//活动参数
	game_state      int //1开始，2下倍，3提交，4结束，0等待开始
	banker_idx      int //专家id
	banker_multiple int //庄倍数
	//
	table_id    int
	table_type  int //类型
	table_num   int //场次
	table_free  int //台费
	seat_count  int //座位数
	banker_time int //抢庄时间
	chip_time   int //下注时间
	commit_time int //提交时间
	over_time   int //结束时间
	base_chip   int //基数
	min_money   int //最小带入（座位）
	max_money   int //最大带入（座位）
	min_chip    int //最小下注倍数
	max_chip    int //最大下注倍数
	min_player  int //最小开始人数
	max_look    int //最大人数
	//其他
	users map[int]*Player
	seats []*Seat
}

func (this *BaseGame) init() {
	this.users = make(map[int]*Player)
	this.seats = make([]*Seat, this.seat_count)
	for i := 0; i < this.seat_count; i++ {
		this.seats[i] = newSeat(i + 1)
	}
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
	for i := range this.seats {
		if this.seats[i].issit() {
			num++
		}
	}
	return num
}

func (this *BaseGame) each_actions(f func(*Seat)) {
	for i := range this.seats {
		if this.seats[i].isplayer() {
			f(this.seats[i])
		}
	}
}

func (this *BaseGame) get_action_num() int {
	num := 0
	for i := range this.seats {
		if this.seats[i].isplayer() {
			num++
		}
	}
	return num
}

func (this *BaseGame) find_with(f func(*Seat) bool) (*Seat, bool) {
	for i := range this.seats {
		seat := this.seats[i]
		if f(seat) {
			return seat, true
		}
	}
	return nil, false
}

func (this *BaseGame) each_pos(pos int, f func(*Seat)) {
	for i := 0; i < this.seat_count; i++ {
		f(this.seats[(pos+i)%this.seat_count])
	}
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
