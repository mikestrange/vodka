package gnet

import "ants/gsys"

//基础的(继承他就好了)
type IBaseProxy interface {
	//代理对象
	IConn
	//作为代理
	INetProxy
	//提供了一些基础
	LivePing()
	Ping() bool
}

type BaseProxy struct {
	NetConn
	pingFlag bool //是否关闭
}

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
	if !this.pingFlag {
		this.pingFlag = true
		return true
	}
	this.Close()
	return false
}

//interface INetProxy
func (this *BaseProxy) Run() {
	tm := gsys.After(PING_TIME, func() {
		this.Send(NewPackArgs(EVENT_HEARTBEAT_PINT))
	})
	defer tm.Stop()
	this.WaitFor()
}

func (this *BaseProxy) OnClose() {

}
