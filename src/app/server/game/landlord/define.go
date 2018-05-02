package landlord

//游戏状态
const (
	GAME_WAIT   = 0
	GAME_BANKER = 1 //抢庄
	GAME_DOUBLE = 2 //加倍
	GAME_PLAY   = 3 //操作
	GAME_STOP   = 4 //结束
)

//牌型
const (
	TYPE_NIL             = 0  // 不能
	TYPE_ONE             = 1  // 单张
	TYPE_PAIR            = 2  // 对牌
	TYPE_THREE           = 3  // 三张
	TYPE_THREE_TO_ONE    = 4  // 三带一
	TYPE_THREE_TO_PAIR   = 5  // 三带对
	TYPE_ONE_S           = 6  // 顺子
	TYPE_PAIR_S          = 7  // 连对
	TYPE_THREE_S         = 8  // 飞机不带
	TYPE_THREE_TO_ONE_S  = 9  // 飞机带一
	TYPE_THREE_TO_PAIR_S = 10 // 飞机带对
	TYPE_FOUR_TO_ONE_S   = 11 // 四带一
	TYPE_FOUR_TO_PAIR_S  = 12 // 四带对
	TYPE_BOMB_SOFT       = 13 // 软癞炸弹(del)
	TYPE_BOMB_HARD       = 14 // 硬炸弹
	TYPE_BOMB_LAZI       = 15 // 纯癞子炸弹(del)
	TYPE_ROCKETS         = 16 // 火箭
)

var TYPE_STR []string = []string{"单张", "对牌", "三张", "三带单", "三带对",
	"单顺", "双顺", "飞机不带", "飞机带单", "飞机带对",
	"四带一", "四带对",
	"软炸弹", "硬炸弹", "变态炸弹*纯癞子", "火箭"}

var CARD_COLORS []string = []string{"-", "方片", "梅花", "红桃", "黑桃", "鬼"}
var CARD_VALS []string = []string{"-", "A", "2", "3", "4", "5",
	"6", "7", "8", "9", "10", "J", "Q", "K", "A", "-", "2", "-", "小鬼", "-", "大鬼"}

const CARD_SIZE = 54
const SmallKing = 0x513
const BigKing = 0x515

//为了让A,2,小王,大王不连续,中间被隔开(11表示2, 13小王, 15大王)
var CARD_LIST [54]int16 = [54]int16{
	0x103, 0x104, 0x105, 0x106, 0x107, 0x108, 0x109, 0x10A, 0x10B, 0x10C, 0x10D, 0x10E, 0x111,
	0x203, 0x204, 0x205, 0x206, 0x207, 0x208, 0x209, 0x20A, 0x20B, 0x20C, 0x20D, 0x20E, 0x211,
	0x303, 0x304, 0x305, 0x306, 0x307, 0x308, 0x309, 0x30A, 0x30B, 0x30C, 0x30D, 0x30E, 0x311,
	0x403, 0x404, 0x405, 0x406, 0x407, 0x408, 0x409, 0x40A, 0x40B, 0x40C, 0x40D, 0x40E, 0x411,
	SmallKing, BigKing}

func Val(val int16) int {
	return int(val & 0xf)
}

func Colr(val int16) int {
	return int(val >> 8)
}

func CardStr(val int16) string {
	k := int(val >> 8)
	c := int(val & 0xf)
	return CARD_COLORS[k] + CARD_VALS[c]
}

func ToStr(vals []int16) string {
	var str string = ""
	for i := range vals {
		str += CardStr(vals[i]) + ","
	}
	return str
}
