package actor

import "fmt"

type IActor interface {
	WillReceive(interface{}) bool
	//扩展就好了
	OnMessage(...interface{})
	OnClosed()
	//下面的不需要改变
	Close()
	Idx() int
	Topics() []int
	Name() string
	Runner() IRunner
	Commit(...interface{}) bool
	Context() IActorSystem
}

//所有的都可以继承它
type BaseActor struct {
	idx    int
	name   string
	runner IRunner
	system IActorSystem
}

//简单的
func NewTestActor() IActor {
	this := new(BaseActor)
	return this.SetMaster(1, "test actor", nil, nil)
}

//protected (必须要设置的项目, 在添加之前)
func (this *BaseActor) SetMaster(idx int, name string, runner IRunner, system IActorSystem) IActor {
	this.idx = idx
	this.name = name
	if runner == nil {
		runner = newRunner()
	}
	this.runner = runner
	if system == nil {
		system = Main
	}
	this.system = system
	return this
}

//interfaces
//继承扩展
func (this *BaseActor) OnMessage(args ...interface{}) {
	fmt.Println(this.name, " >actor message", args)
}

//释放资源的时候继承它,只会被执行一次(所以释放的时候继承它)
func (this *BaseActor) OnClosed() {

}

//校正合法性
func (this *BaseActor) WillReceive(data interface{}) bool {
	return true
}

//constant functions
func (this *BaseActor) Commit(args ...interface{}) bool { //constant
	return this.Runner().Send(args...)
}

func (this *BaseActor) Name() string { //constant
	return this.name
}

func (this *BaseActor) Topics() []int {
	return []int{}
}

func (this *BaseActor) Idx() int { //constant
	return this.idx
}

func (this *BaseActor) Context() IActorSystem { //constant
	return this.system
}

func (this *BaseActor) Runner() IRunner { //constant
	return this.runner
}

func (this *BaseActor) Close() { //constant
	this.Runner().Close()
}
