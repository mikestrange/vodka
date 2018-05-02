package gate

import (
	"ants/base"
	"ants/core"
	"ants/gcode"
	"ants/glog"
	"ants/gnet"
	"app/command"
	"app/conf"
)

var events map[int]interface{} = map[int]interface{}{
	gnet.EVENT_CONN_CLOSE:       on_logout, //系统关闭
	command.CLIENT_LOGON:        on_logon,
	command.CLIENT_LOGOUT:       on_logout,
	command.SERVER_LOGON_RESULT: on_logon_result,
	command.SERVER_KICK_PLAYER:  on_kick,
}

var logon *LogonMode

func init() {
	logon = NewLogonMode()
}

//客户端请求登录
func on_logon(session *GateSession, packet gcode.ISocketPacket) {
	if session.IsBegin() {
		return
	}
	session.LoginWait()
	player := session.Player
	player.UserID = packet.ReadInt()
	player.PassWord = packet.ReadString()
	player.GateID = gate_idx //网关id
	player.SessionID = base.UsedSessionID()
	player.RegTime = base.Timer()
	//加入登陆队列
	logon.CommitLogon(session)
	glog.Debug("Logon Begin# uid=%d, session=%v gate=%d", player.UserID, player.SessionID, player.GateID)
	psend := packet_logon_notice(player.UserID, player.PassWord, player.GateID, player.SessionID)
	core.Main().Send(conf.TOPIC_LOGON, gnet.NewPack(session, psend))
}

//世界返回登录结果
func on_logon_result(packet gcode.ISocketPacket) {
	code, UserID, SessionID := packet.ReadShort(), packet.ReadInt(), packet.ReadUInt64()
	glog.Debug("Login Result# code=%d, uid=%d, session=%v", code, UserID, SessionID)
	if session, ok := logon.CompleteLogon(UserID, SessionID); ok {
		if code == 0 {
			session.LoginOk()
			if oplayer, ok2 := logon.AddUser(UserID, session); ok2 { //踢掉上一个用户
				kick_player(oplayer, 1)
			}
			//推送给自己
			user_info := packet.ReadBytes(0)
			session.Send(pack_logon_result(0, user_info))
		} else {
			session.CloseOf(pack_logon_result(code))
		}
	} else {
		if code == 0 {
			//失败后返回世界移除
			psend := packet_world_delplayer(UserID, gate_idx, SessionID)
			core.Main().Send(conf.TOPIC_WORLD, gnet.NewPack(nil, psend))
		}
	}
}

//客户端主动通知(关闭后自己也会通知一次)
func on_logout(session *GateSession, pack gcode.ISocketPacket) {
	if session.IsLive() {
		player := session.Player
		glog.Debug("close user=%d", player.UserID)
		session.LoginOut()
		//等待列表删除
		logon.CompleteLogon(player.UserID, player.SessionID)
		//登录成功后的删除
		logon.RemoveUserWithSession(player.UserID, player.SessionID)
		//通知世界或者游戏
		psend := packet_world_delplayer(player.UserID, player.GateID, player.SessionID)
		core.Main().Send(conf.TOPIC_WORLD, gnet.NewPack(session, psend))
	}
	session.Close()
}

//被踢(世界通知的命令)
func on_kick(pack gcode.ISocketPacket) {
	UserID, sessionID := pack.ReadInt(), pack.ReadUInt64()
	code := pack.ReadShort()
	if player, ok := logon.RemoveUser(UserID); ok {
		glog.Debug("gate for world kick ok: uid=%d session=%v", UserID, sessionID)
		kick_player(player, code)
	} else {
		glog.Debug("gate for world kick err: uid=%d code=%d", UserID, code)
	}
}

//commom
func kick_player(session *GateSession, code int) {
	//被踢的时候不会上报
	glog.Debug("Kick User ok# code=%d uid=%d, session=%v", code, session.Player.UserID, session.Player.SessionID)
	session.KickOut()
	session.CloseOf(packet_kick_user(code))
}
