package command

//ants
import "ants/gutil"
import "ants/cluster"
import "ants/conf"
import "ants/gnet"
import "ants/gsys"

//基础模块
func SetMode(handle func(interface{}, ...interface{}), events map[int]interface{}, asyn bool) gutil.IModeAccessor {
	var group gsys.IAsynDispatcher = nil
	if asyn {
		group = gsys.NewChannel()
		go group.Loop(nil)
	}
	if handle == nil {
		handle = on_mode_block
	}
	mode := gutil.NewModeWithHandle(func(block interface{}, args ...interface{}) {
		if asyn {
			group.Push(func() {
				handle(block, args...)
			})
		} else {
			handle(block, args...)
		}
	})
	for k, v := range events {
		mode.On(k, v)
	}
	//心跳激活
	mode.On(gnet.EVENT_HEARTBEAT_PINT, func(proxy gnet.IBaseProxy) {
		//		if proxy.Context().AsSocket() {
		//			proxy.Send(gnet.NewPackArgs(gnet.EVENT_HEARTBEAT_PINT))
		//		} else {
		proxy.LivePing()
		//}
	})
	return mode
}

//路由节点
func SetRouter(port int, handle cluster.NodeBlock) cluster.INodeRouter {
	router := cluster.NewMainRouter(cluster.NewData(port), handle)
	conf.EachVo(func(vo *conf.RouteVo) {
		router.AddSet(cluster.NewRouter(cluster.NewDataWithVo(vo)))
	})
	return router
}

func on_mode_block(block interface{}, args ...interface{}) {
	switch block.(type) {
	case func(gnet.ISocketPacket):
		block.(func(gnet.ISocketPacket))(args[0].(gnet.ISocketPacket))
	case func(gnet.IBaseProxy, gnet.ISocketPacket):
		block.(func(gnet.IBaseProxy, gnet.ISocketPacket))(args[1].(gnet.IBaseProxy), args[0].(gnet.ISocketPacket))
	case func(gnet.IBaseProxy):
		block.(func(gnet.IBaseProxy))(args[1].(gnet.IBaseProxy))
	default:
		println("Err: no mode handle ")
	}
}
