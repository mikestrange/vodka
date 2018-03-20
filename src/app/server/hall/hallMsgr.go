package hall

import (
	"ants/gnet"
	"app/command"
	"fmt"
)

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_CHANGE_NAME: on_change_name,
	command.CLIENT_CHANGE_INFO: on_change_info,
}

func on_change_name(pack gnet.ISocketPacket) {
	UserID, GateID, SessionID := pack.ReadInt(), pack.ReadInt(), pack.ReadUInt64()
	//
	name := pack.ReadString()
	//改名
	fmt.Println("用户改名:", UserID, GateID, SessionID, name)
	if ok := change_name(int(UserID), name); ok {
		router.Send(int(GateID), pack_change_name(0, UserID, SessionID, name))
	} else {
		router.Send(int(GateID), pack_change_name(1, UserID, SessionID, ""))
	}
}

func on_change_info(pack gnet.ISocketPacket) {

}
