package kernel

type IGo interface {
	//Try(interface{}, ...interface{}) IGo
	Catch(func(interface{})) IGo
	Die(func()) IGo
	Done()
}

//可以被继承
type GoThread struct {
	target Throw
}

//新的线程
func Go(handle interface{}, args ...interface{}) IGo {
	this := &GoThread{}
	return this.Try(handle, args...)
}

//protected
func (this *GoThread) Try(handle interface{}, args ...interface{}) IGo {
	this.target.Try(func() {
		threadHandle(handle, args)
	})
	return this
}

//public
func (this *GoThread) Catch(block func(interface{})) IGo {
	this.target.Catch(block)
	return this
}

func (this *GoThread) Die(block func()) IGo {
	this.target.Die(block)
	return this
}

func (this *GoThread) Done() {
	go this.target.Done()
}

//支持3个参数
func threadHandle(fun interface{}, args []interface{}) {
	switch f := fun.(type) {
	case func():
		f()
	case func(interface{}):
		f(args[0])
	case func(interface{}, interface{}):
		f(args[0], args[1])
	case func(interface{}, interface{}, interface{}):
		f(args[0], args[1], args[2])
	case func(...interface{}):
		f(args...)
	default:
		panic("not handle func")
	}
}
