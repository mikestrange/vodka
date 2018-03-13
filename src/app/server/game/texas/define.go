package texas

//计时器状态
const (
	TIME_CHIP_STATE = 1
	TIME_OVER_STATE = 2
)

//德州牌型尺寸
const (
	POCKET_CARD_SIZE = 5
)

//操作状态
const (
	SEAT_CHIP_CHECK = 1 //过牌
	SEAT_CHIP_FOLD  = 2 //弃牌
	SEAT_CHIP_CALL  = 3 //跟注
	SEAT_CHIP_RAISE = 4 //加注
)

//玩家状态
const (
	PLAYER_READY_STATE = 0 //等待游戏
	PLAYER_PLAYING     = 1 //游戏中
	PLAYER_FOLD_CARD   = 2 //弃牌
	PLAYER_ALLIN       = 3 //全下
)

//游戏状态
const (
	GAME_STATE_READY = 1 //准备开始(人少的时候)
	GAME_STATE_START = 2 //开始(发手牌)
	GAME_START_FLOP  = 3 //翻牌
	GAME_START_TURN  = 4 //转牌
	GAME_START_RIVER = 5 //河牌
	GAME_START_STOP  = 6 //游戏结束
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

var POKER_TYPE_STR []string = []string{"nil", "杂牌", "高牌", "一对", "二对", "三条", "顺子", "同花", "葫芦", "四条", "同花顺", "皇家同花顺"}
var POKER_COLORS []string = []string{"nil", "方片", "梅花", "红桃", "黑桃"}
var POKER_VALS []string = []string{"nil", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

const POKER_SIZE = 52

var GLOBAL_POKER [52]int16 = [52]int16{
	0x102, 0x103, 0x104, 0x105, 0x106, 0x107, 0x108, 0x109, 0x10A, 0x10B, 0x10C, 0x10D, 0x10E,
	0x202, 0x203, 0x204, 0x205, 0x206, 0x207, 0x208, 0x209, 0x20A, 0x20B, 0x20C, 0x20D, 0x20E,
	0x302, 0x303, 0x304, 0x305, 0x306, 0x307, 0x308, 0x309, 0x30A, 0x30B, 0x30C, 0x30D, 0x30E,
	0x402, 0x403, 0x404, 0x405, 0x406, 0x407, 0x408, 0x409, 0x40A, 0x40B, 0x40C, 0x40D, 0x40E}

func PokerString(val int16) string {
	k := int(val >> 8)
	c := int(val & 0xf)
	return POKER_COLORS[k] + POKER_VALS[c]
}

func PokerStrings(vals []int16) string {
	var str string = ""
	for _, val := range vals {
		str += PokerString(val) + ","
	}
	return str
}
