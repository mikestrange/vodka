package main

import "app/command"
import "app/config"
import "app/server"
import "fat/gsys"
import "fat/glog"
import "fat/gnet"

//import "fat/gnet/nsc"
import "fat/gutil"

const host string = "127.0.0.1:8081" //"120.77.149.74:8081" //

func test_send(idx int) {
	my_uid := int32(1 + idx)
	t := gutil.GetNano()
	if tx, ok := gnet.NewSocket(host); ok {
		go gnet.LoopWithHandle(tx, func(tx gnet.INetContext, data interface{}) {
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
					println("客户端登录: err=", code, ",UID=", my_uid, ",body=", body, ",runtime=", (gutil.GetNano()-t)/1000000, "毫秒")
					psend := gnet.NewPacketWithTopic(command.CLIENT_JOIN_CHANNEL, config.TOPIC_CHAT, int32(10086), "test1")
					tx.Send(psend)
					psend2 := gnet.NewPacketWithTopic(command.CLIENT_NOTICE_CHANNEL, config.TOPIC_CHAT, int32(10086), int16(1), "我是谁")
					tx.Send(psend2)
					psend3 := gnet.NewPacketWithTopic(command.CLIENT_ENTER_TEXAS_ROOM, config.TOPIC_GAME, int16(102))
					tx.Send(psend3)
					//psend4 := gnet.NewPacketWithTopic(command.CLIENT_QUIT_CHANNEL, config.TOPIC_CHAT, int32(10086))
					//tx.Send(psend4)
				}
			case command.CLIENT_NOTICE_CHANNEL:
				{
					cid, fromid, mtype, message := packet.ReadInt(), packet.ReadInt(), packet.ReadShort(), packet.ReadString()
					//ts := gutil.ParseInt(message, 0)
					println(my_uid, ">收到[", fromid, "]在频道[", cid, "]发送的消息:", message, mtype)
				}
			case command.CLIENT_MOVE:
				{
					uid, x, y, z := packet.ReadInt(), packet.ReadShort(), packet.ReadShort(), packet.ReadShort()
					println(uid, "移动:", x, y, z)
				}
			default:
				println("客户端未处理:", packet.Cmd())
			}
		})
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

func test() {
	go func() {
		gutil.Sleep(1000)
		//for {
		for i := 0; i < 100; i++ {
			test_send(i)
		}
		//test_send(1)
		//}

	}()
}

func main() {
	glog.LogAndRunning(glog.LOG_DEBUG, 100000)
	glog.Debug("运行路径=%s", gutil.Pwd())
	//启动服务器
	if gutil.MatchSys(1, "ser") {
		server.Launch()
	}
	if gutil.MatchSys(1, "cli") {
		go test()
	}
	//go http_echo()
	//
	gutil.Add(1)
	gutil.Wait()
}
