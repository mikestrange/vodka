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
	var refLogic actor.IActorRef = actor.NewRefRunning(new(LogicActor))
	//服务器快速启动
	gnet.ListenAndRunServer(port, func(session gnet.IBaseProxy) {
		session.SetHandle(func(b []byte) {
			refLogic.Router(gnet.NewPackBytes(b), session)
		})
	})
}

//逻辑块
type LogicActor struct {
	actor.BaseActor
}

func (this *LogicActor) OnReady(ref actor.IActorRef) {
	ref.SetMqNum(5000)
	ref.SetThreadNum(1000)
	ref.Open()
}

func (this *LogicActor) OnClose() {

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
	UserID, GateID, SessionID := pack.ReadInt(), pack.ReadInt(), pack.ReadUInt64()
	//
	name := pack.ReadString()
	//改名
	fmt.Println("用户改名:", UserID, GateID, SessionID, name)
	if ok := change_name(int(UserID), name); ok {
		actor.Main.Send(int(GateID), pack_change_name(0, UserID, SessionID, name))
	} else {
		actor.Main.Send(int(GateID), pack_change_name(1, UserID, SessionID, ""))
	}
}

func on_change_info(pack gnet.ISocketPacket) {

}
