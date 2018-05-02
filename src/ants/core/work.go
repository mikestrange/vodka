package core

import "sync"

const THREAD_DEF_SIZE = 1   //并行默认
const QUEUE_MIN_SIZE = 100  //队列最小
const QUEUE_DEF_SIZE = 1000 //队列默认

func check_queue_num(size int) int {
	if size < QUEUE_MIN_SIZE {
		return QUEUE_DEF_SIZE
	}
	return size
}

//interface
type IWork interface {
	Make(interface{}) bool    //建立空间
	Close() bool              //关闭
	Put(interface{}) bool     //放消息
	Pop() (interface{}, bool) //取消息
	Recv() <-chan interface{} //接收
	Loop(func(interface{}))   //轮询
	NewTimer() ITimer         //基于时间
}

type Work struct {
	openFlag bool
	l        sync.Mutex
	c        chan interface{}
}

//protected
func (this *Work) Lock() {
	this.l.Lock()
}

func (this *Work) Unlock() {
	this.l.Unlock()
}

//继承可以使用
func (this *Work) NewChannel(val interface{}) bool {
	if this.openFlag {
		return false
	}
	switch c := val.(type) {
	case int:
		this.c = make(chan interface{}, check_queue_num(c))
	case chan interface{}:
		this.c = c
	default:
		return this.NewChannel(QUEUE_MIN_SIZE)
	}
	this.openFlag = true
	return true
}

//interface public
func (this *Work) Make(val interface{}) bool {
	this.Lock()
	ok := this.NewChannel(val)
	this.Unlock()
	return ok
}

func (this *Work) Close() bool {
	ok := false
	this.Lock()
	if this.openFlag {
		this.openFlag = false
		close(this.c)
		ok = true
	}
	this.Unlock()
	return ok
}

func (this *Work) Recv() <-chan interface{} {
	return this.c
}

func (this *Work) Put(data interface{}) bool {
	this.Lock()
	if this.openFlag {
		this.c <- data
	} else {
		//println("not live work")
	}
	this.Unlock()
	return this.openFlag
}

func (this *Work) Pop() (interface{}, bool) {
	v, ok := <-this.c
	if ok {
		return v, true
	}
	return nil, false
}

func (this *Work) Loop(block func(interface{})) {
	for v := range this.c {
		block(v)
	}
}

func (this *Work) NewTimer() ITimer {
	return newTimer(this)
}
