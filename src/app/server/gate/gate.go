package gate

import (
	"app/command"
	"app/config"
	"fat/gnet"
	"fat/gnet/nsc"
	"fat/gutil"
	"fmt"
)

var logon *LogonMode
var router nsc.IRemoteScheduler
var mode gutil.IModeAccessor

func init() {
	logon = NewLogonMode()
	//基础设施
	router, mode = command.SetRouter(config.GATE_PORT, events, nil)
}

//服务器的启动
func ServerLaunch(port int) {
	gnet.ListenAndRunServer(port, func(conn interface{}) {
		context := gnet.NewConn(conn)
		defer context_close_handle(context)
		gnet.LoopConnWithPing(context, func(tx gnet.INetContext, data interface{}) {
			pack := data.(gnet.ISocketPacket)
			//println("gate>", pack.Cmd(), pack.Topic())
			switch pack.Topic() {
			case config.TOPIC_GATE:
				on_context_handle(tx, pack)
			case config.TOPIC_CLIENT:
				on_client_handle(pack)
			default:
				on_topic_handle(tx, pack)
			}
		})
	})
}

//环境自身的处理
func on_context_handle(tx gnet.INetContext, pack gnet.ISocketPacket) {
	//如果登录走登录协议，不是登录校正登录后再处理(目前都可以)
	mode.Done(pack.Cmd(), tx, pack)
}

//通知其他模块(注入uid,serid,sessionid)
func on_topic_handle(tx gnet.INetContext, pack gnet.ISocketPacket) {
	if data, ok := check_logon(tx); ok {
		body := pack.GetBody()
		psend := gnet.NewPacketWithArgs(pack.Cmd(), data.UserID, data.GateID, data.SessionID, body)
		router.Send(pack.Topic(), psend)
	} else {
		fmt.Println("连接用户尚未登录")
	}
}

//派送给用户(这里属于异步,要校正会话ID)
func on_client_handle(pack gnet.ISocketPacket) {
	//获取客户端的头
	UserID, SessionID := pack.ReadInt(), pack.ReadUInt64()
	//获取用户
	if player, ok := logon.GetUserBySession(UserID, SessionID); ok {
		body := pack.ReadBytes(0)
		player.Conn.Send(gnet.NewPacketWithArgs(pack.Cmd(), body))
	}
}

//是否登陆
func check_logon(tx gnet.INetContext) (*UserData, bool) {
	client := tx.Client()
	switch client.(type) {
	case *UserData:
		return client.(*UserData), client.(*UserData).Status == LOGON_OK
	default:
	}
	return nil, false
}

//关闭时候处理
func context_close_handle(tx gnet.INetContext) {
	if data, ok := check_logon(tx); ok {
		tx.SetClient(nil)
		//登录等待列表
		logon.CompleteLogon(data.UserID, data.SessionID)
		//登录成功后的删除
		logon.RemoveUserWithSession(data.UserID, data.SessionID)
		//通知世界
		router.Send(config.TOPIC_WORLD, packet_world_delplayer(data.UserID, data.GateID, data.SessionID))
	} else {
		println("不是登录的客户端>")
	}
}
