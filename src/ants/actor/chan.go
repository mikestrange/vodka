package actor

import "sync"
import "fmt"

type buffItem struct {
	closed bool
	data   interface{}
}

type buffChannel struct {
	closeFlag bool
	buff      chan *buffItem
	mutex     sync.Mutex
}

//func newBuffer(sz int) *buffChannel {
//	this := new(buffChannel)
//	this.init(sz)
//	return this
//}

func (this *buffChannel) init(sz int) {
	this.buff = make(chan *buffItem, sz)
}

func (this *buffChannel) Lock() {
	this.mutex.Lock()
}

func (this *buffChannel) Unlock() {
	this.mutex.Unlock()
}

func (this *buffChannel) Push(data interface{}) bool {
	return this.doPush(false, data)
}

func (this *buffChannel) Pull() (interface{}, bool) {
	v, ok := <-this.buff
	if ok && !v.closed {
		return v.data, true
	}
	return nil, false
}

func (this *buffChannel) Close() {
	this.Lock()
	if !this.closeFlag {
		this.closeFlag = true
		close(this.buff)
	}
	this.Unlock()
}

func (this *buffChannel) Loop(handle func(interface{})) {
	for {
		if v, ok := this.Pull(); ok {
			handle(v)
		} else {
			break
		}
	}
}

//private
func (this *buffChannel) doPush(closed bool, data interface{}) bool {
	this.Lock()
	defer this.Unlock()
	if this.closeFlag || this.isOverFull() {
		fmt.Println("is close or is full")
		return false
	}
	this.buff <- &buffItem{closed, data}
	return true
}

func (this *buffChannel) isOverFull() bool {
	return len(this.buff) == cap(this.buff)
}
