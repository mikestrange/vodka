package hippo

import "time"

type NetActor struct {
	caller    interface{}
	handle    IHandle
	openFlag  bool
	rdelay    int
	sdelay    int
	used      int
	send_chan chan interface{}
	read_chan chan interface{}
	msgr_chan chan IEvent
	exit_chan chan bool
}

func (this *NetActor) Make(size int) bool {
	if !this.openFlag {
		this.used = 0
		this.openFlag = true
		this.send_chan = make(chan interface{}, size)
		this.read_chan = make(chan interface{}, size)
		this.msgr_chan = make(chan IEvent, size)
		this.exit_chan = make(chan bool, size)
		return true
	}
	return false
}

func (this *NetActor) SetHandle(handle IHandle) {
	this.handle = handle
}

func (this *NetActor) SetCaller(val interface{}) {
	this.caller = val
}

func (this *NetActor) SetTimeout(val int) {
	this.rdelay = val
	this.sdelay = val
}

func (this *NetActor) SetReadTimeout(val int) {
	this.rdelay = val
}

func (this *NetActor) SetSendTimeout(val int) {
	this.sdelay = val
}

func (this *NetActor) PushRead(event interface{}) bool {
	if this.openFlag {
		this.read_chan <- event
		return true
	}
	return false
}

func (this *NetActor) PushSend(event interface{}) bool {
	if this.openFlag {
		this.send_chan <- event
		return true
	}
	return false
}

func (this *NetActor) PushEvent(event IEvent) bool {
	if this.openFlag {
		this.msgr_chan <- event
		return true
	}
	return false
}

func (this *NetActor) Exit(code int) bool {
	if this.openFlag {
		this.used = code
		this.openFlag = false
		close(this.exit_chan)
		return true
	}
	return false
}

func (this *NetActor) close_all() {
	close(this.send_chan)
	close(this.read_chan)
	close(this.msgr_chan)
}

func (this *NetActor) done(code int, data interface{}) {
	this.handle.Handle(code, this.caller, data)
}

func (this *NetActor) LoopWithTimeout() {
LOOP:
	for {
		select {
		case v := <-this.read_chan:
			this.done(EVENT_READ, v)
		case v := <-this.send_chan:
			this.done(EVENT_SEND, v)
		case v := <-this.msgr_chan:
			v.Perform()
		case <-time.After(time.Second * check_delay(this.rdelay)):
			this.done(EVENT_TIMEOUT, CODE_TIMEOUT_READ)
		case <-time.After(time.Second * check_delay(this.sdelay)):
			this.done(EVENT_TIMEOUT, CODE_TIMEOUT_SEND)
		case <-this.exit_chan:
			break LOOP
		default:
			//无
		}
	}
	//释放本身(防止中途重新使用)
	this.close_all()
	//释放后关闭
	this.done(EVENT_CLOSED, this.used)
}

func (this *NetActor) LoopWithNothing() {
LOOP:
	for {
		select {
		case v := <-this.read_chan:
			this.done(EVENT_READ, v)
		case v := <-this.send_chan:
			this.done(EVENT_SEND, v)
		case v := <-this.msgr_chan:
			v.Perform()
		case <-this.exit_chan:
			break LOOP
		default:
			//无
		}
	}
	//释放本身(防止中途重新使用)
	this.close_all()
	//释放后关闭
	this.done(EVENT_CLOSED, this.used)
}

//如果不设置，那么默认10分钟
func check_delay(val int) time.Duration {
	if val < 1 {
		return 60 * 10
	}
	return time.Duration(val)
}
