package kernel

import (
	"ants/gsys"
	"fmt"
	"time"
)

/*
计时器接口
*/
type ITimer interface {
	Delay() int
	Retcount() int
	TimeID() uint64
	//动作
	Stop()                       //停止
	Reset(interface{})           //按照之前的方式运行
	Start(int, int, interface{}) //开始运行
	Forever(int, interface{})
	Onec(int, interface{})
	//回调
	SetHandle(interface{})
	//工作线程
	SetWorker(IWorkPusher)
}

/*
计时器Class
*/
type ClockTimer struct {
	live     bool
	delay    int
	retcount int
	timeid   uint64
	timer    *time.Ticker
	handle   TimeDelegate
	work     IWorkPusher
	locked   gsys.Locked
}

func NewTimer() ITimer {
	this := new(ClockTimer)
	return this
}

func NewClockWithWork(work IWorkPusher) ITimer {
	this := NewTimer()
	this.SetWorker(work)
	return this
}

//private
func (this *ClockTimer) lock() {
	this.locked.Lock()
}

func (this *ClockTimer) unlock() {
	this.locked.Unlock()
}

//public
func (this *ClockTimer) SetHandle(block interface{}) {
	this.handle = newDelegate(block)
}

func (this *ClockTimer) SetWorker(work IWorkPusher) {
	this.work = work
}

func (this *ClockTimer) Delay() int {
	return this.delay
}

func (this *ClockTimer) Retcount() int {
	return this.retcount
}

func (this *ClockTimer) TimeID() uint64 {
	return this.timeid
}

//handle
func (this *ClockTimer) Reset(data interface{}) {
	this.Start(this.delay, this.retcount, data)
}

func (this *ClockTimer) Forever(delay int, data interface{}) {
	this.Start(delay, 0, data)
}

func (this *ClockTimer) Onec(delay int, data interface{}) {
	this.Start(delay, 1, data)
}

func (this *ClockTimer) Start(delay int, retcount int, data interface{}) {
	this.Stop()
	this.lock()
	this.live = true
	this.delay = delay
	this.retcount = retcount
	this.timeid, this.timer = clockHandler(delay, retcount, data, this.complete)
	this.unlock()
}

func (this *ClockTimer) Stop() {
	this.lock()
	if this.live {
		this.timeid = 0
		this.live = false
		this.timer.Stop()
	}
	this.unlock()
}

//private functions (这里不判断是否结束)
func (this *ClockTimer) complete(idx uint64, data interface{}) {
	if this.work == nil {
		this.check_complete(idx, data)
	} else {
		this.work.PushBlock(func() {
			this.check_complete(idx, data)
		})
	}
}

func (this *ClockTimer) check_complete(idx uint64, data interface{}) {
	if this.check_timer(idx) {
		if this.handle == nil {
			fmt.Println("No Timer Hanle")
		} else {
			this.handle.OnTimeOutHandle(data)
		}
	} else {
		fmt.Println("Timer work Over")
	}
}

func (this *ClockTimer) check_timer(idx uint64) bool {
	this.lock()
	ok := this.timeid == idx
	this.unlock()
	return ok
}
