package kernel

type TimeDelegate interface {
	OnTimeOutHandle(interface{})
}

type FuncDelegate struct {
	handle interface{}
}

func newDelegate(block interface{}) TimeDelegate {
	return &FuncDelegate{handle: block}
}

func (this *FuncDelegate) OnTimeOutHandle(data interface{}) {
	switch f := this.handle.(type) {
	case func():
		f()
	case func(interface{}):
		f(data)
	case TimeDelegate:
		f.OnTimeOutHandle(data)
	default:
		println("not timer handle")
	}
}
