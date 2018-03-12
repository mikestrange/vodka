package world

import "app/command"
import "fat/gnet"
import "fmt"

var events map[int]interface{} = map[int]interface{}{
	command.SERVER_WORLD_ADD_PLAYER:     on_add_player,
	command.SERVER_WORLD_REMOVE_PLAYER:  on_remove_player,
	command.SERVER_WORLD_NOTICE_PLAYERS: on_notice_players,
	command.SERVER_WORLD_KICK_PLAYER:    on_kick_player,
	command.SERVER_WORLD_GET_ONLINE_NUM: on_online_player,
	command.SERVER_WORLD_NOTICE_TEST:    on_notice_test,
}

//登录发送来(登录通知)
func on_add_player(packet gnet.ISocketPacket) {
	code := packet.ReadShort()
	player := NewPlayer(packet)
	body := packet.ReadBytes(0)
	//失败
	if code != 0 {
		router.Send(player.SerID(), packet_logon_result(code, player.UserID, player.SessionID, body))
		return
	}
	//添加玩家
	if kick_player, ok := SetUser(player); ok {
		if kick_player.GateID != player.GateID {
			//如果是自己的重复登陆，网关控制下
			fmt.Println(fmt.Sprintf("Send kick Ok# uid=%d, session=%v", kick_player.UserID, kick_player.SessionID))
			router.Send(kick_player.SerID(), packet_kick_player(1, kick_player))
		} else {
			fmt.Println("Send kick Err# uid same gate login:", player.UserID, player.GateID)
		}
	}
	fmt.Println(fmt.Sprintf("Enter World Ok# uid=%d, session=%v, gate=%d", player.UserID, player.SessionID, player.GateID))
	router.Send(player.SerID(), packet_logon_result(code, player.UserID, player.SessionID, body))
}

//移除玩家(网关通知)
func on_remove_player(packet gnet.ISocketPacket) {
	//头部(网关通知)
	uid, gateid, session := packet.ReadInt(), packet.ReadInt(), packet.ReadUInt64()
	if player, ok := GetUser(uid); ok {
		if session == player.SessionID && gateid == player.GateID { //同一网关和同一会话id才行
			fmt.Println("Remove Ok# user=", uid)
			RemoveUser(uid)
		} else {
			fmt.Println(uid, "No match user# get:", gateid, session, " local:", player.SessionID, player.GateID)
		}
	} else {
		fmt.Println("Rmove Err# no user:", uid)
	}
}

//直接踢掉用户(任何地方)
func on_kick_player(pack gnet.ISocketPacket) {
	code, uid := pack.ReadShort(), pack.ReadInt()
	if player, ok := RemoveUser(uid); ok {
		router.Send(player.SerID(), packet_kick_player(code, player))
	} else {
		fmt.Println("Kick Err# no user:", uid)
	}
}

//通知世界所有角色(可以直接连世界)
func on_notice_players(pack gnet.ISocketPacket) {
	uid, cmd, body := pack.ReadInt(), int(pack.ReadInt()), pack.ReadBytes(0)
	fmt.Println("Notice World # user=", uid, ",cmd=", cmd, ",size=", len(body))
	NoticeAllUser(func(player *GamePlayer) interface{} {
		return packet_send_client(player.UserID, player.SessionID, cmd, uid, body)
	})
}

func on_online_player(tx gnet.INetContext, pack gnet.ISocketPacket) {
	onlines := int32(len(players))
	fmt.Println("当前在线人数:", onlines)
	tx.Send(gnet.NewPacketWithArgs(command.SERVER_WORLD_GET_ONLINE_NUM, onlines))
}

func on_notice_test(pack gnet.ISocketPacket) {

}
