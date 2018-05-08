package gnet

import "ants/base"
import "ants/gcode"

//基础环境(轻量级别)
type Context interface {
	IAgent
	SetReceiver(func([]byte))
	Close()
	Send(...interface{}) bool
	CloseOf(...interface{})
	//Conn() IConn
	LivePing()
	Ping() bool
}

type NetContext struct {
	Conn
	ping   bool //是否关闭
	handle func([]byte)
}

func NewContext() Context {
	return new(NetContext)
}

func (this *NetContext) SetReceiver(block func([]byte)) {
	this.handle = block
}

//interfaces
func (this *NetContext) OnReady(conn interface{}) {
	this.SetConn(conn)
	this.SetTimeout(3) //1秒不登录断开
	this.Listen(this, 100)
}

func (this *NetContext) OnDie() {
	println("context close")
}

func (this *NetContext) Wait() {
	//todo
}

func (this *NetContext) Run() {
	base.Wraps(func() {
		this.Conn.Run()
	}, func() {
		this.Loop()
		//this.LoopTimeout()
	})
}

//interfaces
func (this *NetContext) LivePing() {
	this.ping = false
}

func (this *NetContext) Ping() bool {
	if this.ping {
		return false
	}
	this.ping = true
	return true
}

func (this *NetContext) DoHandle(v []byte) {
	this.handle(v)
}

//
func (this *NetContext) OnMessage(code int, data interface{}) {
	if code == EVENT_CONN_READ {
		this.handle(base.ToBytes(data))
	} else if code == EVENT_CONN_HEARTBEAT {
		this.Send(gcode.NewPackArgs(EVENT_CONN_HEARTBEAT))
	} else if code == EVENT_CONN_SEND {
		this.WriteBytes(base.ToBytes(data))
	}
}
