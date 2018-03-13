package texas

//座位列表
type SeatPlayerList []*SeatPlayer

//排序
func (this SeatPlayerList) Len() int {
	return len(this)
}

//升序
func (this SeatPlayerList) Less(i, j int) bool {
	return this[i].m_init_chip < this[j].m_init_chip
}

func (this SeatPlayerList) Swap(i, j int) {
	temp := this[i]
	this[i] = this[j]
	this[j] = temp
}

//0 - max
type SeatPlayer struct {
	m_seat_id int
	m_is_sit  bool
	m_state   int
	auto_buy  bool
	Player    *PlayerVo //玩家信息
	//
	m_pot_point    int   //参与池
	m_round_chip   int32 //回合下注
	m_current_chip int32 //当前加注
	m_init_chip    int32 //每回合带入的金币
	m_seat_money   int32 //当前坐下金币
	m_cards        [2]int16
}

func NewTexasSeatPlayer(idx int) *SeatPlayer {
	this := new(SeatPlayer)
	this.InitSeat(idx)
	return this
}

func (this *SeatPlayer) InitSeat(idx int) {
	this.m_seat_id = idx
	this.m_state = PLAYER_READY_STATE
	this.Stand()
}

func (this *SeatPlayer) SeatID() int {
	return this.m_seat_id
}

func (this *SeatPlayer) CheckSeatID(idx int) bool {
	return this.m_seat_id == idx
}

func (this *SeatPlayer) SeatMoney() int32 {
	return this.m_seat_money
}

func (this *SeatPlayer) RoundChip() int32 {
	return this.m_round_chip
}

func (this *SeatPlayer) SetAutoBuy(mtype int) {
	this.auto_buy = mtype == 1
}

func (this *SeatPlayer) IsAutoBuy() bool {
	return this.auto_buy
}

//游戏状态
func (this *SeatPlayer) Begin() {
	this.m_state = PLAYER_PLAYING
	this.m_init_chip = this.m_seat_money
	this.m_current_chip = 0
	this.m_round_chip = 0
	this.m_pot_point = 0
}

func (this *SeatPlayer) Over() {
	this.m_state = PLAYER_READY_STATE
	this.m_round_chip = 0
	this.m_current_chip = 0
}

func (this *SeatPlayer) Turn() {
	this.m_round_chip = 0
	this.m_current_chip = 0
}

func (this *SeatPlayer) SubTableFree(val int32) {
	if val > this.m_seat_money {
		println("服务费不足:", this.m_seat_money, val)
		return
	}
	this.m_seat_money -= val
}

func (this *SeatPlayer) SitDown(money int32, data *PlayerVo) bool {
	if this.IsSit() {
		return false
	}
	this.m_is_sit = true
	this.m_state = PLAYER_READY_STATE
	this.m_seat_money = money
	this.Player = data
	return true
}

func (this *SeatPlayer) Stand() {
	this.m_is_sit = false
	this.Player = nil
	this.m_state = PLAYER_READY_STATE
}

func (this *SeatPlayer) SetPot(val int) {
	this.m_pot_point = val
}

func (this *SeatPlayer) CheckUserID(uid int32) bool {
	if this.IsSit() {
		return this.Player.UserID == uid
	}
	return false
}

//cards
func (this *SeatPlayer) SetCard(idx int, val int16) {
	this.m_cards[idx] = val
}

func (this *SeatPlayer) GetCard(idx int) int16 {
	return this.m_cards[idx]
}

//返回真实下注的筹码
func (this *SeatPlayer) SubChip(val int32) int32 {
	if this.m_round_chip > val {
		println("this seat chip is less to:", val, this.m_round_chip)
		return 0
	}
	//加注的值(扣)
	var chip int32 = val - this.m_round_chip
	//allin
	if chip >= this.m_seat_money {
		chip = this.m_seat_money
		this.SetAllin()
	}
	this.m_seat_money -= chip
	this.m_round_chip = this.m_round_chip + chip
	this.m_current_chip = chip
	return this.m_round_chip
}

//结果赢得
func (this *SeatPlayer) Result(val int64) {
	this.m_seat_money += int32(val)
}

//action
func (this *SeatPlayer) SetFold() {
	this.m_state = PLAYER_FOLD_CARD
	this.m_current_chip = 0
	this.m_round_chip = 0
}

func (this *SeatPlayer) SetAllin() {
	this.m_state = PLAYER_ALLIN
}

//is
func (this *SeatPlayer) IsPlayer() bool {
	if this.IsStand() || this.m_state == PLAYER_READY_STATE {
		return false
	}
	if this.IsFold() {
		return false
	}
	return true
}

func (this *SeatPlayer) IsActionPlayer() bool {
	return this.IsPlayer() && !this.IsAllin()
}

func (this *SeatPlayer) IsAllin() bool {
	return this.m_state == PLAYER_ALLIN
}

func (this *SeatPlayer) IsFold() bool {
	return this.m_state == PLAYER_FOLD_CARD
}

func (this *SeatPlayer) IsStand() bool {
	return this.m_is_sit == false
}

func (this *SeatPlayer) IsSit() bool {
	return this.m_is_sit
}
