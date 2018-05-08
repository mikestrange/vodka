package hippo

import "net"
import "fmt"

//协议
type Conn struct {
	conn net.Conn
}

func Dial(addr string) (interface{}, bool) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", conn.RemoteAddr().String())
		return conn, true
	}
	fmt.Println("Socket Connect Err:", err)
	return nil, false
}

func NewConn(conn interface{}) IConn {
	this := new(Conn)
	this.SetNet(conn)
	return this
}

func (this *Conn) SetNet(conn interface{}) {
	this.conn = conn.(net.Conn)
}

func (this *Conn) ReadBytes(b []byte) (int, error) {
	return this.conn.Read(b)
}

func (this *Conn) SendBytes(b []byte) bool {
	_, err := this.conn.Write(b)
	return err == nil
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
