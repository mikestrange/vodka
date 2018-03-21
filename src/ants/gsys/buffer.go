package gsys

//缓冲通道(可以用于发送)
import (
	"sync"
)

type buffItem struct {
	data interface{}
	next *buffItem
}

//消息队列
type msgBuffer struct {
	size      int
	closeFlag bool
	//其他
	current int
	bitem   *buffItem
	eitem   *buffItem
	cond    *sync.Cond
}

func newBuffer(size int) IAsynDispatcher {
	var mutex sync.Mutex
	this := &msgBuffer{size: size, closeFlag: false, current: 0, cond: sync.NewCond(&mutex)}
	return this
}

func (this *msgBuffer) Close() {
	this.cond.L.Lock()
	if !this.closeFlag {
		this.closeFlag = true
		this.cond.Broadcast()
	}
	this.cond.L.Unlock()
}

func (this *msgBuffer) AsynClose() {
	//println("The buff do not asyn close!")
}

//目前不用
func (this *msgBuffer) Start() bool {
	this.cond.L.Lock()
	if this.closeFlag {
		this.closeFlag = false
		this.current = 0
	}
	this.cond.L.Unlock()
	return this.closeFlag == false
}

func (this *msgBuffer) Push(data interface{}) bool {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	if this.closeFlag || this.current >= this.size {
		return false
	}
	val := &buffItem{data, nil}
	if this.current == 0 {
		this.bitem, this.eitem = val, val
	} else {
		this.eitem.next, this.eitem = val, val
	}
	this.current++
	this.cond.Signal()
	return true
}

func (this *msgBuffer) Pull() (interface{}, bool) {
	for {
		this.cond.L.Lock()
		if this.current > 0 {
			val := this.bitem.data
			this.bitem = this.bitem.next
			this.current--
			this.cond.L.Unlock()
			return val, true
		}
		if this.closeFlag {
			this.cond.L.Unlock()
			return nil, false
		}
		this.cond.Wait()
		this.cond.L.Unlock()
	}
}

func (this *msgBuffer) Loop(block func(interface{})) {
	for {
		if d, ok := this.Pull(); ok {
			handleAsynData(block, d)
		} else {
			break
		}
	}
}
