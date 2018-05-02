package hall

import "ants/gnet"
import "ants/core"
import "ants/gcode"
import "app/conf"

//弱连接服务器，不用管心跳
func ServerLaunch(port int) {
	//数据服务器链接
	init_dber()
	//模块调度
	ref := core.NewBox(new(LogicActor), "大厅服务")
	if conf.LOCAL_TEST {
		core.Main().Join(conf.TOPIC_HALL, ref)
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
}

func (this *LogicActor) Handle(event interface{}) {
	pack := event.(*gnet.SocketEvent).BeginPack()
	if f, ok := events[pack.Cmd()]; ok {
		f.(func(gcode.ISocketPacket))(pack)
	}
}
