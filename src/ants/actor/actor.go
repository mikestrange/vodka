package actor

//对象
type IActor interface {
	OnReady(IActorRef)
	OnMessage(...interface{}) //处理对象
	OnClose()                 //退出调度
}

//为了封装
type BaseActor struct {
}

func (this *BaseActor) OnReady(ref IActorRef) {
	//准备阶段，需要对它进行操作
	ref.Open() //打开ref
}

func (this *BaseActor) OnMessage(args ...interface{}) {
	println("test msg:", len(args))
}

func (this *BaseActor) OnClose() {
	println("test close")
}

//函数回调
type FuncActor struct {
	BaseActor
	onMsgr func(...interface{})
	onEnd  func()
}

func NewFunc(msgr func(...interface{}), end func()) IActor {
	return &FuncActor{onMsgr: msgr, onEnd: end}
}

func (this *FuncActor) OnMessage(args ...interface{}) {
	this.onMsgr(args...)
}

func (this *FuncActor) OnClose() {
	if this.onEnd != nil {
		this.onEnd()
	}
}
