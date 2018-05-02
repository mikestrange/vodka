package landlord

import "ants/base"
import "fmt"

//大到小排序
const not_val = -1

//降序
func sortDrop(cards []int16) {
	base.Sort(cards, func(i, j int) bool {
		return Val(cards[i]) > Val(cards[j])
	})
}

//升序
func sortUper(cards []int16) {
	base.Sort(cards, func(i, j int) bool {
		return Val(cards[i]) < Val(cards[j])
	})
}

//是否连顺(必须连续性)
func check_even(cards []int16, num int) bool {
	if len(cards)%num != 0 {
		return false
	}
	val := not_val
	for i := 0; i < len(cards)/num; i++ {
		if val == not_val {
			val = Val(cards[i*num])
		} else {
			if val-i != Val(cards[i*num]) {
				return false
			}
		}
	}
	return true
}

//值是否一样
func check_some_val(cards []int16) bool {
	idx := not_val
	for i := range cards {
		if idx == not_val {
			idx = Val(cards[i])
		} else {
			if idx != Val(cards[i]) {
				return false
			}
		}
	}
	return true
}

//连续的值是否一样
func check_some_val_num(cards []int16, num int) bool {
	if len(cards)%num != 0 {
		return false
	}
	for i := 0; i < len(cards)/num; i++ {
		b := i * num     // 0
		e := i*num + num // 2-3-4
		if !check_some_val(cards[b:e]) {
			return false
		}
	}
	return true
}

//获取相同值的牌
func get_some_card(cards []int16, num int) int {
	val := not_val
	size := 0
	for i := range cards {
		if val == Val(cards[i]) {
			size++
		} else {
			val = Val(cards[i])
			size = 1
		}
		//
		if size == num {
			return i + 1
		}
	}
	return not_val
}

//取最大相同数目的牌列表
func get_some_cards(cards []int16, num int) ([]int16, []int16, int) {
	idx := 0
	var types []int16 //成型
	var subs []int16  //散牌
	current := cards[0:]
	for {
		current = current[idx:]
		new_idx := get_some_card(current, num)
		if new_idx == not_val {
			subs = append(subs, current...)
			break
		} else {
			types = append(types, current[(new_idx-num):new_idx]...)
			subs = append(subs, current[:(new_idx-num)]...)
			idx = new_idx
		}
	}
	length := len(types) / num
	//types = append(types, subs...)
	return types, subs, length
}

//抽取最大的牌组(3个，4个等)
func take_out_same(cards []int16, num int) ([]int16, []int16, bool) {
	idx := get_some_card(cards, num)
	if idx == not_val {
		return nil, nil, false
	}
	types := append([]int16{}, cards[(idx-num):idx]...)
	subs := append([]int16{}, cards[:(idx-num)]...)
	return types, subs, true
}

/*
各种牌型
*/
func check_type_one(cards []int16) bool { //单牌
	return len(cards) == 1
}

func check_type_pair(cards []int16) bool { //对子
	return len(cards) == 2 && check_some_val(cards)
}

func check_type_three(cards []int16) bool { //三不带
	return len(cards) == 3 && check_some_val(cards)
}

func check_type_three_one(cards []int16) bool { //三带单(可能是炸弹)
	if len(cards) != 4 {
		return false
	}
	if _, subs, ok := take_out_same(cards, 3); ok {
		return check_type_one(subs)
	}
	return false
}

func check_type_three_pair(cards []int16) bool { //三带对
	if len(cards) != 5 {
		return false
	}
	if _, subs, ok := take_out_same(cards, 3); ok {
		return check_type_pair(subs)
	}
	return false
}

//s
func check_type_one_s(cards []int16) bool { //顺子
	if len(cards) < 5 {
		return false
	}
	return check_even(cards, 1)
}

func check_type_pair_s(cards []int16) bool { //连对，必须3对及以上
	if len(cards) < 6 {
		return false
	}
	return check_even(cards, 2)
}

func check_type_three_s(cards []int16) bool { //飞机不带
	if len(cards) != 6 {
		return false
	}
	return check_even(cards, 3)
}

func check_type_three_one_s(cards []int16) bool { //飞机带单,有可能4个三对／三对带单(这里考虑为带牌,由客户端自行提交)
	types, _, count := get_some_cards(cards, 3)
	if count > 1 {
		//连顺
		if check_even(types[:count*3], 3) {
			return len(types[count*3:]) == count
		}
	}
	return false
}

func check_type_three_pair_s(cards []int16) bool { //飞机带对 (难点2)
	types, _, count := get_some_cards(cards, 3)
	if count > 1 {
		//连顺
		if check_even(types[:count*3], 3) {
			subs := types[count*3:]
			if len(subs) == count*2 {
				return check_some_val_num(subs, 2)
			}
		}
	}
	return false
}

//four
func check_type_four_one(cards []int16) bool { //4带两张(可以是一对)
	if _, subs, ok := take_out_same(cards, 4); ok {
		return len(subs) == 2
	}
	return false
}

func check_type_four_pair(cards []int16) bool { //4带两对(注释：可以为两炸)
	if _, subs, ok := take_out_same(cards, 4); ok {
		return len(subs) == 4 && check_some_val_num(subs, 2)
	}
	return true
}

//bomb
func check_type_bomb(cards []int16) bool { //炸弹
	if len(cards) != 4 {
		return false
	}
	return check_some_val(cards)
}

func check_type_rockets(cards []int16) bool { //王炸
	if len(cards) != 2 {
		return false
	}
	return cards[0] == SmallKing && cards[1] == BigKing
}

//获取牌型(优先级别)
func flush_card_type(cards []int16) int {
	sortDrop(cards)
	if check_type_rockets(cards) {
		return TYPE_ROCKETS
	} else if check_type_bomb(cards) {
		return TYPE_BOMB_HARD
	} else if check_type_four_pair(cards) { //可能两炸
		return TYPE_FOUR_TO_PAIR_S
	} else if check_type_four_one(cards) {
		return TYPE_FOUR_TO_ONE_S
	} else if check_type_three_pair_s(cards) {
		return TYPE_THREE_TO_PAIR_S
	} else if check_type_three_one_s(cards) { //可能飞机全部
		return TYPE_THREE_TO_ONE_S
	} else if check_type_three_s(cards) { //可能飞机带单,比如:333,444,555,666(特殊情况)
		return TYPE_THREE_S
	} else if check_type_pair_s(cards) {
		return TYPE_PAIR_S
	} else if check_type_one_s(cards) {
		return TYPE_ONE_S
	} else if check_type_three_pair(cards) {
		return TYPE_THREE_TO_PAIR
	} else if check_type_three_one(cards) {
		return TYPE_THREE_TO_ONE
	} else if check_type_three(cards) {
		return TYPE_THREE
	} else if check_type_pair(cards) {
		return TYPE_PAIR
	} else if check_type_one(cards) {
		return TYPE_ONE
	}
	return TYPE_NIL
}

//判断类型
func check_card_type(ctype int, cards []int16) bool {
	sortDrop(cards)
	switch ctype {
	case TYPE_ROCKETS:
		return check_type_rockets(cards)
	case TYPE_BOMB_HARD:
		return check_type_bomb(cards)
	case TYPE_FOUR_TO_PAIR_S:
		return check_type_four_pair(cards)
	case TYPE_FOUR_TO_ONE_S:
		return check_type_four_one(cards)
	case TYPE_THREE_TO_PAIR_S:
		return check_type_three_pair_s(cards)
	case TYPE_THREE_TO_ONE_S:
		return check_type_three_one_s(cards)
	case TYPE_THREE_S:
		return check_type_three_s(cards)
	case TYPE_PAIR_S:
		return check_type_pair_s(cards)
	case TYPE_ONE_S:
		return check_type_one_s(cards)
	case TYPE_THREE_TO_PAIR:
		return check_type_three_pair(cards)
	case TYPE_THREE_TO_ONE:
		return check_type_three_one(cards)
	case TYPE_THREE:
		return check_type_three(cards)
	case TYPE_PAIR:
		return check_type_pair(cards)
	case TYPE_ONE:
		return check_type_one(cards)
	}
	return false
}

func _init() {
	ll := []int16{0x0103, 0x0103, 0x0104, 0x0104, 0x0104, 0x0105, 0x0105, 0x0105, 0x0105, 0x0105}
	sortDrop(ll)
	fmt.Println(ToStr(ll))
	if check_type_three_pair_s(ll) {
		fmt.Println("飞机双带")
	} else {
		fmt.Println("不是飞机")
	}
}
