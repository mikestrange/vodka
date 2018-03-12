package gsys

//###########################异步任务
type TaskBlock func(interface{})
type CloseBlock func()

//同步异步机制
type IAsynDispatcher interface {
	Name() string
	Size() int
	Close() int
	Start() bool
	CloseSign(int) int
	AsynPush(interface{})
	SetHandle(TaskBlock)
	SetCloseHandle(CloseBlock)
}

//任务处理(这里应该处理错误，免得影响其他工作)
func onAsynTaskHandle(handle TaskBlock, data interface{}) {
	switch data.(type) {
	case func():
		data.(func())()
	default:
		if handle != nil {
			handle(data)
		} else {
			println("No Handle Data:", data)
		}
	}
}

//关闭处理
func onAsynCloseHandle(closed CloseBlock) {
	if closed != nil {
		closed()
	}
}

//#####################class
//实例化一个通道(默认使用系统方式)
func NewChannel() IAsynDispatcher {
	return NewChannelWithSize(0, 10000, "nil")
}

func NewChannelWithSize(mtype int, size int, name string) IAsynDispatcher {
	if mtype == 1 {
		return newBufferWithSize(size, name)
	}
	return newChanWithSize(size, name)
}
