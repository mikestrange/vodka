package gsys

const ASYN_SIZE = 10000

//同步异步机制
type IAsynDispatcher interface {
	Push(interface{}) bool
	Pull() (interface{}, bool)
	Loop(func(interface{}))
	Close()
	AsynClose()
}

//处理异步回调
func handleAsynData(block func(interface{}), data interface{}) {
	switch data.(type) {
	case func():
		data.(func())()
	default:
		block(data)
	}
}

//#####################class
func NewChannel() IAsynDispatcher {
	return newChannel(ASYN_SIZE)
}

func NewBuffer() IAsynDispatcher {
	return newBuffer(ASYN_SIZE)
}
