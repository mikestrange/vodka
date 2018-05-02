package gnet

import "time"
import "ants/base"

type IHandle interface {
	OnMessage(int, interface{})
	OnDestroy()
}

//
type INetBase interface {
	//Listen(IHandle, int, int) bool
	Done(interface{}) bool
	Write([]byte) bool
	Read([]byte) bool
	Exit(int) int
}

//网络基础
type NetBase struct {
	INetBase
	base.Locked
	openFlag bool
	delay    int
	handle   IHandle
	used     int
	rc       chan []byte
	sc       chan []byte
	ec       chan int
	dc       chan interface{}
}

func newBase() INetBase {
	return new(NetBase)
}

func (this *NetBase) Listen(handle IHandle, size int, delay int) bool {
	ok := false
	this.Lock()
	if !this.openFlag {
		ok = true
		this.delay = delay
		this.handle = handle
		this.rc = make(chan []byte, size)
		this.sc = make(chan []byte, size)
		this.dc = make(chan interface{}, size)
		this.ec = make(chan int)
		this.openFlag = true
		go this.LoopWork()
	}
	this.Unlock()
	return ok
}

func (this *NetBase) LoopWork() {
	LoopHandle(this.handle, this.rc, this.sc, this.dc, this.ec, this.delay)
}

func (this *NetBase) Write(b []byte) bool {
	ok := false
	this.Lock()
	if this.openFlag {
		ok = true
		this.sc <- b
	}
	this.Unlock()
	return ok
}

func (this *NetBase) Read(b []byte) bool {
	ok := false
	this.Lock()
	if this.openFlag {
		ok = true
		this.rc <- b
	}
	this.Unlock()
	return ok
}

func (this *NetBase) Done(event interface{}) bool {
	ok := false
	this.Lock()
	if this.openFlag {
		ok = true
		this.dc <- event
	}
	this.Unlock()
	return ok
}

func (this *NetBase) Exit(code int) int {
	this.Lock()
	if this.openFlag {
		this.used = code
		this.openFlag = false
		this.ec <- code
		//close(this.ec)
	}
	this.Unlock()
	return this.used
}

//最终
func LoopHandle(handle IHandle, read chan []byte, send chan []byte, done chan interface{}, exit chan int, delay int) {
	handle.OnMessage(EVENT_CONN_CONNECT, nil)
	for {
		select {
		case v := <-read:
			handle.OnMessage(EVENT_CONN_READ, v)
		case v := <-send:
			handle.OnMessage(EVENT_CONN_SEND, v)
		case v := <-done:
			handle.OnMessage(EVENT_CONN_SIGN, v)
		case v := <-exit:
			handle.OnMessage(EVENT_CONN_CLOSE, v)
			goto End
		case <-time.After(time.Second * time.Duration(delay)):
			handle.OnMessage(EVENT_CONN_HEARTBEAT, nil)
		}
	}
End:
	{
		close(read)
		close(send)
		close(done)
		close(exit)
		handle.OnDestroy()
	}
}
func init() {
	println("init")
	t := newBase()
	t.Exit(1)
	//t.Listen(t, 100, 1)
}
