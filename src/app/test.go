package app

import "app/command"

//
import "ants/gnet"
import "ants/base"
import "ants/gcode"
import "app/conf"
import "fmt"
import "strings"

//世界删除
func Test_remove_player(uid int) {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		var code int16 = 1
		tx.SetReceiver(func(b []byte) {
			pack := gcode.NewPackBytes(b)
			pack.ReadBegin()
			if pack.ReadShort() == 0 {
				fmt.Println("踢出用户成功:", pack.ReadInt())
			} else {
				fmt.Println("踢出用户失败:", pack.ReadInt())
			}
		})
		//gnet.RunAndThrowAgent(tx)
		tx.Send(gcode.NewPackArgs(command.SERVER_WORLD_KICK_PLAYER, code, uid))
	}
}

//通知世界派送消息
func Test_send_all() {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		var uid int32 = 100000
		cmd := command.CLIENT_NOTICE_CHANNEL
		var cid int32 = 10086
		var fromid int32 = uid
		var mtype int16 = 0
		message := base.Ltoa(base.Nano())
		message += "|"
		for i := 0; i < 10; i++ {
			message += "abcde"
		}
		//gnet.RunAndThrowAgent(tx)
		tx.Send(gcode.NewPackArgs(command.SERVER_WORLD_NOTICE_PLAYERS, uid, cmd, cid, fromid, mtype, message))
		tx.SetReceiver(func(b []byte) {
			println("on read")
		})
	}
}

//获取在线用户
func Test_get_online() {
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_WORLD).Addr); ok {
		tx.Send(gcode.NewPackArgs(command.SERVER_WORLD_GET_ONLINE_NUM))
		tx.SetReceiver(func(b []byte) {
			pack := gcode.NewPackBytes(b)
			pack.ReadBegin()
			fmt.Println("当前在线人数 socket:", pack.ReadInt())
		})
		//gnet.RunAndThrowAgent(tx)
	}
}

func Test_max_login(idx int) {
	i := idx
	for Test_login_send(i, "") {
		base.Sleep(5)
		i++
		if i > 5000 {
			return
		}
	}
}

func Test_login_send(idx int, pwd string) bool {
	uid := int32(idx)
	if tx, ok := gnet.Socket(conf.GetRouter(conf.PORT_GATE).Addr); ok {
		test_socket(uid, pwd, tx)
		return true
	}
	return false
}

func test_socket(uid int32, pwd string, tx gnet.Context) {
	tx.Send(gcode.NewPackArgs(command.CLIENT_LOGON, uid, pwd))
	//
	t := base.Nano()
	tx.SetReceiver(func(bits []byte) {
		packet := gcode.NewPackBytes(bits)
		packet.ReadBegin()
		switch packet.Cmd() {
		case gnet.EVENT_CONN_HEARTBEAT:
			{
				tx.Send(gcode.NewPackArgs(gnet.EVENT_CONN_HEARTBEAT))
			}
		case command.CLIENT_LOGON:
			{
				//packet.Print()
				code := packet.ReadShort()
				//body := packet.ReadBytes(0)
				//info
				//				if code == 0 {
				//					name := packet.ReadString()
				//					exp := packet.ReadInt()
				//					money := packet.ReadInt64()
				//					vipexp := packet.ReadInt()
				//					viptype := packet.ReadInt()
				//					pion := packet.ReadInt()
				//					fmt.Println("登录成功, 用户数据:", name, exp, money, vipexp, viptype, pion)
				//				}
				fmt.Println("客户端登录: err=", code, ",UID=", uid, ",runtime=", base.NanoStr(base.Nano()-t))
				//改名
				//psend := gcode.NewPackTopic(command.CLIENT_CHANGE_NAME, conf.TOPIC_HALL, "我不是谁，谁不是我")
				//tx.Send(psend)
				//
				//psend1 := gcode.NewPackTopic(command.CLIENT_JOIN_CHANNEL, conf.TOPIC_CHAT, int32(10086), "test1")
				//tx.Send(psend1)
				//str := base.Int64ToString(base.GetTimer())
				//psend2 := gcode.NewPackTopic(command.CLIENT_NOTICE_CHANNEL, conf.TOPIC_CHAT, int32(10086), int16(1), str)
				//tx.Send(psend2)
				//psend3 := gcode.NewPackTopic(command.CLIENT_GAME_ENTER, conf.TOPIC_GAME, 101)
				//tx.Send(psend3)
				//psend5 := gcode.NewPackTopic(command.CLIENT_GAME_SIT, conf.TOPIC_GAME, 101, int8(uid%6), 1024, true)
				//tx.Send(psend5)
				//psend4 := gcode.NewPackTopic(command.CLIENT_QUIT_CHANNEL, conf.TOPIC_CHAT, int32(10086))
				//tx.Send(psend4)
			}
		case command.CLIENT_NOTICE_CHANNEL:
			{
				cid, fromid, _, message := packet.ReadInt(), packet.ReadInt(), packet.ReadShort(), packet.ReadString()
				strs := strings.Split(message, "|")
				chat_time := base.Atol(strs[0])
				fmt.Println(uid, ">收到[", fromid, "]在频道[", cid, "] 消息延迟:", base.NanoStr(base.Nano()-chat_time), ", size=", len(strs[1]))
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
	//gnet.RunAndThrowAgent(tx)
}

func Test(b int) {
	fmt.Println("200勇士1秒入侵服务器:", b)
	base.Sleep(1000)
	for i := b * 200; i < b*200+200; i++ {
		Test_login_send(i, "abc123")
	}
}
