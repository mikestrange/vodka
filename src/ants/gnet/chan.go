package gnet

import "time"

type IHandle interface {
	OnMessage(int, interface{})
	OnDestroy()
}

type Channels struct {
	read   chan interface{}
	send   chan interface{}
	done   chan interface{}
	exit   chan int
	used   int
	delay  int
	handle IHandle
}

func newChannel(size int) *Channels {
	this := new(Channels)
	this.Make(size)
	return this
}

func (this *Channels) Make(size int) {
	this.read = make(chan interface{}, size)
	this.send = make(chan interface{}, size)
	this.done = make(chan interface{}, size)
	this.exit = make(chan int)
}

func (this *Channels) Read(v interface{}) {
	this.read <- v
}

func (this *Channels) Send(v interface{}) {
	this.send <- v
}

func (this *Channels) Done(v interface{}) {
	this.done <- v
}

func (this *Channels) CloseOf(v int) {
	this.used = v
	close(this.exit)
	//this.exit <- v
}

func (this *Channels) SetTimeout(delay int) {
	this.delay = delay
}

func (this *Channels) SetHandle(handle IHandle) {
	this.handle = handle
}

func (this *Channels) Close() {
	close(this.read)
	close(this.send)
	close(this.done)
	//close(this.exit)
}

//有超时时间的轮询
func (this *Channels) LoopWithTimeout() {
	this.handle.OnMessage(EVENT_CONN_CONNECT, nil)
	for {
		select {
		case v := <-this.read:
			this.handle.OnMessage(EVENT_CONN_READ, v)
		case v := <-this.send:
			this.handle.OnMessage(EVENT_CONN_SEND, v)
		case v := <-this.done:
			this.handle.OnMessage(EVENT_CONN_SIGN, v)
		case <-this.exit:
			this.handle.OnMessage(EVENT_CONN_CLOSE, this.used)
			goto End
		case <-time.After(time.Second * time.Duration(check_delay(this.delay))):
			this.handle.OnMessage(EVENT_CONN_HEARTBEAT, nil)
		}
	}
End:
	{
		this.Close()
		this.handle.OnDestroy()
	}
}

//无超时时间的轮询
func (this *Channels) Loop() {
	this.handle.OnMessage(EVENT_CONN_CONNECT, nil)
	for {
		select {
		case v := <-this.read:
			this.handle.OnMessage(EVENT_CONN_READ, v)
		case v := <-this.send:
			this.handle.OnMessage(EVENT_CONN_SEND, v)
		case v := <-this.done:
			this.handle.OnMessage(EVENT_CONN_SIGN, v)
		case <-this.exit:
			this.handle.OnMessage(EVENT_CONN_CLOSE, this.used)
			goto End
		}
	}
End:
	{
		this.Close()
		this.handle.OnDestroy()
	}
}

func check_delay(val int) int {
	if val < 1 {
		//如果没有，那么默认10分钟心跳
		return 60 * 10
	}
	return val
}
