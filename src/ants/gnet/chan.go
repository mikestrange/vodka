package gnet

import "ants/gsys"

type buffSign struct {
	data   []byte
	closed bool
}

type buffChan struct {
	gsys.Locked
	closeFlag bool
	buff      chan *buffSign
}

func newChan(sz int) INetChan {
	this := new(buffChan)
	this.init(sz)
	return this
}

func (this *buffChan) init(sz int) {
	this.buff = make(chan *buffSign, sz)
	this.closeFlag = false
}

func (this *buffChan) Push(data []byte) bool {
	return this.doPush(false, data)
}

func (this *buffChan) Pull() ([]byte, bool) {
	v, ok := <-this.buff
	if ok && !v.closed {
		return v.data, true
	}
	return nil, false
}

//只会被关闭一次
func (this *buffChan) Close() {
	this.Lock()
	if !this.closeFlag {
		this.closeFlag = true
		close(this.buff)
	}
	this.Unlock()
}

func (this *buffChan) AsynClose() {
	this.doPush(true, nil)
}

func (this *buffChan) Loop(block func([]byte)) {
	for {
		if d, ok := this.Pull(); ok {
			block(d)
		} else {
			break
		}
	}
}

//push
func (this *buffChan) doPush(closed bool, data []byte) bool {
	this.Lock()
	if this.closeFlag {
		this.Unlock()
		return false
	}
	this.buff <- &buffSign{data: data, closed: closed}
	this.Unlock()
	return true
}

func (this *buffChan) isOverFull() bool {
	return len(this.buff) == cap(this.buff)
}
