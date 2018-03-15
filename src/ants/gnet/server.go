package gnet

//一个服务器如此简单
import (
	"fmt"
	"net"
	"sync"
	"time"
)

type ConnSet map[net.Conn]int

type TCPServer struct {
	Port           int
	NewProxyHandle ServerBlock
	//可选
	MaxConnNum int
	WriteNum   int
	//private
	ln         net.Listener
	conns      ConnSet
	mutexConns sync.Mutex
	wgLn       sync.WaitGroup
	wgConns    sync.WaitGroup
}

func NewTcpServer(port int, handle ServerBlock) INetServer {
	this := &TCPServer{Port: port, NewProxyHandle: handle}
	return this
}

func (this *TCPServer) Start() {
	if this.init() {
		this.run()
	}
}

func (this *TCPServer) init() bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", this.Port))
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Open Service:", this.Port)
	if this.MaxConnNum <= 0 {
		this.MaxConnNum = NET_SERVER_CONN_SIZE
	}
	if this.WriteNum <= 0 {
		this.WriteNum = NET_CHAN_SIZE
	}
	this.ln = ln
	this.conns = make(ConnSet)
	return true
}

//操作conn
func (this *TCPServer) closeMany(conn net.Conn) bool {
	this.mutexConns.Lock()
	defer this.mutexConns.Unlock()
	if this.ConnSize() >= this.MaxConnNum {
		conn.Close()
		return true
	}
	this.conns[conn] = this.ConnSize()
	return false
}

func (this *TCPServer) deleteConn(conn net.Conn) {
	this.mutexConns.Lock()
	delete(this.conns, conn)
	this.mutexConns.Unlock()
}

func (this *TCPServer) cleanConns() {
	this.mutexConns.Lock()
	for conn := range this.conns {
		conn.Close()
	}
	this.conns = nil
	this.mutexConns.Unlock()
}

func (this *TCPServer) ConnSize() int {
	return len(this.conns)
}

//运行
func (this *TCPServer) run() {
	this.wgLn.Add(1)
	defer this.wgLn.Done()
	for {
		conn, err := this.ln.Accept()
		if err != nil {
			if this.check_accept(err) {
				break
			}
			continue
		}
		if this.closeMany(conn) {
			fmt.Println("to many conn ;max is ", this.ConnSize())
		} else {
			this.handleConn(conn)
		}
	}
}

func (this *TCPServer) check_accept(err error) bool {
	var tempDelay time.Duration = 0
	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}
		fmt.Println("Accept error: ", err, "; retrying in ", tempDelay)
		time.Sleep(tempDelay)
		return false
	}
	fmt.Println("Accept Err:", err)
	return true
}

//关键处理
func (this *TCPServer) handleConn(conn net.Conn) {
	this.wgConns.Add(1)
	tcpConn := Context(conn, this.WriteNum)
	proxy := this.NewProxyHandle(tcpConn)
	go func() {
		proxy.Run()
		this.deleteConn(conn)
		proxy.OnClose(tcpConn.Close())
		this.wgConns.Done()
	}()
}

func (this *TCPServer) Close() {
	this.ln.Close()
	this.wgLn.Wait()
	this.cleanConns()
	this.wgConns.Wait()
}
