package gnet

import (
	"ants/base"
	"ants/core"
	"fmt"
	"net"
	"sync"
	"time"
)

type NewAgent func(interface{}) IAgent

type INetServer interface {
	core.Fork
	SetMaxNum(int)
	Start(int) bool
	Listen(NewAgent)
}

type TCPServer struct {
	conn_size int
	handle    NewAgent
	_m        sync.Mutex
	ln        net.Listener
	wgNet     base.WaitGroup
	wgConn    base.WaitGroup
	conns     map[net.Conn]int
}

func (this *TCPServer) SetMaxNum(num int) {
	this.conn_size = num
}

func (this *TCPServer) Start(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Open Err:", err)
		return false
	}
	fmt.Println("Open Ok:", port)
	this.SetMaxNum(check_conn_size(this.conn_size))
	this.ln = ln
	return true
}

func (this *TCPServer) Listen(handle NewAgent) {
	this.handle = handle
}

func (this *TCPServer) Wrap(f func()) {
	this.wgNet.Wrap(f)
}

func (this *TCPServer) WrapChild(f func()) {
	this.wgConn.Wrap(f)
}

func (this *TCPServer) Run() {
	for {
		conn, err := this.ln.Accept()
		if err == nil {
			if this.checkConnect(conn) {
				this.handleAgent(conn)
			} else {
				fmt.Println("To many conn ;max is ", this.connSize())
			}
		} else {
			if check_server_error(err) {
				break
			}
		}
	}
}

func (this *TCPServer) handleAgent(conn net.Conn) {
	tx := this.handle(conn)
	tx.OnReady(conn)
	this.WrapChild(func() {
		tx.Run()
		this.delConn(conn)
		tx.OnDie()
	})
}

func (this *TCPServer) OnDie() {

}

func (this *TCPServer) Wait() {
	this.wgNet.Wait()
	this.cleanConns()
	this.wgConn.Wait()
}

func (this *TCPServer) Die() bool {
	return this.ln.Close() == nil
}

//connSets
func (this *TCPServer) checkConnect(conn net.Conn) bool {
	this.lock()
	defer this.unlock()
	if this.conns == nil {
		this.conns = make(map[net.Conn]int)
	} else {
		if len(this.conns) >= this.conn_size {
			conn.Close()
			return false
		}
	}
	this.conns[conn] = 1
	return true
}

func (this *TCPServer) delConn(conn net.Conn) {
	this.lock()
	if this.conns != nil {
		delete(this.conns, conn)
	}
	this.unlock()
	//避免其他地方没关闭
	conn.Close()
}

func (this *TCPServer) cleanConns() {
	this.lock()
	if this.conns != nil {
		for conn := range this.conns {
			conn.Close()
		}
	}
	this.conns = nil
	this.unlock()
}

func (this *TCPServer) connSize() int {
	size := 0
	this.lock()
	if this.conns != nil {
		size = len(this.conns)
	}
	this.unlock()
	return size
}

func (this *TCPServer) lock() {
	this._m.Lock()
}

func (this *TCPServer) unlock() {
	this._m.Unlock()
}

//查看端口是否错误，错误就关闭服务
func check_server_error(err error) bool {
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
		return true
	}
	return false
}

//启动一个服务
func RunAndThrowServer(ser INetServer, port int, handle NewAgent, args ...func()) (INetServer, bool) {
	if !ser.Start(port) {
		return ser, false
	}
	ser.Listen(handle)
	//自身进程
	ser.Wrap(ser.Run)
	//守护进程
	core.WaitDaemon(ser, args...)
	return ser, true
}
