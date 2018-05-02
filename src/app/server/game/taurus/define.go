package taurus

type CardList []int16

//0>无牛
//10>牛牛
//11>银牛	五张牌全由10～k组成且只有一张10，例如10、j、j、q、k。
//12>金牛	五张牌全由j～k组成，例如j、j、q、q、k。
//13>炸弹	五张牌中有4张牌点数相同的牌型，例如：2、2、2、2、k。
//14>五小牛	(目前不支持)五张牌的点数加起来小于10，且每张牌点数都小于5，例如a、3、2、a、2。

const (
	WIN  = 0
	LOSE = 1
	DRAW = 2
	//
	TYPE_NONIU  = 0
	TYPE_BEINIU = 6 //大于他就翻倍
	TYPE_NIUNIU = 10
	TYPE_YINNIU = 11
	TYPE_JINNIU = 12
	TYPE_BOMB   = 13
	TYPE_SMALL  = 14 //目前不存在
)

const CARD_SIZE = 52

var cardTypes = []string{"无牛", "牛1", "牛2", "牛3", "牛4", "牛5", "牛6", "牛7", "牛8", "牛9", "牛牛", "银牛", "金牛", "炸弹", "五小牛"}
var CARD_COLORS []string = []string{"-", "方片", "梅花", "红桃", "黑桃"}
var CARD_VALS []string = []string{"-", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

var CARD_LIST = [52]int16{
	0x101, 0x102, 0x103, 0x104, 0x105, 0x106, 0x107, 0x108, 0x109, 0x10A, 0x10B, 0x10C, 0x10D,
	0x201, 0x202, 0x203, 0x204, 0x205, 0x206, 0x207, 0x208, 0x209, 0x20A, 0x20B, 0x20C, 0x20D,
	0x301, 0x302, 0x303, 0x304, 0x305, 0x306, 0x307, 0x308, 0x309, 0x30A, 0x30B, 0x30C, 0x30D,
	0x401, 0x402, 0x403, 0x404, 0x405, 0x406, 0x407, 0x408, 0x409, 0x40A, 0x40B, 0x40C, 0x40D}

//一张
func CardStr(val int16) string {
	k := int(val >> 8)
	c := int(val & 0xf)
	return CARD_COLORS[k] + CARD_VALS[c]
}

//多张
func ToStr(args ...int16) string {
	var str string = ""
	for i := range args {
		str += CardStr(args[i]) + ","
	}
	return str
}

func CardVal(num int16) int {
	return int(num & 0xf)
}

func CardColor(num int16) int {
	return int(num >> 8)
}

func TypeStr(mtype int) string {
	return cardTypes[mtype]
}
