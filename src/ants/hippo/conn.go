package hippo

import "net"
import "time"

//具体网络(重写而已)
type IConn interface {
	Read([]byte) ([]byte, bool)
	Write([]byte) bool
	Close() bool
	//地址
	Local() string
	Remote() string
	//超时
	SetReadDeadDelay(int)
	SetWriteDeadDelay(int)
	SetDeadDelay(int)
}

//简单的封装一下
type Conn struct {
	conn net.Conn
}

//环境实例
func newContext(conn net.Conn) IContext {
	return &Context{conn: newConn(conn), code: NewBigCoding()}
}

//一个连接
func newConn(conn net.Conn) IConn {
	return &Conn{conn: conn}
}

func (this *Conn) Read(b []byte) ([]byte, bool) {
	ret, err := this.conn.Read(b)
	if err == nil {
		return b[:ret], true
	}
	return nil, false
}

func (this *Conn) Write(b []byte) bool {
	_, err := this.conn.Write(b)
	if err == nil {
		return true
	}
	return false
}

func (this *Conn) Close() bool {
	return this.conn.Close() == nil
}

func (this *Conn) Local() string {
	return this.conn.LocalAddr().String()
}

func (this *Conn) Remote() string {
	return this.conn.RemoteAddr().String()
}

func (this *Conn) SetReadDeadDelay(delay int) {
	this.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(delay)))
}

func (this *Conn) SetWriteDeadDelay(delay int) {
	this.conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(delay)))
}

func (this *Conn) SetDeadDelay(delay int) {
	this.conn.SetDeadline(time.Now().Add(time.Second * time.Duration(delay)))
}
