package gnet

import "net"
import "ants/gcode"

//关闭命令
type workChan chan []interface{}

type IConn interface {
	//1,协议处理
	SetProcesser(gcode.IByteCoder)
	//直接关闭
	Close()
	//异步写
	Send(...interface{}) bool
	//写入
	WriteBytes([]byte)
	//轮询
	Run()
	//远端地址
	Remote() string
	//近端地址
	Local() string
}

//网络环境实例
type Conn struct {
	NetBase
	conn  net.Conn
	coder gcode.IByteCoder
}

//protected 继承能获得
func (this *Conn) SetConn(conn interface{}) {
	this.conn = conn.(net.Conn)
	this.coder = gcode.NewProtocoler() //默认
}

// interface INetConn
func (this *Conn) SetProcesser(val gcode.IByteCoder) {
	this.coder = val
}

func (this *Conn) Send(args ...interface{}) bool {
	decode := this.coder
	if bits, err := decode.Marshal(args...); err == nil {
		return this.Write(bits)
	}
	return false
}

func (this *Conn) Close() {
	this.Exit(CLOSE_SIGN_SELF)
}

//conn oper
func (this *Conn) Remote() string {
	return this.conn.RemoteAddr().String()
}

func (this *Conn) Local() string {
	return this.conn.LocalAddr().String()
}

func (this *Conn) Run() {
	decode := this.coder
	bit := make([]byte, decode.BuffSize())
	for {
		if ret, err := this.conn.Read(bit); err == nil {
			if list, ok := decode.Unmarshal(bit[:ret]); ok == nil {
				for i := range list {
					this.Read(list[i])
				}
			} else {
				this.Exit(CLOSE_SIGN_ERROR)
				break
			}
		} else {
			this.Exit(CLOSE_SIGN_CLIENT)
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

	case EVENT_CONN_READ:

	case EVENT_CONN_SEND:

	case EVENT_CONN_HEARTBEAT:

	case EVENT_CONN_SIGN:

	}
}
