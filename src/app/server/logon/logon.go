package logon

import "ants/gnet"
import "ants/core"
import "ants/gcode"
import "app/conf"

//弱连接服务器，不用管心跳
func ServerLaunch(port int) {
	//数据服务器链接
	init_dber()
	//模块调度
	ref := new(LogicActor)
	if conf.LOCAL_TEST {
		core.Main().Join(conf.TOPIC_LOGON, ref)
	} else {
		core.RunAndThrowBox(ref, nil, func() {
			//重启
		})
		run_service(port, ref)
	}
}

func run_service(port int, ref core.IBox) {
	//服务器快速启动
	gnet.RunAndThrowServer(new(gnet.TCPServer), port, func(conn interface{}) gnet.IAgent {
		session := gnet.NewProxy(conn)
		session.SetReceiver(func(b []byte) {
			ref.Push(gnet.NewBytes(session, b))
		})
		return session
	}, func() {
		ref.Die()
	})
}

//逻辑块
type LogicActor struct {
	core.Box
}

func (this *LogicActor) OnReady() {
	this.SetName("登录服务器")
	this.SetAgent(this)
}

func (this *LogicActor) Handle(event interface{}) {
	//并发处理
	this.Wrap(func() {
		this.OnMessage(event)
	})
}

func (this *LogicActor) OnMessage(event interface{}) {
	pack := event.(*gnet.SocketEvent).BeginPack()
	if f, ok := events[pack.Cmd()]; ok {
		f.(func(gcode.ISocketPacket))(pack)
	}
}
