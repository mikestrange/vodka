package gnet

import "ants/gtime"

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
	return false
}

//interface INetProxy
func (this *BaseProxy) Run() {
	tm := gtime.After(PING_TIME, func() {
		if this.Ping() {
			this.Send(NewPackArgs(EVENT_HEARTBEAT_PINT))
		} else {
			this.Close()
		}
	})
	defer tm.Stop()
	this.WaitFor()
}

func (this *BaseProxy) OnClose() {

}
