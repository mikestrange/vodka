package chat

import (
	"ants/gnet"
	"app/command"
	"fmt"
)

//列表管理
var tables *TableManager = NewManager()

var events map[int]interface{} = map[int]interface{}{
	command.SERVER_BUILD_CHANNEL:  on_build_channel,
	command.SERVER_REMOVE_CHANNEL: on_remove_channel,
	command.CLIENT_JOIN_CHANNEL:   on_join_channel,
	command.CLIENT_QUIT_CHANNEL:   on_quit_channel,
	command.CLIENT_NOTICE_CHANNEL: on_message,
}

func on_build_channel(packet gnet.ISocketPacket) {
	cid, ctype := packet.ReadInt(), packet.ReadByte()
	tables.CreateTable(cid, ctype)
}

func on_remove_channel(packet gnet.ISocketPacket) {
	cid := packet.ReadInt()
	if _, ok := tables.RemoveTable(cid); ok {
		//通知频道所有的人

	}
}

func on_join_channel(packet gnet.ISocketPacket) {
	header := NewHeader(packet)
	//频道ID
	cid := packet.ReadInt()
	if table, ok := tables.GetTable(cid); ok {
		player := &GameUser{Player: header, UserName: packet.ReadString()}
		if table.AddUser(player) {
			fmt.Println("Enter Chat ok: uid=", player.Player.UserID, ",cid=", cid, ", gate=", player.SerID())
			table.NoticeAllUser(func(data *GameUser) interface{} {
				return pack_join_table(data.Player.UserID, data.Player.SessionID, cid, player.Player.UserID)
			})
		} else {
			fmt.Println("Enter Chat Update# user = ", header.UserID, ", cid=", cid)
		}
	} else {
		fmt.Println("Enter Chat Err# no table ", cid)
	}
}

func on_quit_channel(packet gnet.ISocketPacket) {
	header := NewHeader(packet)
	//频道ID
	cid := packet.ReadInt()
	if table, ok := tables.GetTable(cid); ok {
		if player, ok2 := table.RemoveUser(header.UserID); ok2 {
			fmt.Println("Exit Chat Ok# user=", header.UserID, ", cid=", cid)
			//通知所有
			table.NoticeAllUser(func(data *GameUser) interface{} {
				return pack_exit_table(data.Player.UserID, data.Player.SessionID, cid, player.Player.UserID)
			})
			//通知自己
			psend := pack_exit_table(player.Player.UserID, player.Player.SessionID, cid, player.Player.UserID)
			router.Send(player.SerID(), psend)
		} else {
			fmt.Println("Exit Chat Err# no user [", header.UserID, "]in cid=", cid)
		}
	} else {
		fmt.Println("Exit Chat Err# no table ", cid)
	}
}

func on_message(packet gnet.ISocketPacket) {
	header := NewHeader(packet)
	//频道ID
	cid := packet.ReadInt()
	if table, ok := tables.GetTable(cid); ok {
		mtype, message := packet.ReadShort(), packet.ReadString()
		if player, ok2 := table.GetUser(header.UserID); ok2 {
			//通知所有
			fmt.Println("Notice Chat Ok: user=", header.UserID, ", cid=", cid, ", size=", len(message))
			table.NoticeAllUser(func(data *GameUser) interface{} {
				return pack_message(data.Player.UserID, data.Player.SessionID, cid, player.Player.UserID, mtype, message)
			})
		} else {
			fmt.Println("Notice Chat Err: no user [", header.UserID, "] in cid=", cid, ", size=", len(message))
		}
	} else {
		fmt.Println("Notice Chat Err# no table ", cid)
	}
}
