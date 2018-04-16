package hall

import "fmt"
import "app/command"

//弱连接服务器，不用管心跳
import "ants/gnet"
import "ants/actor"

//服务器的启动(快速启动)
func ServerLaunch(port int) {
	init_dber()
	//模块调度
	refLogic := actor.RunAndThrowBox(new(LogicActor), nil)
	//服务器快速启动
	gnet.ListenAndRunServer(port, func(session gnet.IBaseProxy) {
		session.SetHandle(func(b []byte) {
			refLogic.Router(gnet.NewPackBytes(b), session)
		})
	})
}

//逻辑块
type LogicActor struct {
	actor.BaseBox
}

func (this *LogicActor) OnReady() {
	this.SetActor(this)
}

func (this *LogicActor) OnDie() {

}

func (this *LogicActor) PerformRunning() {
	this.Worker().ReadRound(this, 1000)
}

//message
func (this *LogicActor) OnMessage(args ...interface{}) {
	pack := args[0].(gnet.ISocketPacket)
	switch pack.Cmd() {
	case command.CLIENT_CHANGE_NAME:
		on_change_name(pack)
	case command.CLIENT_CHANGE_INFO:
		on_change_info(pack)
	}
}

func on_change_name(pack gnet.ISocketPacket) {
	uid, GateID, SessionID := pack.ReadInt(), pack.ReadInt(), pack.ReadUInt64()
	//
	name := pack.ReadString()
	//改名
	fmt.Println("用户改名:", uid, GateID, SessionID, name)
	if ok := change_name(uid, name); ok {
		actor.Main().Send(GateID, pack_change_name(0, uid, SessionID, name))
	} else {
		actor.Main().Send(GateID, pack_change_name(1, uid, SessionID, ""))
	}
}

func on_change_info(pack gnet.ISocketPacket) {

}
