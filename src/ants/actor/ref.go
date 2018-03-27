package actor

import "ants/gwork"

//代理接口
type IActorRef interface {
	SetActor(IActor)
	ActorObj() IActor
	Router(...interface{}) bool
	Close()
	Run()
	//继承
	Open() bool
	SetMqNum(int)
	SetThreadNum(int)
}

//代理类:能实现对进程
type ActorRef struct {
	gwork.WorkSelector
	obj IActor
}

//适合多态
func NewRef() *ActorRef {
	this := new(ActorRef)
	return this
}

//直接运行(适合单个Actor)
func NewRefRunning(obj IActor) IActorRef {
	this := new(ActorRef)
	this.SetActor(obj)
	RunWithActor(this)
	return this
}

//interfaces
func (this *ActorRef) SetActor(obj IActor) {
	this.obj = obj
}

func (this *ActorRef) ActorObj() IActor {
	return this.obj
}

func (this *ActorRef) Run() {
	this.Join(this.obj.OnMessage)
}

func (this *ActorRef) Router(args ...interface{}) bool {
	return this.Push(args...)
}

func (this *ActorRef) Close() {
	this.Exit()
}

//最后的挣扎
func RunWithActor(ref IActorRef) {
	obj := ref.ActorObj()
	obj.OnReady(ref)
	go func() {
		ref.Run()
		obj.OnClose()
	}()
}
