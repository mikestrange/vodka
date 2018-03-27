package chat

import "app/command"

//非网关模块通用
import "ants/gnet"
import "ants/actor"

//服务器的启动(快速启动)
func ServerLaunch(port int) {
	//模块调度
	var refLogic *ChatActor = new(ChatActor)
	refLogic.init()
	//服务器快速启动
	gnet.ListenAndRunServer(port, func(session gnet.IBaseProxy) {
		session.SetHandle(func(b []byte) {
			refLogic.Router(gnet.NewPackBytes(b), session)
		})
	})
}

//处理模块
type ChatActor struct {
	actor.ActorSystem
	actor.BaseActor
}

func (this *ChatActor) init() {
	this.Init()
	this.SetActor(this)
	actor.RunWithActor(this)
}

func (this *ChatActor) OnReady(ref actor.IActorRef) {
	//test
	this.ActorOf(10086, NewTable(10086, 0))
	//open
	ref.Open()
}

func (this *ChatActor) OnMessage(args ...interface{}) {
	pack := args[0].(gnet.ISocketPacket)
	switch pack.Cmd() {
	case command.SERVER_BUILD_CHANNEL:
		this.on_build_channel(pack.ReadInt(), pack.ReadByte())
	case command.SERVER_REMOVE_CHANNEL:
		this.on_remove_channel(pack.ReadInt())
	default:
		this.on_table_message(pack)
	}
}

func (this *ChatActor) on_build_channel(cid int32, ctype int8) {
	this.ActorOf(int(cid), NewTable(cid, ctype))
}

func (this *ChatActor) on_remove_channel(cid int32) {
	this.Remove(int(cid))
}

func (this *ChatActor) on_table_message(pack gnet.ISocketPacket) {
	header := NewHeader(pack)
	cid := pack.ReadInt()
	this.Send(int(cid), header, pack)
}

func (this *ChatActor) OnClose() {
	this.Shutdown()
}
