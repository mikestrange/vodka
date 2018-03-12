package gate

import "fat/gnet"
import "app/command"

//登录通知
func packet_logon_notice(args ...interface{}) gnet.IByteArray {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.CLIENT_LOGON)
	psend.WriteValue(args...)
	psend.WriteEnd()
	return psend
}

//登录返回
func pack_logon_result(code int16, args ...interface{}) gnet.IByteArray {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.CLIENT_LOGON)
	psend.WriteValue(code)
	psend.WriteValue(args...)
	psend.WriteEnd()
	return psend
}

//踢出
func packet_kick_user(code int16) gnet.IByteArray {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.SERVER_KICK_PLAYER)
	psend.WriteValue(code)
	psend.WriteEnd()
	return psend
}

//世界添加用户
func packet_world_addplayer(args ...interface{}) gnet.IByteArray {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.SERVER_WORLD_ADD_PLAYER)
	psend.WriteValue(args...)
	psend.WriteEnd()
	return psend
}

//世界移除用户
func packet_world_delplayer(args ...interface{}) gnet.IByteArray {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.SERVER_WORLD_REMOVE_PLAYER)
	psend.WriteValue(args...)
	psend.WriteEnd()
	return psend
}

//通知世界转发给所有玩家(测试用)
func packet_world_notices(cmd int, uid int32, bits []byte) gnet.IByteArray {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.SERVER_WORLD_NOTICE_PLAYERS)
	psend.WriteInt(uid) //发送对象
	psend.WriteInt(int32(cmd))
	psend.WriteBytes(bits)
	psend.WriteEnd()
	return psend
}
