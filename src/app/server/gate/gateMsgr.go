package gate

import (
	"ants/conf"
	"ants/gnet"
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
	player.GateID = int32(router.Data().RouteID())
	player.SessionID = router.UsedSessionID()
	player.RegTime = gutil.GetTimer()
	//加入登陆队列
	logon.CommitLogon(session)
	fmt.Println(fmt.Sprintf("Logon Begin# uid=%d, session=%v serid=%d", player.UserID, player.SessionID, player.GateID))
	router.Send(conf.TOPIC_LOGON, packet_logon_notice(player.UserID, player.PassWord, player.GateID, player.SessionID))
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
			body := packet.ReadBytes(0) //自己的一些信息
			session.Send(pack_logon_result(0, body))
		} else {
			session.Send(pack_logon_result(code))
			session.Close()
		}
	}
}

//客户端主动通知
func on_logout(session *GateSession) {
	//直接被关闭了
	session.Kill()
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
	session.Send(packet_kick_user(code))
	session.Close()
}
