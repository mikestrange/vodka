package core

import (
	"ants/base"
	"fmt"
	"sync"
	"sync/atomic"
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
	SetTimeout(interface{})
	//工作线程
	SetWorker(IWork)
}

/*
计时器Class
*/
type ClockTimer struct {
	live       bool
	delay      int
	retcount   int
	timeid     uint64
	clock_step uint64
	timer      *time.Ticker
	handle     TimeDelegate
	work       IWork
	_m         sync.Mutex
}

//基于线程
func newTimer(work IWork) ITimer {
	this := new(ClockTimer)
	this.SetWorker(work)
	return this
}

//private
func (this *ClockTimer) lock() {
	this._m.Lock()
}

func (this *ClockTimer) unlock() {
	this._m.Unlock()
}

//public
func (this *ClockTimer) SetTimeout(block interface{}) {
	this.handle = newDelegate(block)
}

func (this *ClockTimer) SetWorker(work IWork) {
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
	this.timeid, this.timer = this.create(delay, retcount, data)
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
		this.work.Put(func() {
			this.check_complete(idx, data)
		})
	}
}

func (this *ClockTimer) check_complete(idx uint64, data interface{}) {
	if this.check_timer(idx) {
		if this.handle == nil {
			fmt.Println("Warn: not step handle")
		} else {
			this.handle.OnTimeOutHandle(data)
		}
	} else {
		fmt.Println("Warn: step is stop")
	}
}

func (this *ClockTimer) check_timer(idx uint64) bool {
	this.lock()
	ok := this.timeid == idx
	this.unlock()
	return ok
}

func (this *ClockTimer) create(delay int, count int, arg interface{}) (uint64, *time.Ticker) {
	timeid := atomic.AddUint64(&this.clock_step, 1)
	timer := base.SetTimeout(delay, count, func() {
		this.complete(timeid, arg) //回调
	})
	return timeid, timer
}

//delegate interface
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
		panic("Unable to identify time handle")
	}
}
