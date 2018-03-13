package app

import "app/config"
import "app/command"
import "fat/gnet"
import "fat/gutil"
import "fat/gsys"
import "fmt"
import "strings"

//用户链接测试
const host string = "120.77.149.74:8081" //"127.0.0.1:8081" //

//世界删除
func Test_remove_player(uid int) {
	if tx, ok := gnet.NewSocket(config.GetDataRouter(config.WORLD_PORT).Addr()); ok {
		var code int16 = 1
		tx.Send(gnet.NewPacketWithArgs(command.SERVER_WORLD_KICK_PLAYER, code, int32(uid)))
		go func() {
			gutil.Sleep(100)
			tx.Close()
		}()
		tx.ReadBytes(1024, func(bits []byte) {
			pack := gnet.NewPacketWithBytes(bits)
			pack.ReadBegin()
			if pack.ReadShort() == 0 {
				fmt.Println("踢出用户成功:", pack.ReadInt())
			} else {
				fmt.Println("踢出用户失败:", pack.ReadInt())
			}
		})
	}
}

//通知世界派送消息
func Test_send_all() {
	if tx, ok := gnet.NewSocket(config.GetDataRouter(config.WORLD_PORT).Addr()); ok {
		var uid int32 = 100000
		var cmd int32 = int32(command.CLIENT_NOTICE_CHANNEL)
		var cid int32 = 10086
		var fromid int32 = uid
		var mtype int16 = 0
		message := gutil.Int64ToString(gutil.GetTimer())
		message += "|"
		for i := 0; i < 950; i++ {
			message += "A"
		}
		go func() {
			gutil.Sleep(100)
			tx.Close()
		}()
		tx.Send(gnet.NewPacketWithArgs(command.SERVER_WORLD_NOTICE_PLAYERS,
			uid, cmd, cid, fromid, mtype, message))
	}
}

//获取在线用户
func Test_get_online() {
	if tx, ok := gnet.NewSocket(config.GetDataRouter(config.WORLD_PORT).Addr()); ok {
		tx.Send(gnet.NewPacketWithArgs(command.SERVER_WORLD_GET_ONLINE_NUM))
		go func() {
			gutil.Sleep(100)
			tx.Close()
		}()
		tx.ReadBytes(1024, func(bits []byte) {
			pack := gnet.NewPacketWithBytes(bits)
			pack.ReadBegin()
			fmt.Println("当前在线人数:", pack.ReadInt())
		})
	}
}

func Test_login_send(idx int) {
	gutil.Sleep(50)
	uid := int32(idx)
	if tx, ok := gnet.NewSocket(host); ok {
		go test_socket(uid, tx)
		//登录
		psend := gnet.NewPacket()
		psend.WriteBegin(command.CLIENT_LOGON)
		psend.WriteValue(uid, "abc123")
		psend.WriteEnd()
		tx.Send(psend)
	}
}

func test_socket(uid int32, tx gnet.INetContext) {
	defer fmt.Println("客户端关闭: uid=", uid)
	t := gutil.GetNano()
	gnet.LoopWithHandle(tx, func(tx gnet.INetContext, data interface{}) {
		packet := data.(gnet.ISocketPacket)
		switch packet.Cmd() {
		case gnet.GNET_HEARTBEAT_PINT:
			{
				tx.Send(gnet.PacketWithHeartBeat)
			}
		case command.CLIENT_LOGON:
			{
				code := packet.ReadShort()
				body := packet.ReadBytes(0)
				tt := gutil.TimeStr((gutil.GetNano() - t) / 1000000)
				fmt.Println("客户端登录: err=", code, ",UID=", uid, ",body=", body, ",runtime=", tt)
				//psend := gnet.NewPacketWithTopic(command.CLIENT_JOIN_CHANNEL, config.TOPIC_CHAT, int32(10086), "test1")
				//tx.Send(psend)
				//str := gutil.Int64ToString(gutil.GetTimer())
				//psend2 := gnet.NewPacketWithTopic(command.CLIENT_NOTICE_CHANNEL, config.TOPIC_CHAT, int32(10086), int16(1), str)
				//tx.Send(psend2)
				//psend3 := gnet.NewPacketWithTopic(command.CLIENT_ENTER_TEXAS_ROOM, config.TOPIC_GAME, int16(1024))
				//tx.Send(psend3)
				//psend5 := gnet.NewPacketWithTopic(command.CLIENT_TEXAS_SITDOWN, config.TOPIC_GAME, int16(1024), int8(uid), int32(1024), int8(1))
				//tx.Send(psend5)
				//psend4 := gnet.NewPacketWithTopic(command.CLIENT_QUIT_CHANNEL, config.TOPIC_CHAT, int32(10086))
				//tx.Send(psend4)
			}
		case command.CLIENT_NOTICE_CHANNEL:
			{
				packet.ReadInt()
				cid, fromid, _, message := packet.ReadInt(), packet.ReadInt(), packet.ReadShort(), packet.ReadString()
				strs := strings.Split(message, "|")
				chat_time := gutil.ParseInt(strs[0], 0)
				fmt.Println(uid, ">收到[", fromid, "]在频道[", cid, "] 消息延迟:", gutil.TimeStr(gutil.GetTimer()-chat_time), ", size=", len(strs[1]))
			}
		default:
			fmt.Println("客户端未处理:", packet.Cmd())
		}
	})
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
