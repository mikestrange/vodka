package world

//非网关模块通用(状态服务器)
import "ants/gnet"
import "app/command"

import "ants/actor"
import "ants/gutil"

//服务器的启动(快速启动)
func ServerLaunch(port int) {
	refLogic := actor.RunAndThrowBox(new(LogicActor), nil)
	//服务器快速启动
	gnet.ListenAndRunServer(port, func(session gnet.IBaseProxy) {
		session.SetHandle(func(b []byte) {
			refLogic.Router(gnet.NewPackBytes(b), session)
		})
	})
}

//逻辑块(单线)
type LogicActor struct {
	actor.BaseBox
	mode gutil.IModeAccessor
}

func (this *LogicActor) OnReady() {
	this.mode = command.SetMode(nil, events)
	this.SetActor(this)
}

func (this *LogicActor) OnMessage(args ...interface{}) {
	pack := args[0].(gnet.ISocketPacket)
	this.mode.Done(pack.Cmd(), args...)
}

func (this *LogicActor) OnDie() {

}
