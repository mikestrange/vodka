package hall

//非网关模块通用(处理一些无关痛痒的事务)
import "ants/cluster"
import "ants/gnet"
import "app/command"

//服务器的启动(快速启动)
func ServerLaunch(port int) {
	init_dber()
	//模块调度
	mode := command.SetMode(nil, events, false)
	//分布式路由
	router = command.SetRouter(port, on_router_block)
	//服务器快速启动
	gnet.ListenAndRunServer(port, func(session gnet.IBaseProxy) {
		session.SetHandle(func(bits []byte) {
			pack := gnet.NewPackBytes(bits)
			mode.Done(pack.Cmd(), pack, session)
		})
	})
}

var router cluster.INodeRouter

func on_router_block(client interface{}, data interface{}) {
	node := client.(cluster.INodeRouter)
	pack := data.(gnet.ISocketPacket)
	if pack.Cmd() == gnet.EVENT_HEARTBEAT_PINT {
		node.Push(gnet.NewPackArgs(gnet.EVENT_HEARTBEAT_PINT))
	}
}
