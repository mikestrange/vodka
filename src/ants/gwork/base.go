package gwork

const THREAD_DEF_SIZE = 1   //并行默认(也就是单线程)
const QUEUE_MIN_SIZE = 100  //队列最小
const QUEUE_DEF_SIZE = 1000 //队列默认

//同步异步机制
type IAsynPusher interface {
	Push(...interface{}) bool
	Pull() ([]interface{}, bool)
}

//基本
type IWorkThread interface {
	IAsynPusher
	Open() bool
	Exit() bool
	//加入轮询，这里堵塞
	Join(func(...interface{}))
	//设置线程数目
	SetThreadNum(int)
	//设置每个线程队列数目
	SetMqNum(int)
}

func NewChannel() IWorkThread {
	this := NewWorks(0, 0)
	this.Open()
	return this
}

//其他形式
func Handle(block func(...interface{}), args []interface{}) {
	//第一个参数为函数就回调
	if len(args) > 0 {
		v := args[0]
		switch v.(type) {
		case func():
			v.(func())()
			return
		}
	}
	block(args...)
}
