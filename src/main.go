package main

import "app/command"
import "app/config"
import "app/server"
import "fat/gsys"
import "fat/glog"
import "fat/gnet"
import "fmt"

//import "fat/gnet/nsc"
import "fat/gutil"

const host string = "120.77.149.74:8081" //"127.0.0.1:8081" //

func test_send(idx int) {
	my_uid := int32(1 + idx)

	if tx, ok := gnet.NewSocket(host); ok {
		go test_socket(my_uid, tx)
		//登录
		go gsys.After(10, func() {
			psend := gnet.NewPacket()
			psend.WriteBegin(command.CLIENT_LOGON)
			psend.WriteValue(my_uid, "abc123")
			psend.WriteEnd()
			tx.Send(psend)
		})
	}
}

func test_socket(my_uid int32, tx gnet.INetContext) {
	defer fmt.Println("客户端关闭: uid=", my_uid)
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
				fmt.Println("客户端登录: err=", code, ",UID=", my_uid, ",body=", body, ",runtime=", (gutil.GetNano()-t)/1000000, "毫秒")
				//psend := gnet.NewPacketWithTopic(command.CLIENT_JOIN_CHANNEL, config.TOPIC_CHAT, int32(10086), "test1")
				//tx.Send(psend)
				//str := gutil.Int64ToString(gutil.GetTimer())
				//psend2 := gnet.NewPacketWithTopic(command.CLIENT_NOTICE_CHANNEL, config.TOPIC_CHAT, int32(10086), int16(1), str)
				//tx.Send(psend2)
				psend3 := gnet.NewPacketWithTopic(command.CLIENT_ENTER_TEXAS_ROOM, config.TOPIC_GAME, int16(102))
				tx.Send(psend3)
				//psend4 := gnet.NewPacketWithTopic(command.CLIENT_QUIT_CHANNEL, config.TOPIC_CHAT, int32(10086))
				//tx.Send(psend4)
			}
		case command.CLIENT_NOTICE_CHANNEL:
			{
				packet.ReadInt()
				cid, fromid, _, message := packet.ReadInt(), packet.ReadInt(), packet.ReadShort(), packet.ReadString()
				chat_time := gutil.ParseInt(message, 0)
				fmt.Println(my_uid, ">收到[", fromid, "]在频道[", cid, "] 消息延迟:", gutil.GetTimer()-chat_time, "毫秒")
			}
		default:
			fmt.Println("客户端未处理:", packet.Cmd())
		}
	})
}

func test(b int) {
	go func() {
		gutil.Sleep(1000)
		for i := b * 200; i < b*200+200; i++ {
			test_send(i)
			gutil.Sleep(50)
		}
	}()
}

func test_send_all() {
	if tx, ok := gnet.NewSocket(config.GetDataRouter(config.WORLD_PORT).Addr()); ok {
		var uid int32 = 100000
		var cmd int32 = int32(command.CLIENT_NOTICE_CHANNEL)
		var cid int32 = 10086
		var fromid int32 = uid
		var mtype int16 = 0
		message := gutil.Int64ToString(gutil.GetTimer())
		tx.Send(gnet.NewPacketWithArgs(command.SERVER_WORLD_NOTICE_PLAYERS,
			uid, cmd, cid, fromid, mtype, message))
		tx.Close()
	}
}

func main() {
	glog.LogAndRunning(glog.LOG_DEBUG, 100000)
	glog.Debug("运行路径=%s", gutil.Pwd())
	//启动服务器
	if gutil.MatchSys(1, "cli") {
		idx := gutil.Atoi(gutil.GetArgs(2))
		//size := gutil.Atoi(gutil.GetArgs(3))
		go test(idx)
	} else if gutil.MatchSys(1, "online") {
		if tx, ok := gnet.NewSocket(config.GetDataRouter(config.WORLD_PORT).Addr()); ok {
			tx.Send(gnet.NewPacketWithArgs(command.SERVER_WORLD_GET_ONLINE_NUM))
			tx.Close()
		}
	} else if gutil.MatchSys(1, "all") {
		test_send_all()
	} else {
		server.Launch()
		//gutil.Sleep(10)
		//test_send(10000)
		//gutil.Sleep(100)
		//test_send_all()
	}
	//go http_echo()
	//
	gutil.Add(1)
	gutil.Wait()
}
