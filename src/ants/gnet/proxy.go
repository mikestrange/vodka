package gnet

import "ants/kernel"

//基础的(继承他就好了)
type IBaseProxy interface {
	//代理对象
	Context
	//作为代理
	INetProxy
	//提供了一些基础
	LivePing()
	Ping() bool
}

type BaseProxy struct {
	NetConn
	pingFlag bool //是否关闭
	handle   func([]byte)
}

//代理Conn
func NewProxy(conn interface{}) IBaseProxy {
	this := &BaseProxy{pingFlag: false}
	this.SetConn(conn)
	return this
}

//interfaces IBaseProxy
func (this *BaseProxy) LivePing() {
	this.pingFlag = false
}

func (this *BaseProxy) Ping() bool {
	if this.pingFlag {
		return false
	}
	this.pingFlag = true
	return true
}

//interface INetProxy
func (this *BaseProxy) Run() {
	tm := kernel.After(PING_TIME, func() {
		if this.Ping() {
			this.Send(NewPackArgs(EVENT_HEARTBEAT_PINT))
		} else {
			this.Close()
		}
	})
	defer tm.Stop()
	this.Join(this.on_read_handle)
}

func (this *BaseProxy) SetHandle(block func([]byte)) {
	this.handle = block
}

func (this *BaseProxy) OnClose() {

}

//private
func (this *BaseProxy) on_read_handle(b []byte) {
	if this.handle == nil {
		println("no read handle")
	} else {
		this.handle(b)
	}
}
