package gate

import (
	"ants/actor"
	"ants/conf"
	"ants/gnet"
	"ants/gutil"
	"app/command"
)

var gate_idx int
var logon *LogonMode
var refLogic actor.IActorRef

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
func ServerLaunch(port int, gid int) {
	refLogic = actor.NewRefRunning(new(LogicActor))
	gate_idx = gid
	gnet.NewTcpServer(port, func(conn interface{}) gnet.INetProxy {
		session := NewSession(conn)
		//handle func begin
		session.SetHandle(func(b []byte) {
			refLogic.Router(session, gnet.NewPackBytes(b))
		})
		//handle func end
		return session
	}).Start()
}

//逻辑块(单线)
type LogicActor struct {
	actor.BaseActor
	mode gutil.IModeAccessor
}

func (this *LogicActor) OnReady(ref actor.IActorRef) {
	this.mode = command.SetMode(on_gate_handle, events)
	//--
	ref.Open()
}

func (this *LogicActor) OnClose() {

}

//不同的消息放入不同的线程
func (this *LogicActor) OnMessage(args ...interface{}) {
	session := args[0].(*GateSession)
	pack := args[1].(gnet.ISocketPacket)
	switch pack.Topic() {
	case conf.TOPIC_SELF:
		//自身处理
		this.mode.Done(pack.Cmd(), args...)
	case conf.TOPIC_CLIENT:
		//直接推送给客户端
		UserID, SessionID, body := pack.ReadInt(), pack.ReadUInt64(), pack.ReadBytes(0)
		if target, ok := logon.GetUserBySession(UserID, SessionID); ok {
			target.Send(gnet.NewPackArgs(pack.Cmd(), body))
		}
	default:
		//通知其他模块
		if session.IsLogin() {
			player := session.Player
			body := pack.GetBody()
			psend := gnet.NewPackArgs(pack.Cmd(), player.UserID, player.GateID, player.SessionID, body)
			actor.Main.Send(pack.Topic(), psend)
		}
	}
}
