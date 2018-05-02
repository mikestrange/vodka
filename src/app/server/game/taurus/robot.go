package taurus

import "math/rand"

//发牌机器人(顶部抽取)
type CardDealer struct {
	cards []int16
	pos   int
	size  int
}

func NewDealer() *CardDealer {
	this := new(CardDealer)
	this.init()
	return this
}

func (this *CardDealer) init() {
	this.size = CARD_SIZE
	this.cards = make([]int16, this.size)
	copy(this.cards[0:], CARD_LIST[0:this.size])
}

//随机
func (this *CardDealer) Random() {
	this.pos = this.size
	for m := 0; m < this.size; m++ {
		for i := 0; i < this.size; i++ {
			r := rand.Int() % this.size
			temp := this.cards[i]
			this.cards[i] = this.cards[r]
			this.cards[r] = temp
		}
	}
}

func (this *CardDealer) Pop() int16 {
	if this.Available() == 0 {
		return 0
	}
	this.pos--
	return this.cards[this.pos]
}

func (this *CardDealer) Available() int {
	return this.pos
}
