package actor

import "ants/kernel"

//就这么简单
type BoxRef struct {
	//处理器
	actor IBoxActor
	//运行机器
	work kernel.WorkPusher
	//父级
	father IBoxSystem
}

//直接就运行了
func NewBox(target IBoxActor, val interface{}) IBoxRef {
	this := new(BoxRef)
	this.SetActor(target)
	RunAndThrowBox(this, val)
	return this
}

//需要make
func (this *BoxRef) Worker() kernel.IWorkPusher {
	return &this.work
}

//interfaces
func (this *BoxRef) Make(val interface{}) bool {
	_, ok := this.work.Make(val)
	return ok
}

func (this *BoxRef) OnReady() {
	//初始化(如果没有继承，那么默认最小)
	this.Make(nil)
}

func (this *BoxRef) Router(args ...interface{}) bool {
	return this.work.Push(args...)
}

func (this *BoxRef) SetActor(act IBoxActor) {
	this.actor = act
}

//一般在于重写他 >>启动在其他线程
func (this *BoxRef) PerformRunning() {
	//可以设置多个口
	this.work.ReadMsg(this)
}

func (this *BoxRef) Die() {
	this.work.Die()
}

func (this *BoxRef) OnReceiver(args ...interface{}) {
	//这里捕获错误吧
	this.actor.OnMessage(args...)
}

//被释放(基于父亲)
func (this *BoxRef) OnRelease() {
	this.actor.OnDie()
}

//一般情况下不使用
func (this *BoxRef) Father() IBoxSystem {
	return this.father
}

//不允许外界设置
func (this *BoxRef) setFather(father IBoxSystem) {
	this.father = father
}

func (this *BoxRef) Main() IBoxSystem {
	return main_actor
}

//独立运行的盒子1
func RunAndThrowBox(ref IBoxRef, val interface{}) IBoxRef {
	ref.Make(val)
	ref.OnReady()
	go func() {
		ref.PerformRunning()
		ref.OnRelease()
	}()
	return ref
}

//独立运行的盒子2
func RunAndThrowActor(ref IBoxRef, ator IBoxActor, val interface{}) IBoxRef {
	ref.SetActor(ator)
	return RunAndThrowBox(ref, val)
}
