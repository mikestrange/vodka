package gnet

//一个服务器如此简单
import (
	"ants/gsys"
	"fmt"
	"net"
	"sync"
	"time"
)

type ServerBlock func(interface{}) INetProxy
type ConnSet map[net.Conn]int

type TCPServer struct {
	gsys.Locked
	Port       int
	ConnHandle ServerBlock
	//可选
	MaxConnNum int
	//private
	ln      net.Listener
	conns   ConnSet
	wgLn    sync.WaitGroup
	wgConns sync.WaitGroup
}

func NewTcpServer(port int, handle ServerBlock) INetServer {
	this := &TCPServer{Port: port, ConnHandle: handle}
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
	this.ln = ln
	this.conns = make(ConnSet)
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
	//conn.Close()
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
	this.wgLn.Wait()
	this.cleanConns()
	this.wgConns.Wait()
}

//运行
func (this *TCPServer) run() {
	this.wgLn.Add(1)
	defer this.wgLn.Done()
	for {
		conn, err := this.ln.Accept()
		if err == nil {
			if this.checkMany(conn) {
				fmt.Println("To many conn ;max is ", this.ConnSize())
			} else {
				this.wgConns.Add(1)
				this.handleConn(conn)
			}
		} else {
			if !this.check_accept(err) {
				break
			}
		}
	}
}

func (this *TCPServer) handleConn(conn net.Conn) {
	agent := this.ConnHandle(conn)
	//
	go func() {
		agent.Run()
		this.deleteConn(conn)
		agent.OnClose()
		this.wgConns.Done()
	}()
}

//没有什么意义的一段
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
	} else {
		fmt.Println("Accept Err:", err)
		return false
	}
	return true
}
