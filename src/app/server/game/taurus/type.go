package taurus

import "ants/base"

type CardDataList []*CardData

//每一张牌的数据
type CardData struct {
	num   int16 //原始
	color int8
	value int8
}

func NewData(val int16) *CardData {
	return &CardData{val, int8(val >> 8), int8(val & 0xf)}
}

func (this *CardData) Num(val int16) {
	this.num = val
}

func (this *CardData) Val() int {
	if this.value > 10 {
		return 10
	}
	return int(this.value)
}

//true a1大
func CompareCard(a1 int16, a2 int16) int {
	if a1 == a2 {
		return DRAW
	}
	v1 := CardVal(a1)
	v2 := CardVal(a2)
	if v1 == v2 {
		if CardColor(a1) > CardColor(a2) {
			return WIN
		}
		return LOSE
	}
	if v1 > v2 {
		return WIN
	}
	return LOSE
}

//判断是否一样
func check_some_vals(args ...int8) bool {
	var val int8 = -1
	for i := range args {
		if val == -1 {
			val = args[i]
		} else {
			if args[i] != val {
				return false
			}
		}
	}
	return true
}

//炸弹
func typeBomb(c CardDataList) bool {
	if check_some_vals(c[0].value, c[1].value, c[2].value, c[3].value) ||
		check_some_vals(c[4].value, c[1].value, c[2].value, c[3].value) {
		return true
	}
	return false
}

//金牛
func typeTaurus(c CardDataList) bool {
	return c[1].value > 10 && c[2].value > 10 && c[3].value > 10 && c[4].value > 10 && c[0].value > 10
}

//银牛
func typeSilverbull(c CardDataList) bool {
	return c[1].value >= 10 && c[2].value >= 10 && c[3].value >= 10 && c[4].value >= 10 && c[0].value >= 10
}

//其他类型
func typeOther(cards CardDataList) int {
	ctype := 0
	size := len(cards)
	for i := 0; i < size; i++ {
		//前3位
		c1 := cards[i]
		c2 := cards[(i+1)%size]
		c3 := cards[(i+2)%size]
		//后两位
		c4 := cards[(i+3)%size]
		c5 := cards[(i+4)%size]
		//println(c1.Val(), c2.Val(), c3.Val(), c4.Val(), c5.Val())
		sum := c1.Val() + c2.Val() + c3.Val()
		if num := sum % 10; num == 0 {
			//后两位
			mtype := (c4.Val() + c5.Val()) % 10
			//牛牛特需
			if mtype == 0 {
				mtype = TYPE_NIUNIU
			}
			if mtype > ctype {
				ctype = mtype
			}
		}
	}
	return ctype
}

//5张牌 按大小排序
func getNiuNiuType(cards []int16) int {
	//最大的排前面
	base.Sort(cards, func(i, j int) bool {
		return CardVal(cards[i]) > CardVal(cards[j])
	})
	var list CardDataList
	for i := 0; i < len(cards); i++ {
		list = append(list, NewData(cards[i]))
	}
	//----
	if typeBomb(list) {
		return TYPE_BOMB
	} else if typeTaurus(list) {
		return TYPE_JINNIU
	} else if typeSilverbull(list) {
		return TYPE_YINNIU
	} else {
		return typeOther(list)
	}
	return TYPE_NONIU
}

//func init() {
//	arr := []int16{0x201, 0x303, 0x405, 0x102, 0x10b}
//	t := getNiuNiuType(arr)
//	println("牛几:", t)
//	println(ToStr(arr...))
//}
