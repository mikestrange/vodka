package texas

//一些基础的操作(无关逻辑)

type BaseGame struct {
	//通用动态
	chip_idx   int //当前操作的玩家
	attack_idx int //攻击的位置 -1为没有攻击
	round_idx  int //1个开始的位置(没有攻击位置，回到本位置就下一回合)
	game_state int //游戏状态
	//房间数据
	room_id        int
	table_times    int //房间场次
	room_type      int
	seat_count     int
	bet_time       int //下注时间
	table_free     int //服务费
	small_blind    int
	big_blind      int
	min_buyin      int //最小带入金额
	max_buyin      int //最大带入金额
	max_look       int //观看人数限制
	min_player_num int //最小开始人数
	//
	seats []*Seat
	users map[int]*Player
	//
	public_cards []int16
}

func (this *BaseGame) init() {
	//seats 第一个座位为空位 1 - size
	this.seats = make([]*Seat, this.seat_count)
	for i := 0; i < this.seat_count; i++ {
		this.seats[i] = newSeat(i + 1)
	}
	//users
	this.users = make(map[int]*Player)
	//cards
	this.public_cards = make([]int16, 0, 5)
}

func (this *BaseGame) pushCard(num int16) {
	this.public_cards = append(this.public_cards, num)
}

func (this *BaseGame) clearCards() {
	if len(this.public_cards) > 0 {
		this.public_cards = make([]int16, 0, 5)
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

func (this *BaseGame) get_seat(seat_id int) (*Seat, bool) {
	if seat_id > this.seat_count || seat_id < 1 {
		return nil, false
	}
	return this.seats[seat_id-1], true
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

func (this *BaseGame) get_player(uid int) (*Player, bool) {
	user, ok := this.users[uid]
	return user, ok
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
	return this.game_state > GAME_STATE_WAIT && this.game_state < GAME_START_STOP
}

func (this *BaseGame) each_players(f func(*Player)) {
	for i := range this.users {
		f(this.users[i])
	}
}

//所有在玩的玩家
func (this *BaseGame) each_actions(f func(*GameAction)) {
	for i := range this.seats {
		if this.seats[i].isPlaying() {
			f(this.seats[i].Action())
		}
	}
}

//攻击
func (this *BaseGame) hasAttack() bool {
	return this.attack_idx != -1
}

func (this *BaseGame) setAttack(idx int) {
	this.attack_idx = idx
}

func (this *BaseGame) setNoAttack() {
	this.attack_idx = -1
}

func (this *BaseGame) checkAttack(idx int) bool {
	return this.attack_idx == idx
}

func (this *BaseGame) setRound(idx int) {
	this.round_idx = idx
}

func (this *BaseGame) setNoRound() {
	this.round_idx = -1
}

func (this *BaseGame) checkRoundset(idx int) bool {
	if this.round_idx == -1 {
		this.round_idx = idx
		return false
	}
	return this.round_idx == idx
}

//当前操作
func (this *BaseGame) setCurrent(idx int) {
	this.chip_idx = idx
}

func (this *BaseGame) isCurrent(idx int) bool {
	return this.chip_idx == idx
}

//从pos(庄)开始找到一个活人:pos是屏蔽的那一位
func (this *BaseGame) find_action(pos int) (*GameAction, bool) {
	for i := 0; i < this.seat_count; i++ {
		idx := (pos + i) % this.seat_count
		seat := this.seats[idx]
		if act, ok := seat.checkAction(); ok {
			if act.isAction() {
				return act, true
			}
		}
	}
	return nil, false
}

//当前操作的玩家
func (this *BaseGame) current_action() (*GameAction, bool) {
	seat, ok := this.get_seat(this.chip_idx)
	if ok && seat.isPlaying() {
		return seat.Action(), true
	}
	return nil, false
}

//获取下一个游戏人
func (this *BaseGame) find_player(pos int) (*GameAction, bool) {
	for i := 0; i < this.seat_count; i++ {
		idx := (pos + i) % this.seat_count
		seat := this.seats[idx]
		if act, ok := seat.checkAction(); ok {
			return act, true
		}
	}
	return nil, false
}

//是否就他未弃牌
func (this *BaseGame) get_only_player() (*GameAction, bool, bool) {
	player_num := 0 //能行动的人数
	allin_num := 0  //allin人数
	var cur *GameAction
	for i := 0; i < this.seat_count; i++ {
		//idx := (pos + i) % this.seat_count
		seat := this.seats[i]
		if act, ok := seat.checkAction(); ok {
			if act.isNoFold() {
				player_num++ //包括allin和活着的人
				cur = act
			}
			if act.isAllin() {
				allin_num++ //allin人数
			}
		}
	}
	//活人
	ret := player_num - allin_num
	if ret == 0 { //无人可操作
		if allin_num > 1 { //多人allin
			return nil, false, true
		} else { //一人allin
			return cur, true, false
		}
	} else if ret == 1 { //一人可操作
		if allin_num == 0 { //无人allin
			return cur, true, false
		} else { //有人allin
			return nil, false, true
		}
	}
	return nil, false, false
}

//获取下一个坐下的玩家
func (this *BaseGame) get_next_siter(pos int) int {
	for i := 0; i < this.seat_count; i++ {
		idx := (pos + i) % this.seat_count
		if this.seats[idx].issit() {
			return this.seats[idx].seat_id
		}
	}
	panic("not next siter")
}

func (this *BaseGame) check_buyin(val int) bool {
	return val >= this.min_buyin && val <= this.max_buyin
}
