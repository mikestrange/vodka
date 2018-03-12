package game

import "fmt"
import "fat/gnet"
import "app/command"

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_ENTER_TEXAS_ROOM: on_enter,
	command.CLIENT_LEAVE_TEXAS_ROOM: on_exit,
	command.CLIENT_TEXAS_SITDOWN:    on_sitdown,
	command.CLIENT_TEXAS_STAND:      on_stand,
	command.CLIENT_TEXAS_CHIP:       on_chip,
}

func on_enter(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	//进入()
	room_id := pack.ReadShort()
	//玩家基本信息(可以异步获取)
	body := pack.ReadBytes(0)
	//---
	fmt.Println("Enter Texas:", header, room_id, len(body))
}

func on_exit(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	//退出uid
	fmt.Println("Leave Texas:", header)
}

func on_sitdown(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	//操作
	seat_id, seat_chip, auto_buy := pack.ReadByte(), pack.ReadInt(), pack.ReadByte()
	fmt.Println("Sit Texas:", header, seat_id, seat_chip, auto_buy)
}

func on_stand(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	//站起来uid
	fmt.Println("Stand Texas:", header)
}

func on_chip(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	//操作类型+金额+uid(根据uid查找seatid)
	mtype, chip := pack.ReadByte(), pack.ReadInt()
	fmt.Println("Chip Texas:", header, mtype, chip)
}
