package texas

import "math/rand"

//发牌机器人(顶部抽取)
type RobotCard struct {
	m_cards [52]int16
	m_pos   int
}

func NewRobotCard() *RobotCard {
	this := new(RobotCard)
	this.InitRobotCard()
	return this
}

func (this *RobotCard) InitRobotCard() {
	copy(this.m_cards[0:], GLOBAL_POKER[0:])
	this.RandomCards()
}

func (this *RobotCard) RandomCards() {
	this.m_pos = POKER_SIZE
	for m := 0; m < POKER_SIZE; m++ {
		for i := 0; i < POKER_SIZE; i++ {
			r := rand.Int() % POKER_SIZE
			temp := this.m_cards[i]
			this.m_cards[i] = this.m_cards[r]
			this.m_cards[r] = temp
		}
	}
}

func (this *RobotCard) NextCard() int16 {
	if this.CardAvailable() == 0 {
		return 0
	}
	this.m_pos--
	return this.m_cards[this.m_pos]
}

func (this *RobotCard) CardAvailable() int {
	return this.m_pos
}

func (this *RobotCard) Trace() {
	for i := 0; i < 52; i++ {
		print(PokerString(this.m_cards[i]), ",")
	}
	print("\n")
}
