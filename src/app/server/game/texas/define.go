package texas

//德州牌型尺寸
const (
	POCKET_CARD_SIZE = 5
)

//输赢结果
const (
	WIN  = 0
	LOSE = 1
	DRAW = 2
)

//计时器标志
const (
	TIME_ACTION = 1
	TIME_START  = 2
)

//结束标志
const (
	OVER_NONE  = 0 //正常结束
	OVER_EARLY = 1 //提前结束
)

//操作状态
const (
	SEAT_CHIP_CHECK = 1 //过牌
	SEAT_CHIP_FOLD  = 2 //弃牌
	SEAT_CHIP_CALL  = 3 //跟注
	SEAT_CHIP_RAISE = 4 //加注
)

//游戏状态
const (
	GAME_STATE_WAIT  = 0 //准备开始(人少的时候)
	GAME_STATE_START = 1 //开始(发手牌)
	GAME_START_FLOP  = 2 //翻牌
	GAME_START_TURN  = 3 //转牌
	GAME_START_RIVER = 4 //河牌
	GAME_START_STOP  = 5 //游戏结束
)

//牌型
const (
	TYPE_NONE_CARD      = 0  // 无型
	TYPE_PIE_CARD       = 1  // 杂牌
	TYPE_HIGH_CARD      = 2  // 高牌
	TYPE_PAIR           = 3  // 一对
	TYPE_TWO_PAIRS      = 4  // 二对
	TYPE_THREE_KIND     = 5  // 三条
	TYPE_STRAIGHT       = 6  // 顺子
	TYPE_FLUSH          = 7  // 同花
	TYPE_FULL_HOUSE     = 8  // 葫芦
	TYPE_FOUR_KIND      = 9  // 四条
	TYPE_STRAIGHT_FLUSH = 10 // 同花顺
	TYPE_ROYAL_FLUSH    = 11 // 皇家同花顺
)

var POKER_TYPE_STR []string = []string{"-", "杂牌", "高牌", "一对", "二对", "三条", "顺子", "同花", "葫芦", "四条", "同花顺", "皇家同花顺"}
var POKER_COLORS []string = []string{"-", "方片", "梅花", "红桃", "黑桃"}
var POKER_VALS []string = []string{"-", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

const POKER_SIZE = 52

var GLOBAL_POKER [52]int16 = [52]int16{
	0x102, 0x103, 0x104, 0x105, 0x106, 0x107, 0x108, 0x109, 0x10A, 0x10B, 0x10C, 0x10D, 0x10E,
	0x202, 0x203, 0x204, 0x205, 0x206, 0x207, 0x208, 0x209, 0x20A, 0x20B, 0x20C, 0x20D, 0x20E,
	0x302, 0x303, 0x304, 0x305, 0x306, 0x307, 0x308, 0x309, 0x30A, 0x30B, 0x30C, 0x30D, 0x30E,
	0x402, 0x403, 0x404, 0x405, 0x406, 0x407, 0x408, 0x409, 0x40A, 0x40B, 0x40C, 0x40D, 0x40E}

func TypeStr(val int) string {
	return POKER_TYPE_STR[val]
}

func CardStr(val int16) string {
	k := int(val >> 8)
	c := int(val & 0xf)
	return POKER_COLORS[k] + POKER_VALS[c]
}

func ToStr(vals []int16) string {
	var str string = ""
	for i := range vals {
		str += CardStr(vals[i]) + ","
	}
	return str
}

func CardVal(val int16) int {
	return int(val & 0xf)
}

func CardColor(val int16) int {
	return int(val >> 8)
}
