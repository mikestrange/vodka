package world

import "ants/gnet"
import "ants/core"
import "ants/gcode"
import "app/conf"

func ServerLaunch(port int) {
	ref := new(LogicActor)
	if conf.LOCAL_TEST {
		core.Main().Join(conf.TOPIC_WORLD, ref)
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

type LogicActor struct {
	core.BaseBox
}

func (this *LogicActor) OnReady() {
	this.SetName("世界服务器")
	this.SetAgent(this)
	this.SetBlock(this.OnMessage)
}

func (this *LogicActor) OnMessage(event interface{}) {
	data := event.(*gnet.SocketEvent)
	pack := data.BeginPack()
	if block, ok := events[pack.Cmd()]; ok {
		switch f := block.(type) {
		case func(gnet.IAgent, gcode.ISocketPacket):
			f(data.Tx(), pack)
		}
	}
}
