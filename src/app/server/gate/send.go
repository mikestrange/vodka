package gate

import "ants/gnet"
import "app/command"

//登录通知
func packet_logon_notice(args ...interface{}) gnet.IByteArray {
	return gnet.NewPackArgs(command.CLIENT_LOGON, args...)
}

//登录返回
func pack_logon_result(code int, args ...interface{}) gnet.IByteArray {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.CLIENT_LOGON)
	psend.WriteShort(code)
	psend.WriteValue(args...)
	psend.WriteEnd()
	return psend
}

//踢出
func packet_kick_user(code int) gnet.IByteArray {
	return gnet.NewPackArgs(command.SERVER_KICK_PLAYER, int16(code))
}

//世界添加用户
func packet_world_addplayer(args ...interface{}) gnet.IByteArray {
	return gnet.NewPackArgs(command.SERVER_WORLD_ADD_PLAYER, args...)
}

//世界移除用户
func packet_world_delplayer(args ...interface{}) gnet.IByteArray {
	return gnet.NewPackArgs(command.SERVER_WORLD_REMOVE_PLAYER, args...)
}

//通知世界转发给所有玩家(测试用)
func packet_world_notices(cmd int, uid int, bits []byte) gnet.IByteArray {
	return gnet.NewPackArgs(command.SERVER_WORLD_NOTICE_PLAYERS, uid, cmd, bits)
}
