package gnet

import "net"
import "ants/gcode"
import "ants/base"

type IConn interface {
	//1,协议处理
	SetProcesser(gcode.IByteCoder)
	//处理
	Coder() gcode.IByteCoder
	//直接关闭
	Close()
	//关闭并且写入
	CloseOf(...interface{})
	//异步写
	Send(...interface{}) bool
	//同步写入
	WriteBytes([]byte)
	//读消息
	ReadMsg([]byte) ([]interface{}, int)
	//轮询
	Run()
	//远端地址
	Remote() string
	//近端地址
	Local() string
	//local
	SetTimeout(int)
}

//网络环境实例
type Conn struct {
	NetBase
	conn    net.Conn
	process gcode.IByteCoder
}

//protected 继承能获得
func (this *Conn) SetConn(conn interface{}) {
	this.conn = conn.(net.Conn)
	this.process = gcode.NewProtocoler() //默认
}

// interface INetConn
func (this *Conn) SetProcesser(val gcode.IByteCoder) {
	this.process = val
}

func (this *Conn) Coder() gcode.IByteCoder {
	return this.process
}

func (this *Conn) Send(args ...interface{}) bool {
	if bits, err := this.Coder().Marshal(args...); err == nil {
		return this.Write(bits)
	}
	return false
}

func (this *Conn) Close() {
	this.Exit(CLOSE_SIGN_SELF)
}

func (this *NetContext) CloseOf(args ...interface{}) {
	this.Send(args...)
	this.Close()
}

//conn oper
func (this *Conn) Remote() string {
	return this.conn.RemoteAddr().String()
}

func (this *Conn) Local() string {
	return this.conn.LocalAddr().String()
}

func (this *Conn) ReadMsg(b []byte) ([]interface{}, int) {
	if ret, err := this.conn.Read(b); err == nil {
		if list, ok := this.Coder().Unmarshal(b[:ret]); ok == nil {
			return list, CLOSE_SIGN_SUCCESS
		} else {
			return nil, CLOSE_SIGN_ERROR
		}
	}
	return nil, CLOSE_SIGN_CLIENT
}

func (this *Conn) Run() {
	bit := make([]byte, this.Coder().BuffSize())
	for {
		list, code := this.ReadMsg(bit)
		if code == CLOSE_SIGN_SUCCESS {
			for i := range list {
				this.Read(list[i])
			}
		} else {
			this.Exit(code)
			break
		}
	}
}

func (this *Conn) WriteBytes(b []byte) {
	this.conn.Write(b)
}

//handle
func (this *Conn) OnDestroy() {
	this.conn.Close()
}

func (this *Conn) OnMessage(code int, data interface{}) {
	switch code {
	case EVENT_CONN_CLOSE:
		//关闭处理
	case EVENT_CONN_READ:

	case EVENT_CONN_SEND:
		//发送
		this.WriteBytes(base.ToBytes(data))
	case EVENT_CONN_HEARTBEAT:
		//心跳
	case EVENT_CONN_SIGN:
		//其他事务
	}
}
