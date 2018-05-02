package texas

import "math/rand"

//发牌机器人(顶部抽取)
type RobotCard struct {
	cards [52]int16
	pos   int
}

func newCards() *RobotCard {
	this := new(RobotCard)
	this.init()
	return this
}

func (this *RobotCard) init() {
	copy(this.cards[0:], GLOBAL_POKER[0:])
	this.Random()
}

func (this *RobotCard) Random() {
	this.pos = POKER_SIZE
	for m := 0; m < POKER_SIZE; m++ {
		for i := 0; i < POKER_SIZE; i++ {
			r := rand.Int() % POKER_SIZE
			temp := this.cards[i]
			this.cards[i] = this.cards[r]
			this.cards[r] = temp
		}
	}
}

func (this *RobotCard) Pop() int16 {
	if this.Available() == 0 {
		return 0
	}
	this.pos--
	return this.cards[this.pos]
}

func (this *RobotCard) Available() int {
	return this.pos
}
