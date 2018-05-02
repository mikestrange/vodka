package logon

import "app/command"
import "ants/gcode"
import "app/conf"
import "ants/core"
import "ants/glog"
import "ants/gnet"

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_LOGON: on_logon,
}

//message
func on_logon(pack gcode.ISocketPacket) {
	//header
	UserID, PassWord, SerID, SessionID := pack.ReadInt(), pack.ReadString(), pack.ReadInt(), pack.ReadUInt64()
	//other
	glog.Debug("Logon Info# uid=%d, session=%v gateid=%d", UserID, SessionID, SerID)
	err_code := check_user(UserID, PassWord)
	glog.Debug("Seach Result: code=%d uid=%d", err_code, UserID)
	var body []byte = []byte{}
	if err_code == 0 {
		body = get_user_info(UserID)
	}
	//错误直接返回
	if err_code != 0 {
		psend := pack_logon_result(err_code, UserID, SerID, SessionID, body)
		core.Main().Send(conf.TOPIC_WORLD, gnet.NewPack(nil, psend))
	} else {
		psend := pack_logon_result(err_code, UserID, SerID, SessionID, body)
		core.Main().Send(conf.TOPIC_WORLD, gnet.NewPack(nil, psend))
	}
}

//send world(通知登录结果)
func pack_logon_result(code int, uid int, gate int, session uint64, body []byte) interface{} {
	return gcode.NewPackArgs(command.SERVER_WORLD_ADD_PLAYER, int16(code), uid, gate, session, body)
}
