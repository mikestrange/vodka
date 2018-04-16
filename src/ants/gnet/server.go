package gnet

//一个服务器如此简单
import (
	"ants/gsys"
	"fmt"
	"net"
)

type ServerBlock func(interface{}) INetProxy
type connSet map[net.Conn]int

type TCPServer struct {
	gsys.Locked
	//必须
	Port       int
	ConnHandle ServerBlock
	//可选
	MaxConnNum int
	//private
	ln     net.Listener
	conns  connSet
	wgSer  gsys.WaitGroup
	wgConn gsys.WaitGroup
}

func NewTcpServer(port int, handle ServerBlock) INetServer {
	this := &TCPServer{Port: port, ConnHandle: handle}
	return this
}

func (this *TCPServer) Start() bool {
	if this.init() {
		go this.run()
		return true
	}
	return false
}

func (this *TCPServer) init() bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", this.Port))
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Open Service:", this.Port)
	if this.MaxConnNum <= 0 {
		this.MaxConnNum = NET_SERVER_CONN_SIZE
	}
	this.ln = ln
	this.conns = make(connSet)
	return true
}

//操作conn
func (this *TCPServer) checkMany(conn net.Conn) bool {
	this.Lock()
	if this.ConnSize() >= this.MaxConnNum {
		this.Unlock()
		conn.Close()
		return true
	}
	this.conns[conn] = this.ConnSize()
	this.Unlock()
	return false
}

func (this *TCPServer) deleteConn(conn net.Conn) {
	this.Lock()
	delete(this.conns, conn)
	this.Unlock()
	//避免其他地方没关闭
	conn.Close()
}

func (this *TCPServer) cleanConns() {
	this.Lock()
	for conn := range this.conns {
		conn.Close()
	}
	this.Unlock()
}

func (this *TCPServer) ConnSize() int {
	return len(this.conns)
}

func (this *TCPServer) Close() {
	this.ln.Close()
	this.wgSer.Wait()
	this.cleanConns()
	this.wgConn.Wait()
}

//运行
func (this *TCPServer) run() {
	this.wgSer.Add()
	defer this.wgSer.Done()
	for {
		conn, err := this.ln.Accept()
		if err == nil {
			if this.checkMany(conn) {
				fmt.Println("To many conn ;max is ", this.ConnSize())
			} else {
				this.handleConn(conn)
			}
		} else {
			if check_server_error(err) {
				break
			}
		}
	}
}

func (this *TCPServer) handleConn(conn net.Conn) {
	proxy := this.ConnHandle(conn)
	this.wgConn.Wrap(func() {
		proxy.Run()
		this.deleteConn(conn)
		proxy.OnClose()
	})
}
