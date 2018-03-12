package gate

import (
	"app/command"
	"app/config"
	"fat/gnet"
	"fat/gutil"
	"fmt"
)

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_LOGON:        on_logon,
	command.CLIENT_LOGOUT:       on_logout,
	command.SERVER_LOGON_RESULT: on_logon_result,
	command.SERVER_KICK_PLAYER:  on_kick,
	//RPG
	command.CLIENT_MOVE:   on_move,
	command.CLIENT_ACTION: on_action,
	command.CLIENT_RESULT: on_result,
}

//客户端请求登录
func on_logon(tx gnet.INetContext, packet gnet.ISocketPacket) {
	if tx.Client() != nil {
		fmt.Println("Context is Logining")
		return
	}
	//添加到登陆列表里面
	player := &UserData{
		UserID:    packet.ReadInt(),
		PassWord:  packet.ReadString(),
		GateID:    int32(router.Data().RouteID()),
		SessionID: router.UsedSessionID(),
		RegTime:   gutil.GetTimer()}
	//绑定客户端(超时关闭)
	tx.SetClient(player)
	//加入登陆队列
	logon.CommitLogon(&GamePlayer{tx, player})
	fmt.Println(fmt.Sprintf("Logon Begin# uid=%d, session=%v serid=%d", player.UserID, player.SessionID, player.GateID))
	router.Send(config.TOPIC_LOGON, packet_logon_notice(player.UserID, player.PassWord, player.GateID, player.SessionID))
}

//世界返回登录结果
func on_logon_result(tx gnet.INetContext, packet gnet.ISocketPacket) {
	code, UserID, SessionID := packet.ReadShort(), packet.ReadInt(), packet.ReadUInt64()
	fmt.Println(fmt.Sprintf("Login Result# code=%d, uid=%d, session=%v", code, UserID, SessionID))
	if player, ok := logon.CompleteLogon(UserID, SessionID); ok {
		if code == 0 {
			player.Player.Status = LOGON_OK
			if oplayer, ok2 := logon.AddUser(UserID, player); ok2 { //踢掉上一个用户
				kick_player(oplayer, 1)
			}
			//推送给自己
			body := packet.ReadBytes(0) //自己的一些信息
			player.Conn.Send(pack_logon_result(0, body))
		} else {
			player.Conn.SetClient(nil)
			player.Conn.Send(pack_logon_result(code))
			player.Conn.Close()
		}
	}
}

//客户端主动通知
func on_logout(tx gnet.INetContext, packet gnet.ISocketPacket) {
	tx.Close() //退出的时候会通知世界
}

//被踢(一般世界发送)
func on_kick(packet gnet.ISocketPacket) {
	code, UserID := packet.ReadShort(), packet.ReadInt()
	if player, ok := logon.RemoveUser(UserID); ok {
		kick_player(player, code)
	}
}

//commom
func kick_player(player *GamePlayer, code int16) {
	//被踢的时候不会上报
	fmt.Println(fmt.Sprintf("Kick User ok# code=%d uid=%d, session=%v", code, player.Player.UserID, player.Player.SessionID))
	player.Conn.SetClient(nil)
	player.Conn.Send(packet_kick_user(code))
	player.Conn.Close()
}
