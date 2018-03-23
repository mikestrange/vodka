package actor

import "sync"
import "fmt"

type buffItem struct {
	closed bool
	args   []interface{}
}

type buffChannel struct {
	closeFlag bool
	buff      chan *buffItem
	mutex     sync.Mutex
}

func (this *buffChannel) init(sz int) {
	this.buff = make(chan *buffItem, sz)
}

func (this *buffChannel) Lock() {
	this.mutex.Lock()
}

func (this *buffChannel) Unlock() {
	this.mutex.Unlock()
}

func (this *buffChannel) Push(args ...interface{}) bool {
	return this.doPush(false, args...)
}

func (this *buffChannel) Pull() ([]interface{}, bool) {
	v, ok := <-this.buff
	if ok && !v.closed {
		return v.args, true
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

func (this *buffChannel) Loop(handle func(...interface{})) {
	for {
		if v, ok := this.Pull(); ok {
			handle(v...)
		} else {
			break
		}
	}
}

//private
func (this *buffChannel) doPush(closed bool, args ...interface{}) bool {
	this.Lock()
	defer this.Unlock()
	if this.closeFlag || this.isOverFull() {
		fmt.Println("is close or is full")
		return false
	}
	this.buff <- &buffItem{closed, args}
	return true
}

func (this *buffChannel) isOverFull() bool {
	return len(this.buff) == cap(this.buff)
}
