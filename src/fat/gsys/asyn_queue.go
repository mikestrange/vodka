package gsys

//缓冲通道(可以用于发送)
import (
	"fmt"
	"sync"
)

//元素
type buffItem struct {
	data interface{}
	next *buffItem
}

type msgBuffer struct {
	size        int
	name        string
	used        int
	closed      bool
	handle      TaskBlock
	closeHandle CloseBlock
	//其他
	current int
	bitem   *buffItem
	eitem   *buffItem
	cond    *sync.Cond
}

func newBufferWithSize(size int, name string) IAsynDispatcher {
	var mutex sync.Mutex
	this := new(msgBuffer)
	this.InitBuffWithSize(size, name, &mutex)
	return this
}

func (this *msgBuffer) InitBuffWithSize(size int, name string, mutex sync.Locker) {
	this.size = size
	this.name = name
	this.used = 0
	this.closed = true
	this.handle = nil
	this.closeHandle = nil
	//
	this.bitem = nil
	this.eitem = nil
	this.cond = sync.NewCond(mutex)
}

//
func (this *msgBuffer) SetSize(val int) {
	this.size = val
}

func (this *msgBuffer) SetName(val string) {
	this.name = val
}

func (this *msgBuffer) Size() int {
	return this.size
}

func (this *msgBuffer) Name() string {
	return this.name
}

func (this *msgBuffer) SetHandle(handle TaskBlock) {
	this.handle = handle
}

func (this *msgBuffer) SetCloseHandle(closed CloseBlock) {
	this.closeHandle = closed
}

func (this *msgBuffer) Close() int {
	return this.CloseSign(0)
}

func (this *msgBuffer) CloseSign(used int) int {
	this.cond.L.Lock()
	if false == this.closed {
		this.used = used
		this.closed = true
		this.cond.Broadcast()
	}
	this.cond.L.Unlock()
	return this.used
}

func (this *msgBuffer) AsynPush(data interface{}) {
	this.put(data)
}

func (this *msgBuffer) Start() bool {
	this.cond.L.Lock()
	if this.closed {
		this.closed = false
		this.current = 0
		this.used = 0
		go this.thread(this.handle, this.closeHandle)
	}
	this.cond.L.Unlock()
	return this.closed == false
}

func (this *msgBuffer) put(data interface{}) bool {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	//TODO
	if this.closed || this.current >= this.size {
		if this.closed {
			fmt.Println("This Buffer is Closed")
		} else {
			fmt.Println("This Buffer is OverFull", this.current, this.size)
		}
		return false
	}
	val := &buffItem{data, nil}
	if this.current == 0 {
		this.bitem, this.eitem = val, val
	} else {
		this.eitem.next, this.eitem = val, val
	}
	this.current = this.current + 1
	this.cond.Signal()
	return true
}

func (this *msgBuffer) pop() (interface{}, bool) {
	for {
		this.cond.L.Lock()
		if this.current > 0 {
			val := this.bitem.data
			this.bitem = this.bitem.next
			this.current = this.current - 1
			this.cond.L.Unlock()
			return val, true
		}
		if this.closed {
			this.cond.L.Unlock()
			return nil, false
		}
		this.cond.Wait()
		this.cond.L.Unlock()
	}
}

func (this *msgBuffer) thread(handle TaskBlock, closed CloseBlock) {
	defer onAsynCloseHandle(closed)
	for {
		if data, ok := this.pop(); ok {
			onAsynTaskHandle(handle, data)
		} else {
			break
		}
	}
}
