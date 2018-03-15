package gsys

import (
	"fmt"
	"time"
)

//stop方法必须在回调一个线程(因为如果提前停掉可能会导致致命问题)
type TimerBlock func(interface{})

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
	//处理
	SetHandle(TimerBlock)
	SetChannel(IAsynDispatcher)
}

/*
计时器Class
*/
type ClockTimer struct {
	Locked
	live     bool
	delay    int
	retcount int
	timeid   uint64
	timer    *time.Ticker
	handle   TimerBlock
	channel  IAsynDispatcher
}

func NewTimer() ITimer {
	this := new(ClockTimer)
	this.InitClockTimer()
	return this
}

func NewTimerWithChannel(target IAsynDispatcher) ITimer {
	this := NewTimer()
	this.SetChannel(target)
	return this
}

func (this *ClockTimer) InitClockTimer() {

}

func (this *ClockTimer) SetHandle(block TimerBlock) {
	this.handle = block
}

func (this *ClockTimer) SetChannel(channel IAsynDispatcher) {
	this.channel = channel
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
	this.Lock()
	this.live = true
	this.delay = delay
	this.retcount = retcount
	this.timeid, this.timer = clockHandler(delay, retcount, data, this.complete)
	this.Unlock()
}

func (this *ClockTimer) Stop() {
	this.Lock()
	if this.live {
		this.timeid = 0
		this.live = false
		this.timer.Stop()
	}
	this.Unlock()
}

//private functions (这里不判断是否结束)
func (this *ClockTimer) complete(idx uint64, data interface{}) {
	if this.channel == nil {
		this.check_complete(idx, data)
	} else {
		this.channel.Push(func() {
			this.check_complete(idx, data)
		})
	}
}

func (this *ClockTimer) check_complete(idx uint64, data interface{}) {
	if this.check_timer(idx) {
		if this.handle == nil {
			fmt.Println("No Timer Hanle")
		} else {
			this.handle(data)
		}
	} else {
		fmt.Println("Timer channel Over")
	}
}

func (this *ClockTimer) check_timer(idx uint64) bool {
	this.Lock()
	defer this.Unlock()
	return this.timeid == idx
}
