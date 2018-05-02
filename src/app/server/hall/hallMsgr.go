package hall

import "app/command"
import "ants/glog"
import "ants/core"
import "ants/gcode"
import "ants/gnet"

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_CHANGE_NAME: on_change_name,
	command.CLIENT_CHANGE_INFO: on_change_info,
}

func on_change_name(pack gcode.ISocketPacket) {
	uid, GateID, SessionID := pack.ReadInt(), pack.ReadInt(), pack.ReadUInt64()
	//---
	name := pack.ReadString()
	//改名
	glog.Debug("uid=%d change name=%d", uid, name)
	if ok := change_name(uid, name); ok {
		core.Main().Send(GateID, gnet.NewPack(nil, pack_change_name(0, uid, SessionID, name)))
	} else {
		core.Main().Send(GateID, gnet.NewPack(nil, pack_change_name(1, uid, SessionID, "")))
	}
}

func on_change_info(pack gcode.ISocketPacket) {

}
