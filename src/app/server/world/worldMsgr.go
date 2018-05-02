package world

import "app/command"
import "ants/gnet"
import "ants/gcode"
import "app/conf"
import "ants/core"
import "ants/glog"

var events map[int]interface{} = map[int]interface{}{
	command.SERVER_WORLD_ADD_PLAYER:     on_add_player,
	command.SERVER_WORLD_REMOVE_PLAYER:  on_remove_player,
	command.SERVER_WORLD_NOTICE_PLAYERS: on_notice_players,
	command.SERVER_WORLD_KICK_PLAYER:    on_kick_player,
	command.SERVER_WORLD_GET_ONLINE_NUM: on_online_player,
	command.SERVER_WORLD_NOTICE_TEST:    on_notice_test,
}

//登录发送来(登录通知)
func on_add_player(_ gnet.IAgent, packet gcode.ISocketPacket) {
	code := packet.ReadShort()
	player := NewPlayer(packet)
	body := packet.ReadBytes(0)
	//失败
	if code != 0 {
		core.Main().Send(player.GateID, packet_logon_result(code, player.UserID, player.SessionID, body))
		return
	}
	//添加玩家
	if kick_player, ok := SetUser(player); ok {
		if kick_player.GateID != player.GateID {
			//如果是自己的重复登陆，网关控制下
			glog.Debug("Send kick Ok# uid=%d, session=%v", kick_player.UserID, kick_player.SessionID)
			core.Main().Send(kick_player.GateID, packet_kick_player(1, kick_player))
		} else {
			glog.Debug("Send kick Err# uid=%d same gate=%d login:", player.UserID, player.GateID)
		}
	}
	glog.Debug("Enter World Ok# uid=%d, session=%v, gate=%d", player.UserID, player.SessionID, player.GateID)
	//返回给网关
	psend1 := packet_logon_result(code, player.UserID, player.SessionID, body)
	core.Main().Send(player.GateID, gnet.NewPack(nil, psend1))
	//通知游戏
	psend2 := gcode.NewPackArgs(command.SERVER_ADD_PLAYER, player.UserID, player.GateID, player.SessionID)
	core.Main().Send(conf.TOPIC_GAME, gnet.NewPack(nil, psend2))
}

//移除玩家(网关通知)
func on_remove_player(_ gnet.IAgent, pack gcode.ISocketPacket) {
	//头部(网关通知)
	uid, serid, session := pack.ReadInt(), pack.ReadInt(), pack.ReadUInt64()
	if player, ok := GetUser(uid); ok {
		if session == player.SessionID && serid == player.GateID { //同一网关和同一会话id才行
			glog.Debug("Remove Ok: user=%d", uid)
			RemoveUser(uid)
			//通知游戏
			psend := gcode.NewPackArgs(command.SERVER_DEL_PLAYER, player.UserID)
			core.Main().Send(conf.TOPIC_GAME, gnet.NewPack(nil, psend))
		} else {
			glog.Debug("No match uid=%d gate=%d session=%v and gate=%d session=%v",
				uid, serid, session, player.SessionID, player.GateID)
		}
	} else {
		glog.Debug("Rmove Err# no uid=%d", uid)
	}
}

//直接踢掉用户(任何地方)
func on_kick_player(session gnet.IAgent, pack gcode.ISocketPacket) {
	code, uid := pack.ReadShort(), pack.ReadInt()
	if player, ok := RemoveUser(uid); ok {
		core.Main().Send(player.GateID, gnet.NewPack(nil, packet_kick_player(code, player)))
		//session.CloseOf(gcode.NewPackArgs(pack.Cmd(), int16(0), uid))
	} else {
		glog.Debug("Kick Err# no uid=%d", uid)
		//session.CloseOf(gcode.NewPackArgs(pack.Cmd(), int16(1), uid))
	}
}

//通知世界所有角色(可以直接连世界)
func on_notice_players(session gnet.IAgent, pack gcode.ISocketPacket) {
	uid, cmd, body := pack.ReadInt(), int(pack.ReadInt()), pack.ReadBytes(0)
	glog.Debug("Notice World # uid=%d cmd=%d size=%d", uid, cmd, len(body))
	NoticeAllUser(func(player *GamePlayer) interface{} {
		return packet_send_client(cmd, player.UserID, player.SessionID, body)
	})
	//session.Close()
}

func on_online_player(session gnet.IAgent, pack gcode.ISocketPacket) {
	onlines := len(players)
	glog.Debug("online player num=%d", onlines)
	//session.CloseOf(gcode.NewPackArgs(command.SERVER_WORLD_GET_ONLINE_NUM, onlines))
}

func on_notice_test(session gnet.IAgent, pack gcode.ISocketPacket) {

}
