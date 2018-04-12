package kernel

//静态方法接收
type FuncReceiver struct {
	handle func(...interface{})
}

func NewReceiver(fun func(...interface{})) IWorkReceiver {
	return &FuncReceiver{fun}
}

func (this *FuncReceiver) OnReceiver(args ...interface{}) {
	this.handle(args...)
}
