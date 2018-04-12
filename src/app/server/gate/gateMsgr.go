package gate

import (
	"ants/actor"
	"ants/conf"
	"ants/gnet"
	"ants/gsys"
	"ants/gutil"
	"app/command"
	"fmt"
)

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_LOGON:        on_logon,
	command.CLIENT_LOGOUT:       on_logout,
	command.SERVER_LOGON_RESULT: on_logon_result,
	command.SERVER_KICK_PLAYER:  on_kick,
}

//客户端请求登录
func on_logon(session *GateSession, packet gnet.ISocketPacket) {
	if session.IsBegin() {
		return
	}
	session.LoginWait()
	//添加到登陆列表里面
	player := session.Player
	//
	player.UserID = packet.ReadInt()
	player.PassWord = packet.ReadString()
	player.GateID = int32(gate_idx) //小心识别
	player.SessionID = gsys.MainSession.UsedSessionID()
	player.RegTime = gutil.GetTimer() //可以设计超时
	//加入登陆队列
	logon.CommitLogon(session)
	fmt.Println(fmt.Sprintf("Logon Begin# uid=%d, session=%v serid=%d", player.UserID, player.SessionID, player.GateID))
	actor.Main.Send(conf.TOPIC_LOGON, packet_logon_notice(player.UserID, player.PassWord, player.GateID, player.SessionID))
}

//世界返回登录结果
func on_logon_result(packet gnet.ISocketPacket) {
	code, UserID, SessionID := packet.ReadShort(), packet.ReadInt(), packet.ReadUInt64()
	fmt.Println(fmt.Sprintf("Login Result# code=%d, uid=%d, session=%v", code, UserID, SessionID))
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
			actor.Main.Send(conf.TOPIC_WORLD, packet_world_delplayer(UserID, gate_idx, SessionID))
		}
	}
}

//客户端主动通知(关闭后自己也会通知一次)
func on_logout(session *GateSession) {
	if session.IsLive() {
		player := session.Player
		session.LoginOut()
		//等待列表删除
		logon.CompleteLogon(player.UserID, player.SessionID)
		//登录成功后的删除
		logon.RemoveUserWithSession(player.UserID, player.SessionID)
		//通知世界或者游戏
		actor.Main.Send(conf.TOPIC_WORLD, packet_world_delplayer(player.UserID, player.GateID, player.SessionID))
	}
	session.Close()
}

//被踢(世界通知的命令)
func on_kick(packet gnet.ISocketPacket) {
	UserID, _ := packet.ReadInt(), packet.ReadUInt64()
	code := packet.ReadShort()
	if player, ok := logon.RemoveUser(UserID); ok {
		fmt.Println("World guest kick UID Ok:", UserID, code)
		kick_player(player, code)
	} else {
		fmt.Println("World guest kick UID Err:", UserID, code)
	}
}

//commom
func kick_player(session *GateSession, code int16) {
	//被踢的时候不会上报
	fmt.Println(fmt.Sprintf("Kick User ok# code=%d uid=%d, session=%v", code, session.Player.UserID, session.Player.SessionID))
	session.KickOut()
	session.CloseOf(packet_kick_user(code))
}
