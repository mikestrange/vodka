package gnet

import "ants/gsys"

//基础的
type IBaseProxy interface {
	INetProxy
	//提供了一些基础
	Context() INetContext
	Send(...interface{})
	Close()
	Kill()
	LivePing()
	Ping() bool
}

type BaseProxy struct {
	tx        INetContext
	heartbeat bool
}

func NewBaseProxy(tx INetContext) IBaseProxy {
	this := new(BaseProxy)
	this.SetContext(tx)
	return this
}

func (this *BaseProxy) Context() INetContext {
	return this.tx
}

func (this *BaseProxy) SetContext(tx INetContext) {
	this.tx = tx
}

func (this *BaseProxy) Send(args ...interface{}) {
	this.tx.Send(args...)
}

func (this *BaseProxy) Close() {
	this.tx.Close()
}

func (this *BaseProxy) Kill() {
	this.tx.Shutdown(SIGN_CLOSE_OK)
}

func (this *BaseProxy) LivePing() {
	this.heartbeat = true
}

func (this *BaseProxy) Ping() bool {
	if this.heartbeat {
		this.heartbeat = false
		return true
	}
	this.tx.CloseSign(SIGN_CLOSE_HEARTBEAT)
	return false
}

//interface INetProxy
func (this *BaseProxy) Run() {
	if this.tx.AsSocket() {
		this.tx.Run()
	} else {
		t := gsys.After(PING_TIME, func() {
			this.Send(NewPackArgs(EVENT_HEARTBEAT_PINT))
		})
		defer t.Stop()
		this.tx.Run()
	}
}

func (this *BaseProxy) OnClose(code int) {
	//只关心被关闭
}
