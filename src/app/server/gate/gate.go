package gate

import (
	"ants/core"
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
		session := NewSession(conn)
		if check_blacklist(session.Conn()) {
			glog.Warn("black list link: %s", session.Conn().Remote())
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
		request_topic_local(data.Tx(), pack)
	case conf.TOPIC_CLIENT:
		request_topic_client(pack)
	default:
		request_topic_actor(data.Tx(), pack)
	}
}

//黑名单
func check_blacklist(tx gnet.IConn) bool {
	return false //base.FindOk(tx.Remote(), "127.0.0.1")
}
