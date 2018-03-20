package gnet

import "net"
import "fmt"

type INetServer interface {
	Start()
	Close()
}

//网络代理接口
type INetProxy interface {
	Run()
	OnClose()
}

func ListenAndRunServer(port int, block func(IBaseProxy)) INetServer {
	ser := NewTcpServer(port, func(conn interface{}) {
		session := NewBaseProxy(Context(conn))
		block(session)
		session.Run()
		session.OnClose()
	})
	go ser.Start()
	return ser
}

//快速启动服务器
func ListenRunServer(port int) (net.Listener, bool) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Println("Run Err:", err)
		return nil, false
	}
	fmt.Println("Run Ser:", port)
	return ln, true
}

//直接使用代理就会运行
func Context(conn interface{}) INetContext {
	return NewConn(conn)
}

func Socket(addr string) (INetContext, bool) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", conn.LocalAddr().String())
		return NewConn(conn), true
	}
	fmt.Println("Socket Connect Err:", err)
	return nil, false
}
