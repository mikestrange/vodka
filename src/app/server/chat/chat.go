package chat

import "ants/gnet"
import "ants/core"
import "ants/gcode"
import "app/conf"

//服务器的启动(快速启动)
func ServerLaunch(port int) {
	ref := core.NewBox(new(LogicActor), "聊天服务")
	if conf.LOCAL_TEST {
		core.Main().Join(conf.TOPIC_CHAT, ref)
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

//处理模块
type LogicActor struct {
	//
}

func (this *LogicActor) Handle(event interface{}) {
	pack := event.(*gnet.SocketEvent).BeginPack()
	if block, ok := events[pack.Cmd()]; ok {
		switch f := block.(type) {
		case func(gcode.ISocketPacket):
			f(pack)
		}
	}
}
