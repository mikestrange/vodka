package actor

//static
//静态函数回调
type ActFunc struct {
	onMsgr func(...interface{})
	onEnd  func()
}

func NewFunc(block func(...interface{}), die func()) IBoxActor {
	return &ActFunc{onMsgr: block, onEnd: die}
}

func (this *ActFunc) OnMessage(args ...interface{}) {
	this.onMsgr(args...)
}

func (this *ActFunc) OnDie() {
	if this.onEnd != nil {
		this.onEnd()
	}
}

//继承(未实现)
type UnTypeActor struct {
	IBoxActor
}

func (this *UnTypeActor) OnMessage(args ...interface{}) {}

func (this *UnTypeActor) OnDie() {}

//集成
type BaseBox struct {
	BoxRef
	UnTypeActor
}
