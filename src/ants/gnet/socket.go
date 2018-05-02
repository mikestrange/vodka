package gnet

import "net"
import "fmt"
import "ants/gcode"

//基础的代理
type TcpSocket struct {
	NetContext
	addr string
}

func NewSocket(addr string) *TcpSocket {
	this := new(TcpSocket)
	this.Connect(addr)
	return this
}

func (this *TcpSocket) Connect(addr string) bool {
	this.addr = addr
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", conn.RemoteAddr().String())
		this.SetConn(conn)
		this.Listen(this, 1024, 60*10) //默认3秒
		this.Conn().SetProcesser(gcode.NewClient())
		return true
	}
	fmt.Println("Socket Connect Err:", err)
	return false
}
