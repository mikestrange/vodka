package texas

import "ants/gcode"
import "app/command"
import "app/conf"

//对外消息

func send_enter(uid int, name string) {
	gcode.NewPackTopic(command.CLIENT_GAME_ENTER, conf.TOPIC_CLIENT, uid, name)
}

//推送给用户房间的数据
func send_table_info() {

}

func send_leave(uid int) {
	gcode.NewPackTopic(command.CLIENT_GAME_LEAVE, conf.TOPIC_CLIENT, uid)
}

func send_sit(seatid int, uid int, name string, chip int, gift int) {
	gcode.NewPackTopic(command.CLIENT_GAME_SIT, conf.TOPIC_CLIENT, int8(seatid), uid, name, chip, gift)
}

func send_stand(seatid int) {
	gcode.NewPackTopic(command.CLIENT_GAME_STAND, conf.TOPIC_CLIENT, int8(seatid))
}

func send_start(seat []*Seat) {

}

func send_turn(pots []int) {

}

//true提前结束
func send_over(fast bool) {

}

func send_user_chip_start(seatid int, needcall int, needchip int, maxchip int) {
	gcode.NewPackTopic(command.SERVER_TEXAS_BET_START, conf.TOPIC_CLIENT, int8(seatid), needcall, needchip, maxchip)
}

func send_user_chip_result(uid int, action int, chip int) {
	gcode.NewPackTopic(command.SERVER_TEXAS_BET_RESULT, conf.TOPIC_CLIENT, uid, int8(action), chip)
}

func send_message(uid int, message string) {
	gcode.NewPackTopic(command.CLIENT_TEXAS_MESSAGE, conf.TOPIC_CLIENT, uid, message)
}

//每个池分配的用户
func send_result_pot(pot int, money int, seats []*Seat) {
	pack := gcode.NewPacket()
	pack.WriteBeginWithTopic(command.SERVER_TEXAS_RESULT, conf.TOPIC_CLIENT)
	pack.WriteByte(pot)
	pack.WriteInt(money)
	pack.WriteByte(len(seats))
	//	for i := range seats {

	//	}
	pack.WriteEnd()
}
