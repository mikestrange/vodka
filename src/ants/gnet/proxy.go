package gnet

import "ants/gsys"

//基础的(继承他就好了)
type IBaseProxy interface {
	INetProxy
	//提供了一些基础
	Context() INetContext
	Send(...interface{})
	Close() //异步关闭
	Kill()  //直接关闭
	LivePing()
	Ping() bool
}

type BaseProxy struct {
	tx       INetContext
	pingFlag bool //是否关闭
}

func NewBaseProxy(tx INetContext) IBaseProxy {
	this := &BaseProxy{tx: tx, pingFlag: false}
	return this
}

func (this *BaseProxy) SetContext(tx INetContext) {
	this.tx = tx
}

//interfaces IBaseProxy
func (this *BaseProxy) Context() INetContext {
	return this.tx
}

func (this *BaseProxy) Send(args ...interface{}) {
	this.Context().Send(args...)
}

func (this *BaseProxy) Close() {
	this.Context().Close()
}

func (this *BaseProxy) Kill() {
	this.Context().Shutdown(SIGN_CLOSE_OK)
}

func (this *BaseProxy) LivePing() {
	this.pingFlag = false
}

func (this *BaseProxy) Ping() bool {
	if !this.pingFlag {
		this.pingFlag = true
		return true
	}
	this.Context().CloseSign(SIGN_CLOSE_HEARTBEAT)
	return false
}

//interface INetProxy
func (this *BaseProxy) Run() {
	tx := this.Context()
	if tx.AsSocket() {
		tx.Run()
	} else {
		t := gsys.After(PING_TIME, func() {
			tx.Send(NewPackArgs(EVENT_HEARTBEAT_PINT))
		})
		defer t.Stop()
		tx.Run()
	}
}

func (this *BaseProxy) OnClose(code int) {
	//只关心被关闭
}
