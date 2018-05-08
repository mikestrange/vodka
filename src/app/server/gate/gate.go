package gate

import (
	"ants/core"
	"ants/gcode"
	"ants/glog"
	"ants/gnet"
	"app/conf"
)

var gate_idx int
var ref core.IBox

//服务器的启动
func ServerLaunch(port int, gid int) {
	gate_idx = gid
	ref = new(LogicActor)
	if conf.LOCAL_TEST {
		core.Main().Join(conf.TOPIC_GATE, ref)
	} else {
		core.RunAndThrowBox(ref, nil, func() {
			//通知世界，然后重启
		})
	}
	run_gate_service(port, ref)
}

func run_gate_service(port int, ref core.IBox) {
	gnet.RunAndThrowServer(new(gnet.TCPServer), port, func(conn interface{}) gnet.IAgent {
		session := NewSession()
		session.SetProcesser(gcode.NewClient())
		if check_blacklist(session) {
			glog.Warn("black list link: %s", session.Remote())
			session.Close()
		} else {
			session.SetReceiver(func(b []byte) {
				ref.Push(gnet.NewBytes(session, b))
			})
		}
		return session
	}, func() {
		ref.Die() //服务器关闭，关闭当前的进程
	})
}

//逻辑块(单线)
type LogicActor struct {
	core.BaseBox
}

func (this *LogicActor) OnReady() {
	this.SetName("网关服务器")
	this.SetAgent(this)
	this.SetBlock(this.OnMessage)
}

//不同的消息放入不同的线程
func (this *LogicActor) OnMessage(event interface{}) {
	data := event.(*gnet.SocketEvent)
	pack := data.BeginPack()
	glog.Debug("gate cmd: %d", pack.Cmd())
	switch pack.Topic() {
	case conf.TOPIC_SELF:
		this.local_handle(data.Tx(), pack)
	case conf.TOPIC_CLIENT:
		this.send_to_client(pack)
	default:
		this.send_actor_handle(data.Tx().(*GateSession), pack)
	}
}

//本地处理
func (this *LogicActor) local_handle(session gnet.Context, pack gcode.ISocketPacket) {
	if block, ok := events[pack.Cmd()]; ok {
		switch f := block.(type) {
		case func(*GateSession, gcode.ISocketPacket):
			f(session.(*GateSession), pack)
		case func(gcode.ISocketPacket):
			f(pack)
		default:
			glog.Debug("not handle func:", pack.Topic(), pack.Cmd())
		}
	}
}

//发送给模块
func (this *LogicActor) send_actor_handle(session *GateSession, pack gcode.ISocketPacket) {
	if session.IsLogin() {
		player := session.Player
		body := pack.GetBody()
		psend := gcode.NewPackArgs(pack.Cmd(), player.UserID, player.GateID, player.SessionID, body)
		core.Main().Send(pack.Topic(), gnet.NewPack(session, psend))
	}
}

//发送给客户端
func (this *LogicActor) send_to_client(pack gcode.ISocketPacket) {
	UserID, SessionID, body := pack.ReadInt(), pack.ReadUInt64(), pack.ReadRemaining()
	if target, ok := logon.GetUserBySession(UserID, SessionID); ok {
		target.Send(gcode.NewPackArgs(pack.Cmd(), body))
	}
}

//黑名单
func check_blacklist(tx gnet.IConn) bool {
	return false //base.FindOk(tx.Remote(), "127.0.0.1")
}
