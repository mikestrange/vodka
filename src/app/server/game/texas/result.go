package texas

type SeatResultList []*SeatResult

type SeatResult struct {
	Seat     *SeatPlayer
	PotIdx   int
	CardType int
	CardList CardDataList
}

func NewSeatResult(seat *SeatPlayer, cards []int16) *SeatResult {
	this := new(SeatResult)
	this.Seat = seat
	this.PotIdx = seat.m_pot_point
	this.CardType, this.CardList = CardTypeOfTexas(cards)
	return this
}

func (this *SeatResult) Trace() {
	println("玩家:", this.Seat.SeatID(), POKER_TYPE_STR[this.CardType], TraceBigCards(this.CardList))
}

//排序
func (this SeatResultList) Len() int {
	return len(this)
}

//降序
func (this SeatResultList) Less(i, j int) bool {
	return this[i].CardType > this[j].CardType
}

func (this SeatResultList) Swap(i, j int) {
	temp := this[i]
	this[i] = this[j]
	this[j] = temp
}
