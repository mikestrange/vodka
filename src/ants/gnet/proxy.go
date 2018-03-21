package gnet

import "ants/gsys"

//基础的(继承他就好了)
type IBaseProxy interface {
	INetProxy
	//提供了一些基础
	Tx() INetContext
	Send(interface{})
	CloseWrite() //异步关闭
	Close()      //直接关闭
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

func NewProxyConn(conn interface{}) IBaseProxy {
	this := &BaseProxy{tx: NewConn(conn), pingFlag: false}
	return this
}

func (this *BaseProxy) SetContext(tx INetContext) {
	this.tx = tx
}

//interfaces IBaseProxy
func (this *BaseProxy) Tx() INetContext {
	return this.tx
}

func (this *BaseProxy) Send(data interface{}) {
	this.tx.Send(data)
}

func (this *BaseProxy) CloseWrite() {
	this.tx.CloseWrite()
}

func (this *BaseProxy) Close() {
	this.tx.Close()
}

func (this *BaseProxy) LivePing() {
	this.pingFlag = false
}

func (this *BaseProxy) Ping() bool {
	if !this.pingFlag {
		this.pingFlag = true
		return true
	}
	this.tx.Close()
	return false
}

//interface INetProxy
func (this *BaseProxy) Run() {
	tx := this.Tx()
	t := gsys.After(PING_TIME, func() {
		tx.Send(NewPackArgs(EVENT_HEARTBEAT_PINT))
	})
	defer t.Stop()
	tx.WaitFor()
}

func (this *BaseProxy) OnClose() {

}
