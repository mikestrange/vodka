package gnet

import "ants/base"
import "ants/gcode"

//基础环境(轻量级别)
type Context interface {
	IAgent
	SetReceiver(func([]byte))
	Close()
	Send(...interface{})
	CloseOf(...interface{})
	Conn() IConn
	LivePing()
	Ping() bool
}

type NetContext struct {
	conn   Conn
	ping   bool //是否关闭
	handle func([]byte)
}

func NewContext(conn interface{}) Context {
	this := new(NetContext)
	this.SetConn(conn)
	return this
}

func (this *NetContext) SetConn(conn interface{}) {
	this.conn.SetConn(conn)
	this.conn.Listen(this, 1024, 60*10)
}

func (this *NetContext) SetReceiver(block func([]byte)) {
	this.handle = block
}

func (this *NetContext) OnDie() {
	println("context close")
}

func (this *NetContext) Wait() {
	//todo
}

func (this *NetContext) Run() {
	this.conn.Run()
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

func (this *NetContext) Conn() IConn {
	return &this.conn
}

func (this *NetContext) Send(args ...interface{}) {
	this.conn.Send(args...)
}

func (this *NetContext) Close() {
	this.conn.Close()
}

func (this *NetContext) CloseOf(args ...interface{}) {
	this.conn.Send(args...)
	this.conn.Close()
}

func (this *NetContext) OnDestroy() {
	this.conn.OnDestroy()
}

func (this *NetContext) OnMessage(code int, data interface{}) {
	if code == EVENT_CONN_READ {
		this.handle(base.ToBytes(data))
	} else if code == EVENT_CONN_HEARTBEAT {
		this.Send(gcode.NewPackArgs(EVENT_CONN_HEARTBEAT))
	} else if code == EVENT_CONN_SEND {
		this.conn.WriteBytes(base.ToBytes(data))
	}
}
