package landlord

//牌结果
type CardResult struct {
	card_type int
	cards     []int16 //出的牌
	types     []int16 //牌型
	subs      []int16 //带的牌
}

func newResult() *CardResult {
	this := new(CardResult)
	//this.flush(ctype, cards)
	return this
}

func (this *CardResult) flush(ctype int, cards []int16) bool {
	this.card_type = ctype
	this.cards = cards
	return check_card_type(ctype, cards)
}

//比较>大于
func (this *CardResult) compare(result *CardResult) bool {
	if this.card_type > result.card_type {
		if this.card_type >= TYPE_BOMB_SOFT { //炸弹以上才行
			return true
		}
	} else if this.card_type == result.card_type {
		//比较最大值就可以了
		return Val(this.types[0]) > Val(result.types[0])
	}
	return false
}
