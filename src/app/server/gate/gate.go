package gate

import (
	"ants/cluster"
	"ants/conf"
	"ants/gnet"
	"app/command"
	"fmt"
)

var logon *LogonMode
var router cluster.INetCluster

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
	gnet.NewTcpServer(port, func(tx gnet.INetContext) gnet.INetProxy {
		session := NewSession(tx)
		tx.SetHandle(func(event int, bits []byte) {
			pack := gnet.NewPackBytes(bits)
			switch pack.Topic() {
			case conf.TOPIC_GATE: //自身处理
				{
					mode.Done(pack.Cmd(), session, pack)
				}
			case conf.TOPIC_CLIENT: //直接推送给客户端
				{
					//获取客户端的头
					UserID, SessionID := pack.ReadInt(), pack.ReadUInt64()
					//获取用户
					if session, ok := logon.GetUserBySession(UserID, SessionID); ok {
						body := pack.ReadBytes(0)
						session.Send(gnet.NewPackArgs(pack.Cmd(), body))
					}
				}
			default: //通知模块
				{
					if session.IsLogin() {
						player := session.Player
						body := pack.GetBody()
						psend := gnet.NewPackArgs(pack.Cmd(), player.UserID, player.GateID, player.SessionID, body)
						router.Send(pack.Topic(), psend)
					} else {
						fmt.Println("连接用户尚未登录")
					}
				}
			}
		})
		return session
	}).Start()
}

func on_router_block(node cluster.INodeRouter, data interface{}) {
	pack := data.(gnet.ISocketPacket)
	if pack.Cmd() == gnet.EVENT_HEARTBEAT_PINT {
		node.Push(gnet.NewPackArgs(gnet.EVENT_HEARTBEAT_PINT))
	}
}
