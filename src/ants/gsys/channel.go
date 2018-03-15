package gsys

import "sync"

type chanItem struct {
	closed bool
	data   interface{}
}

type buffChan struct {
	Locked
	closeFlag bool
	buff      chan *chanItem
	mutex     sync.Mutex
}

func newChannel(sz int) IAsynDispatcher {
	this := new(buffChan)
	this.init(sz)
	return this
}

func (this *buffChan) init(sz int) {
	this.buff = make(chan *chanItem, sz)
}

func (this *buffChan) Push(data interface{}) bool {
	return this.done_push(false, data)
}

func (this *buffChan) Pull() (interface{}, bool) {
	v, ok := <-this.buff
	if ok && !v.closed {
		return v.data, true
	}
	return nil, false
}

func (this *buffChan) AsynClose() {
	this.done_push(true, nil)
}

func (this *buffChan) Close() {
	this.mutex.Lock()
	if !this.closeFlag {
		this.closeFlag = true
		close(this.buff)
	}
	this.mutex.Unlock()
}

func (this *buffChan) Loop(block func(interface{})) {
	for {
		if d, ok := this.Pull(); ok {
			handleAsynData(block, d)
		} else {
			break
		}
	}
}

//private
func (this *buffChan) done_push(closed bool, data interface{}) bool {
	this.mutex.Lock()
	if this.closeFlag {
		this.mutex.Unlock()
		return false
	}
	this.buff <- &chanItem{closed, data}
	this.mutex.Unlock()
	return true
}
