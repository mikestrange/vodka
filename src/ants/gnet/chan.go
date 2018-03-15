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
	this := &buffChan{buff: make(chan *buffSign, sz), closeFlag: false}
	return this
}

func (this *buffChan) Push(data []byte) bool {
	return this.done_push(false, data)
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
	this.done_push(true, nil)
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
func (this *buffChan) done_push(closed bool, data []byte) bool {
	this.Lock()
	// len(this.buff) == cap(this.buff) 满了
	if this.closeFlag {
		this.Unlock()
		return false
	}
	this.buff <- &buffSign{data: data, closed: closed}
	this.Unlock()
	return true
}
