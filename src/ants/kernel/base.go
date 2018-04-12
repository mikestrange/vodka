package kernel

const THREAD_DEF_SIZE = 1   //并行默认
const QUEUE_MIN_SIZE = 100  //队列最小
const QUEUE_DEF_SIZE = 1000 //队列默认

const (
	WORK_NONE_SIGN  = 0
	WORK_BLOCK_SIGN = 1
)

func check_thread_num(size int) int {
	if size < THREAD_DEF_SIZE {
		return THREAD_DEF_SIZE
	}
	return size
}

func check_queue_num(size int) int {
	if size < QUEUE_MIN_SIZE {
		return QUEUE_DEF_SIZE
	}
	return size
}

//推送者
type IWorkPusher interface {
	//建立一个通道(关闭的时候不允许)
	Make(interface{}) (interface{}, bool)
	//关闭
	Die() bool
	//推送
	Push(...interface{}) bool
	//推送回调函数
	PushBlock(func()) bool
	//推送信号
	//PushSign(int, interface{}, ...interface{}) bool
	//单线
	ReadMsg(IWorkReceiver)
	//多线程(一个处理器并行)
	ReadRound(IWorkReceiver, int)
	//计时器
	NewClock() ITimer
}

//消费者
type IWorkReceiver interface {
	OnReceiver(...interface{})
}

//------------------------static------------------------------
func workArgsHandler(obser IWorkReceiver, data interface{}) error {
	switch sign := data.(*workSign); sign.code {
	case WORK_BLOCK_SIGN: //回调函数
		sign.handle.(func())()
	default:
		obser.OnReceiver(sign.args...)
	}
	return nil
}
