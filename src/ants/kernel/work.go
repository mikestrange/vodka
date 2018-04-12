package kernel

import "ants/gsys"

//others
//type ChanWork chan []interface{}
//type ChanSend chan<- []interface{}
//type ChanRead <-chan []interface{}

//信号
type workSign struct {
	code   int           //信号
	handle interface{}   //处理
	args   []interface{} //参数
}

type WorkPusher struct {
	openFlag bool
	work     chan interface{}
	locked   gsys.Locked
}

//do it
func NewWork() IWorkPusher {
	return new(WorkPusher)
}

//private
func (this *WorkPusher) lock() {
	this.locked.Lock()
}

func (this *WorkPusher) unlock() {
	this.locked.Unlock()
}

func (this *WorkPusher) Make(val interface{}) (interface{}, bool) {
	this.lock()
	defer this.unlock()
	return this.newChannel(val)
}

func (this *WorkPusher) newChannel(val interface{}) (interface{}, bool) {
	if this.openFlag {
		return val, false
	}
	switch c := val.(type) {
	case int:
		this.work = make(chan interface{}, check_queue_num(c))
	case chan interface{}:
		this.work = c
	default:
		return this.newChannel(QUEUE_MIN_SIZE)
	}
	this.openFlag = true
	return this.work, true
}

func (this *WorkPusher) Die() bool {
	this.lock()
	defer this.unlock()
	if this.openFlag {
		this.openFlag = false
		close(this.work)
		return false
	}
	return false
}

//外部直接调度
func (this *WorkPusher) Push(args ...interface{}) bool {
	return this.do_sign(&workSign{args: args})
}

func (this *WorkPusher) PushBlock(block func()) bool {
	return this.do_sign(&workSign{code: WORK_BLOCK_SIGN, handle: block})
}

//信号
func (this *WorkPusher) do_sign(data *workSign) bool {
	this.lock()
	if this.openFlag {
		this.work <- data
		this.unlock()
		return true
	}
	this.unlock()
	return false
}

//单线>基于阻塞
func (this *WorkPusher) ReadMsg(obser IWorkReceiver) {
	for data := range this.work {
		workArgsHandler(obser, data)
	}
}

//多线程执行>阻塞退出
func (this *WorkPusher) ReadRound(obser IWorkReceiver, size int) {
	wg := gsys.NewGroup()
	for i := 0; i < check_thread_num(size); i++ {
		wg.Wrap(func() {
			this.ReadMsg(obser)
		})
	}
	wg.Wait()
}

//环境计时器
func (this *WorkPusher) NewClock() ITimer {
	return NewClockWithWork(this)
}
