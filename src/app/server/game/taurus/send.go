package taurus

import "ants/gcode"
import "app/command"
import "app/conf"

//对外消息
type GameSend struct {
	BaseGame
}

//发送给用户
func (this *GameSend) send_player(user *Player, pack gcode.ISocketPacket) {

}

func (this *GameSend) send_reconnect(player *Player) {

}

//推送给用户房间的数据
func (this *GameSend) send_table_info(player *Player) {
	pack := gcode.NewPacket()
	pack.WriteBeginWithTopic(command.CLIENT_GAME_TABLE_INFO, conf.TOPIC_CLIENT)
	pack.WriteInt(this.table_id)
	pack.WriteByte(this.table_type)
	pack.WriteInt(this.base_chip)
	pack.WriteByte(this.seat_count)
	pack.WriteByte(this.banker_time)
	pack.WriteByte(this.chip_time)
	pack.WriteByte(this.commit_time)
	pack.WriteByte(this.game_state)
	pack.WriteByte(this.banker_idx)
	pack.WriteByte(this.banker_multiple)
	//
	pack.WriteByte(this.get_sit_num())
	this.each_sits(func(seat *Seat) {
		pack.WriteInt(seat.uid)
		pack.WriteByte(seat.seat_id)
		pack.WriteInt(seat.seat_money)
		pack.WriteString(seat.name)
		pack.WriteInt(seat.gift)
		pack.WriteByte(seat.bet_multiple) //下倍
		pack.WriteByte(seat.game_status)
		//自己的话追加牌
		if seat.uid == player.uid && seat.isplayer() {
			pack.WriteByte(len(seat.cards))
			for p := range seat.cards {
				pack.WriteShort(int(seat.cards[p]))
			}
		}
	})
	pack.WriteEnd()
}

func (this *GameSend) send_enter_err(player *Player, code int) {
	gcode.NewPackTopic(command.CLIENT_GAME_ENTER, conf.TOPIC_CLIENT, int16(code))
}

func (this *GameSend) broadcast_enter(uid int, name string) {
	this.each_players(func(player *Player) {
		p := gcode.NewPackTopic(command.SERVER_BROADCAST_GAME_ENTER, conf.TOPIC_CLIENT, player.uid, player.session,
			uid, name)
		//core.Main().Send(player.gate, p)
		println(p)
	})
}

func (this *GameSend) send_leave_err(player *Player, code int) {
	gcode.NewPackTopic(command.SERVER_BROADCAST_GAME_LEAVE, conf.TOPIC_CLIENT, int16(code))
}

func (this *GameSend) broadcast_leave(uid int) {
	gcode.NewPackTopic(command.CLIENT_GAME_LEAVE, conf.TOPIC_CLIENT, uid)
}

func (this *GameSend) send_sit_err(player *Player, code int) {
	gcode.NewPackTopic(command.CLIENT_GAME_SIT, conf.TOPIC_CLIENT, int16(code))
}

func (this *GameSend) broadcast_sit(seatid int, uid int, money int, name string, gift int) {
	gcode.NewPackTopic(command.SERVER_BROADCAST_GAME_SIT, conf.TOPIC_CLIENT, int8(seatid), uid, money, name, gift)
}

func (this *GameSend) send_stand_err(uid int, code int) {
	gcode.NewPackTopic(command.CLIENT_GAME_STAND, conf.TOPIC_CLIENT, int16(code))
}

func (this *GameSend) broadcast_stand(seatid int) {
	gcode.NewPackTopic(command.SERVER_BROADCAST_GAME_STAND, conf.TOPIC_CLIENT, int8(seatid))
}

//游戏开始
func (this *GameSend) broadcast_start() {
	pack := gcode.NewPacket()
	pack.WriteBeginWithTopic(command.SERVER_NIUNIU_START, conf.TOPIC_CLIENT)
	pack.WriteByte(this.get_action_num())
	this.each_actions(func(seat *Seat) {
		pack.WriteByte(seat.seat_id)
		pack.WriteInt(seat.uid)
		pack.WriteInt(seat.seat_money)
	})
	pack.WriteEnd()
}

//给自己发牌
func (this *GameSend) send_self_cards(seat *Seat) {
	pack := gcode.NewPacket()
	pack.WriteBeginWithTopic(command.SERVER_NIUNIU_USER_CARDS, conf.TOPIC_CLIENT)
	pack.WriteInt(seat.uid)
	list := seat.cards
	pack.WriteByte(len(list))
	for i := range list {
		pack.WriteShort(int(list[i]))
	}
	pack.WriteEnd()
}

//广播谁抢庄
func (this *GameSend) broadcast_rob_banker(seatid int, num int) {
	gcode.NewPackTopic(command.CLIENT_NIUNIU_BANKER, conf.TOPIC_CLIENT, int8(seatid), int8(num))
}

//开始下注(倍数)
func (this *GameSend) broadcast_bet_start() {
	gcode.NewPackTopic(command.SERVER_NIUNIU_BET_START, conf.TOPIC_CLIENT)
}

//下倍
func (this *GameSend) broadcast_user_bet(seatid int, num int) {
	gcode.NewPackTopic(command.CLIENT_NIUNIU_BET, conf.TOPIC_CLIENT, int8(seatid), int8(num))
}

//开始提交
func (this *GameSend) broadcast_commit_start() {
	gcode.NewPackTopic(command.SERVER_NIUNIU_COMMIT_START, conf.TOPIC_CLIENT)
}

//提交
func (this *GameSend) broadcast_user_commit(seatid int) {
	gcode.NewPackTopic(command.CLIENT_NIUNIU_COMMIT, conf.TOPIC_CLIENT, int8(seatid))
}

//结束
func (this *GameSend) broadcast_over() {
	gcode.NewPackTopic(command.SERVER_NIUNIU_OVER, conf.TOPIC_CLIENT)
}

//每个池分配的用户
func (this *GameSend) broadcast_result(list []*Seat) {
	pack := gcode.NewPacket()
	pack.WriteBeginWithTopic(command.SERVER_NIUNIU_RESULT, conf.TOPIC_CLIENT)
	pack.WriteByte(len(list))
	for i := range list {
		seat := list[i]
		pack.WriteInt(seat.uid)
		pack.WriteByte(seat.seat_id)
		pack.WriteInt(seat.seat_money)
		pack.WriteByte(seat.result)
		pack.WriteInt(seat.result_money)
		pack.WriteByte(seat.card_type)
		pack.WriteByte(len(seat.cards))
		for p := range seat.cards {
			pack.WriteShort(int(seat.cards[p]))
		}
	}
	pack.WriteEnd()
}
