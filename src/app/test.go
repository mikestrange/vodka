package app

import "app/command"

//
import "ants/gnet"
import "ants/conf"
import "ants/gutil"
import "ants/gsys"
import "fmt"
import "strings"

//世界删除
func Test_remove_player(uid int) {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		var code int16 = 1
		tx.Send(gnet.NewPackArgs(command.SERVER_WORLD_KICK_PLAYER, code, int32(uid)))
		tx.SetHandle(func(code int, bits []byte) {
			pack := gnet.NewPackBytes(bits)
			if pack.ReadShort() == 0 {
				fmt.Println("踢出用户成功:", pack.ReadInt())
			} else {
				fmt.Println("踢出用户失败:", pack.ReadInt())
			}
		})
		tx.Run()
	}
}

//通知世界派送消息
func Test_send_all() {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		var uid int32 = 100000
		var cmd int32 = int32(command.CLIENT_NOTICE_CHANNEL)
		var cid int32 = 10086
		var fromid int32 = uid
		var mtype int16 = 0
		message := gutil.Int64ToString(gutil.GetNano())
		message += "|"
		for i := 0; i < 100; i++ {
			message += "A"
		}
		tx.Send(gnet.NewPackArgs(command.SERVER_WORLD_NOTICE_PLAYERS, uid, cmd, cid, fromid, mtype, message))
		tx.Close()
	}
}

//获取在线用户
func Test_get_online() {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		tx.SetHandle(func(code int, bits []byte) {
			pack := gnet.NewPackBytes(bits)
			fmt.Println("当前在线人数:", pack.ReadInt())
		})
		tx.Send(gnet.NewPackArgs(command.SERVER_WORLD_GET_ONLINE_NUM))
		tx.Run()
		println("close online find")
	}
}

func Test_max_login(idx int) {
	i := idx
	for Test_login_send(i) {
		gutil.Sleep(5)
		i++
		if i > 5000 {
			return
		}
	}
}

func Test_login_send(idx int) bool {
	//gutil.Sleep(10)
	uid := int32(idx)
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_GATE).Addr); ok {
		go test_socket(uid, tx)
		return true
	}
	return false
}

func test_socket(uid int32, tx gnet.INetContext) {
	defer fmt.Println("客户端关闭: uid=", uid)
	//
	tx.Send(gnet.NewPackArgs(command.CLIENT_LOGON, uid, "abc123"))
	//
	t := gutil.GetNano()
	tx.SetHandle(func(code int, bits []byte) {
		packet := gnet.NewPackBytes(bits)
		switch packet.Cmd() {
		//		case gnet.GNET_HEARTBEAT_PINT:
		//			{
		//				tx.Send(gnet.PacketWithHeartBeat)
		//			}
		case command.CLIENT_LOGON:
			{
				packet.Print()
				code := packet.ReadShort()
				body := packet.ReadBytes(0)
				fmt.Println("客户端登录: err=", code, ",UID=", uid, ",body=", body, ",runtime=", gutil.NanoStr(gutil.GetNano()-t))
				psend := gnet.NewPackTopic(command.CLIENT_JOIN_CHANNEL, conf.TOPIC_CHAT, int32(10086), "test1")
				tx.Send(psend)
				//str := gutil.Int64ToString(gutil.GetTimer())
				//psend2 := gnet.NewPackTopic(command.CLIENT_NOTICE_CHANNEL, conf.TOPIC_CHAT, int32(10086), int16(1), str)
				//tx.Send(psend2)
				//psend3 := gnet.NewPackTopic(command.CLIENT_ENTER_TEXAS_ROOM, conf.TOPIC_GAME, int16(1024))
				//tx.Send(psend3)
				//psend5 := gnet.NewPackTopic(command.CLIENT_TEXAS_SITDOWN, conf.TOPIC_GAME, int16(1024), int8(uid), int32(1024), int8(1))
				//tx.Send(psend5)
				//psend4 := gnet.NewPackTopic(command.CLIENT_QUIT_CHANNEL, conf.TOPIC_CHAT, int32(10086))
				//tx.Send(psend4)
			}
		case command.CLIENT_NOTICE_CHANNEL:
			{
				packet.ReadInt()
				cid, fromid, _, message := packet.ReadInt(), packet.ReadInt(), packet.ReadShort(), packet.ReadString()
				strs := strings.Split(message, "|")
				chat_time := gutil.ParseInt(strs[0], 0)
				fmt.Println(uid, ">收到[", fromid, "]在频道[", cid, "] 消息延迟:", gutil.NanoStr(gutil.GetNano()-chat_time), ", size=", len(strs[1]))
			}
		case command.CLIENT_JOIN_CHANNEL:
			{
				cid, fuid := packet.ReadInt(), packet.ReadInt()
				fmt.Println(fuid, "进入房间:", cid)
			}
		default:
			fmt.Println("客户端未处理:", packet.Cmd())
		}
	})
	//
	tx.Run()
}

func Test(b int) {
	sub_time := 3
	fmt.Println("200勇士3秒入侵服务器:", 3)
	gsys.SetTimeout(1000, 3, func() {
		sub_time--
		fmt.Println("200勇士3秒入侵服务器:", sub_time)
	})
	go func() {
		gutil.Sleep(3000)
		for i := b * 200; i < b*200+200; i++ {
			go Test_login_send(i)
		}
	}()
}
