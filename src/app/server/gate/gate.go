package gate

import (
	"ants/cluster"
	"ants/conf"
	"ants/glog"
	"ants/gnet"
	"app/command"
)

var logon *LogonMode

func init() {
	logon = NewLogonMode()
}

func on_gate_handle(block interface{}, args ...interface{}) {
	switch block.(type) {
	case func(gnet.ISocketPacket):
		block.(func(gnet.ISocketPacket))(args[1].(gnet.ISocketPacket))
	case func(*GateSession, gnet.ISocketPacket):
		block.(func(*GateSession, gnet.ISocketPacket))(args[0].(*GateSession), args[1].(gnet.ISocketPacket))
	case func(*GateSession):
		block.(func(*GateSession))(args[0].(*GateSession))
	case func(gnet.IBaseProxy):
		block.(func(gnet.IBaseProxy))(args[0].(gnet.IBaseProxy))
	default:
		println("Err: no mode handle ")
	}
}

//服务器的启动
func ServerLaunch(port int) {
	mode := command.SetMode(on_gate_handle, events, false)
	//
	router = command.SetRouter(port, on_router_block)
	//
	gnet.NewTcpServer(port, func(conn interface{}) gnet.INetProxy {
		session := NewSession(conn)
		//handle func begin
		session.Tx().SetHandle(func(b []byte) {
			pack := gnet.NewPackBytes(b)
			switch pack.Topic() {
			case 0: //自身处理
				mode.Done(pack.Cmd(), session, pack)
			case conf.TOPIC_CLIENT: //直接推送给客户端
				//获取客户端的头
				UserID, SessionID, body := pack.ReadInt(), pack.ReadUInt64(), pack.ReadBytes(0)
				//获取用户
				if target, ok := logon.GetUserBySession(UserID, SessionID); ok {
					target.Send(gnet.NewPackArgs(pack.Cmd(), body))
				}
			default: //通知模块
				if session.IsLogin() {
					player := session.Player
					body := pack.GetBody()
					psend := gnet.NewPackArgs(pack.Cmd(), player.UserID, player.GateID, player.SessionID, body)
					router.Send(pack.Topic(), psend)
				} else {
					glog.Debug("连接用户尚未登录")
				}
			}
		})
		//handle func end
		return session
	}).Start()
}

var router cluster.INodeRouter

func on_router_block(client interface{}, data interface{}) {
	node := client.(cluster.INodeRouter)
	pack := data.(gnet.ISocketPacket)
	if pack.Cmd() == gnet.EVENT_HEARTBEAT_PINT {
		node.Push(gnet.NewPackArgs(gnet.EVENT_HEARTBEAT_PINT))
	}
}
