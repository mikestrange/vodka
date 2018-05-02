package landlord

type GameAction struct {
	uid       int
	seat_id   int
	open_card bool //明牌
	state     int  //状态(脱管,离线)
	//2倍
	banker_call int //抢庄
	double_call int //加倍
	//手里牌
	cards    []int16 //手牌(出过的牌置后)
	card_pos int     //当前位置
	//牌型
	result *CardResult
	//其他
	timeout_num int //超时次数
}

func newAction(seat *Seat) *GameAction {
	this := new(GameAction)
	this.uid = seat.uid
	this.seat_id = seat.seat_id
	return this
}

func (this *GameAction) init() {
	this.card_pos = 0
	this.timeout_num = 0
	this.result = nil
}

func (this *GameAction) pushCard(val int16) {
	this.cards = append(this.cards, val)
	this.card_pos++
}

//明牌
func (this *GameAction) OpenCards() {
	this.open_card = true
}

//叫庄
func (this *GameAction) Call() {
	this.banker_call++
}

//加倍
func (this *GameAction) Double() {
	this.double_call = 1
}

//是否大于
func (this *GameAction) BigTo(action *GameAction, ctype int, cards []int16) bool {
	if !this.indexOfs(cards) {
		return false
	}
	ret := newResult()
	if ret.flush(ctype, cards) && ret.compare(action.result) {
		this.result = ret
		this.killCards(ret.cards)
		return true
	}
	return false
}

//删除这些牌
func (this *GameAction) killCards(cards []int16) {
	for i := range cards {
		this.killCard(cards[i])
	}
}

func (this *GameAction) indexOf(card int16) int {
	for i := 0; i < this.card_pos; i++ {
		if this.cards[i] == card {
			return i
		}
	}
	return -1
}

func (this *GameAction) indexOfs(cards []int16) bool {
	for i := range cards {
		if this.indexOf(cards[i]) == -1 {
			return false
		}
	}
	return true
}

func (this *GameAction) killCard(card int16) {
	idx := this.indexOf(card)
	if idx != -1 {
		temp := this.cards[this.card_pos-1]
		this.cards[this.card_pos-1] = card
		this.cards[idx] = temp
		this.card_pos = this.card_pos - 1
	}
}

//是否出完了
func (this *GameAction) overShow() bool {
	return this.card_pos == 0
}

//出牌
func (this *GameAction) outCards(ctype int, cards []int16) bool {
	if !this.indexOfs(cards) {
		return false
	}
	ret := newResult()
	if ret.flush(ctype, cards) {
		//是组合
		return true
	}
	return false
}
