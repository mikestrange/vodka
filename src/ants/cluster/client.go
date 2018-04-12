package cluster

import "net"
import "fmt"

type TcpClient struct {
}

func (this *TcpClient) connect(addr string) (net.Conn, bool) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", conn.RemoteAddr().String())
		return conn, true
	}
	fmt.Println("Socket Connect Err:", err)
	return nil, false
}

func (this *TcpClient) OnMessage(args ...interface{}) {
	//这里就直接去发送，不接受回执
}
