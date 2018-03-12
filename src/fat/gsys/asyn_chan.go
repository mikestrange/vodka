package gsys

import (
	"fmt"
	"sync"
)

type msgChan struct {
	size        int
	name        string
	used        int
	closed      bool
	handle      TaskBlock
	closeHandle CloseBlock
	//其他
	mutex *sync.Mutex
	work  chan interface{}
	exit  chan bool
}

//名称和尺寸
func newChanWithSize(size int, name string) IAsynDispatcher {
	this := new(msgChan)
	this.InitChanWithSize(size, name)
	return this
}

func (this *msgChan) InitChanWithSize(size int, name string) {
	this.size = size
	this.name = name
	this.used = 0
	this.closed = true
	this.handle = nil
	this.closeHandle = nil
	//
	this.mutex = new(sync.Mutex)
	this.work = nil
	this.exit = nil
}

func (this *msgChan) Lock() {
	this.mutex.Lock()
}

func (this *msgChan) Unlock() {
	this.mutex.Unlock()
}

//interfaces begin
func (this *msgChan) SetSize(val int) {
	this.size = val
}

func (this *msgChan) SetName(val string) {
	this.name = val
}

func (this *msgChan) Size() int {
	return this.size
}

func (this *msgChan) Name() string {
	return this.name
}

func (this *msgChan) SetHandle(handle TaskBlock) {
	this.handle = handle
}

func (this *msgChan) SetCloseHandle(closed CloseBlock) {
	this.closeHandle = closed
}

func (this *msgChan) Close() int {
	return this.CloseSign(0)
}

func (this *msgChan) CloseSign(used int) int {
	this.Lock()
	if false == this.closed {
		this.used = used
		this.closed = true
		close(this.exit)
	}
	this.Unlock()
	return this.used
}

//interfaces end

func (this *msgChan) Start() bool {
	this.Lock()
	if this.closed {
		this.used = 0
		this.closed = false
		this.exit = make(chan bool)
		this.work = make(chan interface{}, this.size)
		go this.thread(this.work, this.exit, this.handle, this.closeHandle)
	} else {
		fmt.Println(this.Name(), " is Running")
	}
	this.Unlock()
	return false == this.closed
}

func (this *msgChan) AsynPush(topic interface{}) {
	this.Lock()
	if this.closed {
		fmt.Println(this.Name(), " Closed Err Push")
	} else {
		this.work <- topic
	}
	this.Unlock()
}

//private functions
func (this *msgChan) thread(item chan interface{}, exit chan bool, handle TaskBlock, closed CloseBlock) {
	defer onAsynCloseHandle(closed)
	defer close(item)
	for {
		select {
		case data := <-item:
			onAsynTaskHandle(handle, data)
		case <-exit:
			break
		}
	}
}
