package texas

import "sort"

func init() {
	//	a := []int16{1, 3, 5, 7, 2}
	//	sort.Sort(a)
	//	fmt.Println(a)
}

type CardDataList []*CardData

//排序规则：首先按年龄排序（由小到大），年龄相同时按姓名进行排序（按字符串的自然顺序）
func (this CardDataList) Len() int {
	return len(this)
}

//这里使用降序
func (this CardDataList) Less(i, j int) bool {
	return this[i].value > this[j].value
}

func (this CardDataList) Swap(i, j int) {
	temp := this[i]
	this[i] = this[j]
	this[j] = temp
}

//每一张牌的数据
type CardData struct {
	color int8
	value int8
	m_num int16
}

func (this *CardData) IsAce() bool {
	return this.value == 0xE
}

func NewCardData(val int16) *CardData {
	this := new(CardData)
	this.color = int8(val >> 8)
	this.value = int8(val & 0xf)
	this.m_num = val
	return this
}

//class
type card_calculator struct {
	m_point    int
	big_list   [5]*CardData
	color_hash map[int8]CardDataList
	val_hash   map[int8]CardDataList
	sort_list  CardDataList
	card_size  int
}

//计算德州扑克的牌型
func CardTypeOfTexas(list []int16) (int, CardDataList) {
	this := new(card_calculator)
	this.m_point = 0
	this.card_size = len(list)
	this.sort_list = make(CardDataList, this.card_size)
	this.color_hash = make(map[int8]CardDataList, this.card_size)
	this.val_hash = make(map[int8]CardDataList, this.card_size)
	for i := 0; i < this.card_size; i++ {
		data := NewCardData(list[i])
		this.sort_list[i] = data
		add_hash(data.color, this.color_hash, data)
		add_hash(data.value, this.val_hash, data)
	}
	sort.Sort(this.sort_list)
	//	//
	//	for i := 0; i < this.card_size; i++ {
	//		print(PokerString(this.sort_list[i].m_num), ",")
	//	}
	//	print("\n")
	return this.GetType(), this.SupplyCards()
}

func add_hash(val int8, hash map[int8]CardDataList, data *CardData) {
	_, ok := hash[val]
	if !ok {
		hash[val] = make(CardDataList, 0)
	}
	hash[val] = append(hash[val], data)
}

func (this *card_calculator) GetType() int {
	ret := this.straight_flush()
	if ret == 2 { //皇家同花顺
		return TYPE_ROYAL_FLUSH
	} else if ret == 1 { //同花顺
		return TYPE_STRAIGHT_FLUSH
	}
	if this.FourKind() {
		return TYPE_FOUR_KIND
	} else if this.FullHouse() {
		return TYPE_FULL_HOUSE
	} else if this.Flush() {
		return TYPE_FLUSH
	} else if this.Straight() {
		return TYPE_STRAIGHT
	} else if this.ThreeKind() {
		return TYPE_THREE_KIND
	} else if this.TwoPair() {
		return TYPE_TWO_PAIRS
	} else if this.OnePair() {
		return TYPE_PAIR
	} else if this.TopKind() {
		return TYPE_HIGH_CARD
	}
	return TYPE_PIE_CARD
}

//同花顺
func (this *card_calculator) straight_flush() int {
	this.SetBegin()
	for _, list := range this.color_hash {
		if len(list) >= 5 {
			sort.Sort(list)
			ret := this.has_straight(list)
			if ret != nil {
				this.PushTypeCards(ret, 5)
				if ret[0].IsAce() {
					return 2
				}
				return 1
			}
		}
	}
	return 0
}

//4条
func (this *card_calculator) FourKind() bool {
	this.SetBegin()
	return this.SomePair(4, 0) > 0
}

//葫芦
func (this *card_calculator) FullHouse() bool {
	this.SetBegin()
	ret := this.SomePair(3, 0)
	if ret > 0 {
		if this.SomePair(2, ret) > 0 {
			return true
		}
	}
	return false
}

//获取最大的同花(目测最多一个同花)
func (this *card_calculator) Flush() bool {
	this.SetBegin()
	var big_val int8 = 0
	var key int8 = 0
	for color, list := range this.color_hash {
		if len(list) >= 5 {
			if list[0].value > big_val {
				key = color
			}
		}
	}
	//最大同花
	if key > 0 {
		bigs := this.color_hash[key]
		sort.Sort(bigs)
		this.PushTypeCards(bigs[0:5], 5)
	}
	return key > 0
}

//顺子
func (this *card_calculator) Straight() bool {
	this.SetBegin()
	ret := this.has_straight(this.sort_list)
	if ret != nil {
		this.PushTypeCards(ret, 5)
	}
	return ret != nil
}

//获取最大的顺子(A 2 3 4 5也是顺子啊)
func (this *card_calculator) has_straight(list CardDataList) CardDataList {
	var begin_val int8 = 0
	var flow_num int = 0
	var str_list CardDataList = make(CardDataList, 5)
	for _, data := range list {
		if begin_val == data.value {
			continue
		} else if begin_val-1 == data.value {
			flow_num++
		} else {
			flow_num = 1
		}
		str_list[flow_num-1] = data
		begin_val = data.value
		if flow_num == 5 {
			break
		}
	}
	//特需情况，最小的顺子(最后一张是2且最大牌有A)
	if flow_num == 4 {
		if str_list[flow_num-1].value == 2 && list[0].IsAce() {
			flow_num++
			str_list[flow_num-1] = list[0]
		}
	}
	if flow_num == 5 {
		return str_list[0:5]
	}
	return nil
}

func (this *card_calculator) ThreeKind() bool {
	this.SetBegin()
	return this.SomePair(3, 0) > 0
}

func (this *card_calculator) TwoPair() bool {
	this.SetBegin()
	ret := this.SomePair(2, 0)
	if ret > 0 {
		if this.SomePair(2, ret) > 0 {
			return true
		}
	}
	return false
}

func (this *card_calculator) OnePair() bool {
	this.SetBegin()
	return this.SomePair(2, 0) > 0
}

//获取最大相同张数(返回一个值)
func (this *card_calculator) SomePair(size int, del int8) int8 {
	var big int8 = 0
	for val, list := range this.val_hash {
		//剔除一个
		if del == val {
			continue
		}
		if len(list) >= size {
			if val > big {
				big = val
			}
		}
	}
	if big > 0 {
		this.PushTypeCards(this.val_hash[big], size)
	}
	return big
}

func (this *card_calculator) TopKind() bool {
	this.SetBegin()
	//直接获取顶部5张
	var ret bool = false
	for val, _ := range this.val_hash {
		if val > 9 {
			ret = true
		}
	}
	//直接获取顶部5张(不管了)
	this.PushTypeCards(this.sort_list[0:5], 5)
	return ret
}

//最大牌型
func (this *card_calculator) SetBegin() {
	this.m_point = 0
}

func (this *card_calculator) PushTypeCards(list CardDataList, size int) {
	for i := 0; i < size; i++ {
		this.PushTypeCard(list[i])
	}
}

func (this *card_calculator) PushTypeCard(card *CardData) {
	this.big_list[this.m_point] = card
	this.m_point++
}

func (this *card_calculator) AppendBigCard(card *CardData) {
	var ret bool = false
	for i := 0; i < this.m_point; i++ {
		if this.big_list[i] == card {
			ret = true
			break
		}
	}
	if !ret {
		this.PushTypeCard(card)
	}
}

//根据已经给出的牌，找出剩下最大的牌补全5张(德州扑克专用)
func (this *card_calculator) SupplyCards() CardDataList {
	if this.m_point < POCKET_CARD_SIZE {
		//补全剩余的牌
		for _, data := range this.sort_list {
			this.AppendBigCard(data)
			if this.m_point >= POCKET_CARD_SIZE {
				break
			}
		}
	}
	return this.big_list[0:POCKET_CARD_SIZE]
}

//比较5组牌(同类型比较) 1[a>b],0[a==b],-1[a<b]
func CompareCards(lsa CardDataList, lsb CardDataList) int {
	//最简单的办法就是，值相加
	for i := 0; i < POCKET_CARD_SIZE; i++ {
		if lsa[i].value > lsb[i].value {
			return 1
		} else if lsa[i].value < lsb[i].value {
			return -1
		}
	}
	return 0
}

func TraceBigCards(list CardDataList) string {
	str := "最大牌:"
	for i := 0; i < POCKET_CARD_SIZE; i++ {
		str += PokerString(list[i].m_num) + ","
	}
	return str
}
