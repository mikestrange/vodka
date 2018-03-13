package command

import "fat/gutil"
import "fat/gnet/nsc"
import "fat/gnet"
import "app/config"
import "fat/gsys"

//这个可以设置为单线程
func SetRouter(port int, events map[int]interface{}, group gsys.IAsynDispatcher) (nsc.IRemoteScheduler, gutil.IModeAccessor) {
	if group != nil {
		group.Start()
	}
	//设置回调
	mode := gutil.NewModeWithHandle(func(block interface{}, cmd int, args ...interface{}) {
		//时间回执是否异步
		if group == nil {
			on_mode_handle(block, args[0], args[1])
		} else {
			group.AsynPush(func() {
				on_mode_handle(block, args[0], args[1])
			})
		}
	})
	//注册监听
	if events != nil {
		for k, v := range events {
			mode.On(k, v)
		}
	}
	//心跳激活
	mode.On(gnet.GNET_HEARTBEAT_PINT, func(tx gnet.INetContext) {
		tx.LivePing()
		//println("Heart Beat Set Live")
	})
	//建立一个路由节点
	router := nsc.NewRemoteScheduler(config.GetDataRouter(port), func(rotue nsc.IRouter, data interface{}) {
		//心跳处理
		if data.(gnet.ISocketPacket).Cmd() == gnet.GNET_HEARTBEAT_PINT {
			rotue.Push(gnet.PacketWithHeartBeat)
		}
	})
	config.SetServerLists(router)
	return router, mode
}

func on_mode_handle(block interface{}, tx interface{}, pack interface{}) {
	switch block.(type) {
	case func(gnet.INetContext, gnet.ISocketPacket):
		block.(func(gnet.INetContext, gnet.ISocketPacket))(tx.(gnet.INetContext), pack.(gnet.ISocketPacket))
	case func(gnet.ISocketPacket):
		block.(func(gnet.ISocketPacket))(pack.(gnet.ISocketPacket))
	case func(gnet.INetContext):
		block.(func(gnet.INetContext))(tx.(gnet.INetContext))
	}
}
