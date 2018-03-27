package gnet

import "net"
import "fmt"

//服务器启动
type INetServer interface {
	Start() bool
	Close()
}

//网络代理接口
type INetProxy interface {
	Run()
	OnClose()
}

func ListenAndRunServer(port int, block func(IBaseProxy)) INetServer {
	ser := NewTcpServer(port, func(conn interface{}) INetProxy {
		session := NewProxy(conn)
		block(session)
		return session
	})
	ser.Start()
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

//tcp链接
func Socket(addr string) (IConn, bool) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", conn.RemoteAddr().String())
		return NewConn(conn), true
	}
	fmt.Println("Socket Connect Err:", err)
	return nil, false
}

//运行
func RunWithContext(proxy INetProxy, conn IConn) {
	go func() {
		proxy.Run()
		proxy.OnClose()
	}()
}
