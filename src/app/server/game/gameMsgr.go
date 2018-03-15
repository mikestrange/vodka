package game

import "fmt"
import "ants/gnet"
import "app/command"

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_ENTER_TEXAS_ROOM: on_enter,
	command.CLIENT_LEAVE_TEXAS_ROOM: on_exit,
	command.CLIENT_TEXAS_SITDOWN:    on_sitdown,
	command.CLIENT_TEXAS_STAND:      on_stand,
	command.CLIENT_TEXAS_CHIP:       on_chip,
}

func on_enter(pack gnet.ISocketPacket) {
	//头部
	header := NewHeader(pack)
	//进入()
	room_id := int(pack.ReadShort())
	SendTable(room_id, pack, header)
	//---
	fmt.Println("Enter Texas:", header, room_id)
}

func on_exit(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	//如果能同时进入多个房间就需要roomid
	room_id := int(pack.ReadShort())
	SendTable(room_id, pack, header)
	//退出
	fmt.Println("Leave Texas:", header)
}

func on_sitdown(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	room_id := int(pack.ReadShort())
	//操作
	seat_id, seat_chip, auto_buy := pack.ReadByte(), pack.ReadInt(), pack.ReadByte()
	SendTable(room_id, pack, header, seat_id, seat_chip, auto_buy)
	fmt.Println("Sit Texas:", header, seat_id, seat_chip, auto_buy)
}

func on_stand(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	room_id := int(pack.ReadShort())
	SendTable(room_id, pack, header)
	//站起来uid
	fmt.Println("Stand Texas:", header)
}

func on_chip(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	room_id := int(pack.ReadShort())
	//操作类型+金额
	mtype, chip := pack.ReadByte(), pack.ReadInt()
	SendTable(room_id, pack, header, mtype, chip)
	fmt.Println("Chip Texas:", header, mtype, chip)
}
