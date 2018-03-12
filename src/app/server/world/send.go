package world

import "app/command"
import "fat/gnet"
import "app/config"

//踢人
func packet_kick_player(code int16, player *GamePlayer) gnet.IBytes {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.SERVER_KICK_PLAYER)
	psend.WriteValue(player.UserID, player.SessionID)
	psend.WriteShort(code)
	psend.WriteEnd()
	return psend
}

//登录返回
func packet_logon_result(code int16, uid int32, session uint64, body []byte) gnet.IBytes {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.SERVER_LOGON_RESULT)
	psend.WriteValue(code, uid, session, body)
	//成功写入资料
	psend.WriteEnd()
	return psend
}

//直接发送给客户端
func packet_send_client(uid int32, session uint64, cmd int, fromid int32, body []byte) gnet.IBytes {
	psend := gnet.NewPacket()
	psend.WriteBeginWithTopic(cmd, config.TOPIC_CLIENT)
	psend.WriteValue(uid, session) //网关校正
	psend.WriteValue(fromid)
	psend.WriteBytes(body)
	psend.WriteEnd()
	return psend
}
