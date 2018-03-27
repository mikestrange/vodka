package world

//非网关模块通用(状态服务器)
import "ants/gnet"
import "app/command"

import "ants/actor"
import "ants/gutil"

//服务器的启动(快速启动)
func ServerLaunch(port int) {
	var refLogic actor.IActorRef = actor.NewRefRunning(new(LogicActor))
	//服务器快速启动
	gnet.ListenAndRunServer(port, func(session gnet.IBaseProxy) {
		session.SetHandle(func(b []byte) {
			refLogic.Router(gnet.NewPackBytes(b), session)
		})
	})
}

//逻辑块(单线)
type LogicActor struct {
	actor.BaseActor
	mode gutil.IModeAccessor
}

func (this *LogicActor) OnReady(ref actor.IActorRef) {
	this.mode = command.SetMode(nil, events)
	ref.Open()
}

func (this *LogicActor) OnMessage(args ...interface{}) {
	pack := args[0].(gnet.ISocketPacket)
	this.mode.Done(pack.Cmd(), args...)
}

func (this *LogicActor) OnClose() {

}
