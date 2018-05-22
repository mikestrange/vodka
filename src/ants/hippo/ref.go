package hippo

import "sync"
import "time"

type IRef interface {
	OnDie()
	OnReady()
	Context() IContext
	SetHandle(IHandle)
	SetTimeout(int)
	SetCaller(interface{})
	Send(...interface{}) bool
	CloseOf(...interface{})
	Read(interface{})
	Run()
	Exit(int)
	Close()
	Failt()
	Err() int
}

//不重复利用
type NetRef struct {
	context    IContext
	caller     interface{}
	handle     IHandle
	openFlag   bool
	time_delay int
	send_chan  chan []interface{}
	read_chan  chan interface{}
	exit_chan  chan bool
	exit_code  int
	lock       sync.Mutex
}

func (this *NetRef) Lock() {
	this.lock.Lock()
}

func (this *NetRef) Unlock() {
	this.lock.Unlock()
}

func (this *NetRef) Make(val int) bool {
	this.Lock()
	if !this.openFlag {
		this.exit_code = 0
		this.openFlag = true
		this.send_chan = make(chan []interface{}, val)
		this.read_chan = make(chan interface{}, val)
		this.exit_chan = make(chan bool)
		this.Unlock()
		return true
	}
	this.Unlock()
	return false
}

func (this *NetRef) SetContext(tx IContext) {
	this.context = tx
}

func (this *NetRef) Context() IContext {
	return this.context
}

func (this *NetRef) SetHandle(val IHandle) {
	this.handle = val
}

func (this *NetRef) SetTimeout(delay int) {
	this.time_delay = delay
}

func (this *NetRef) SetCaller(val interface{}) {
	this.caller = val
}

func (this *NetRef) notice(code int, val interface{}) {
	this.handle.Handle(NewEvent(code, this.caller, val))
}

func (this *NetRef) Send(args ...interface{}) bool {
	this.Lock()
	if this.openFlag {
		this.send_chan <- args
		this.Unlock()
		return true
	}
	this.Unlock()
	return false
}

func (this *NetRef) Read(data interface{}) {
	this.Lock()
	if this.openFlag {
		this.read_chan <- data
	}
	this.Unlock()
}

func (this *NetRef) CloseOf(args ...interface{}) {
	this.Send(args...)
	this.Close()
}

//func (this *NetRef) ReadLoop() {
//	for {
//		if args, ok := this.context.ReadMsg(); ok {
//			for i := range args {
//				this.notice(EVENT_TYPE_READ, args[i])
//			}
//		} else {
//			break
//		}
//	}
//	this.Failt()
//}

func (this *NetRef) Run() {
	for {
		select {
		case v := <-this.send_chan:
			if !this.context.SendMsg(v...) {
				this.Failt()
			}
		case v := <-this.read_chan:
			this.notice(EVENT_TYPE_READ, v)
		case <-time.After(time.Second * check_delay(this.time_delay)):
			this.notice(EVENT_TYPE_TIMEOUT, nil)
		case <-this.exit_chan:
			break
		}
	}
	this.destroy()
	this.notice(EVENT_TYPE_CLOSE, this.exit_code)
}

func (this *NetRef) destroy() {
	close(this.read_chan)
	close(this.send_chan)
	this.context.Close()
}

func (this *NetRef) Exit(code int) {
	this.Lock()
	if this.openFlag {
		this.openFlag = true
		this.exit_code = code
		close(this.exit_chan)
	}
	this.Unlock()
}

func (this *NetRef) Close() {
	this.Exit(CLOSE_CODE_SELF)
}

func (this *NetRef) Failt() {
	this.Exit(CLOSE_CODE_ERROR)
}

func (this *NetRef) Err() int {
	return this.exit_code
}

func check_delay(val int) time.Duration {
	if val < 1 {
		return 1
	}
	return time.Duration(val)
}
