package command

//ants
import "ants/gutil"
import "ants/gnet"

//基础模块
func SetMode(handle func(interface{}, ...interface{}), events map[int]interface{}) gutil.IModeAccessor {
	if handle == nil {
		handle = on_mode_block
	}
	mode := gutil.NewModeWithHandle(func(block interface{}, args ...interface{}) {
		handle(block, args...)
	})
	for k, v := range events {
		mode.On(k, v)
	}
	//心跳激活
	mode.On(gnet.EVENT_HEARTBEAT_PINT, func(proxy gnet.IBaseProxy) {
		proxy.LivePing()
	})
	return mode
}

//pack 0, session 1
func on_mode_block(block interface{}, args ...interface{}) {
	switch block.(type) {
	case func(gnet.ISocketPacket):
		block.(func(gnet.ISocketPacket))(args[0].(gnet.ISocketPacket))
	case func(gnet.ISocketPacket, gnet.IBaseProxy):
		block.(func(gnet.ISocketPacket, gnet.IBaseProxy))(args[0].(gnet.ISocketPacket), args[1].(gnet.IBaseProxy))
	case func(gnet.IBaseProxy):
		block.(func(gnet.IBaseProxy))(args[1].(gnet.IBaseProxy))
	case func(gnet.IBaseProxy, gnet.ISocketPacket):
		block.(func(gnet.IBaseProxy, gnet.ISocketPacket))(args[1].(gnet.IBaseProxy), args[0].(gnet.ISocketPacket))
	default:
		println("Err: no mode handle ")
	}
}
