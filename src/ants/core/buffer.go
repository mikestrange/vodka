package core

//缓冲通道(可以用于发送)
import (
	"sync"
)

//缓存节点
type node struct {
	data interface{}
	next *node
}

//工作缓冲区
type WorkBuffer struct {
	openFlag bool
	max_size int
	current  int
	bitem    *node
	eitem    *node
	cond     *sync.Cond
}

//protected
func (this *WorkBuffer) Lock() {
	this.cond.L.Lock()
}

func (this *WorkBuffer) Unlock() {
	this.cond.L.Unlock()
}

//interface
func (this *WorkBuffer) Make(val int) bool {
	this.Lock()
	ok := this.NewChannel(val)
	this.Unlock()
	return ok
}

func (this *WorkBuffer) NewChannel(val int) bool {
	if this.openFlag {
		return false
	}
	this.openFlag = true
	this.current = 0
	this.max_size = val
	return true
}

func (this *WorkBuffer) Die() bool {
	this.Lock()
	if this.openFlag {
		this.openFlag = false
		this.cond.Broadcast()
	}
	this.Unlock()
	return true
}

func (this *WorkBuffer) Put(event interface{}) bool {
	this.Lock()
	//over
	if !this.openFlag || this.current >= this.max_size {
		this.Unlock()
		return false
	}
	val := &node{event, nil}
	if this.current == 0 {
		this.bitem, this.eitem = val, val
	} else {
		this.eitem.next, this.eitem = val, val
	}
	this.current++
	this.cond.Signal()
	this.Unlock()
	return true
}

func (this *WorkBuffer) Pop() (interface{}, bool) {
	for {
		this.Lock()
		if this.current > 0 {
			val := this.bitem.data
			this.bitem = this.bitem.next
			this.current--
			this.Unlock()
			return val, true
		}
		if !this.openFlag {
			this.Unlock()
			return nil, false
		}
		this.cond.Wait()
		this.Unlock()
	}
}
