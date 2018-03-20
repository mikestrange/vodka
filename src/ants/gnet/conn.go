package gnet

import "fmt"
import "net"
import "sync"
import "ants/gsys"

const (
	NET_OK = 0
	NET_NO = 1
)

type ConnBlock func([]byte)

//网络接口
type INetContext interface {
	Close() error
	//关闭读写
	CloseRead()
	CloseWrite()
	//读写
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	//异步 (写入缓冲)
	Send(interface{}) bool
	Recv(interface{}) bool
	//轮询写入
	LoopWrite(int, func(interface{}))
	LoopRead(int, func(interface{}))
	//处理
	SetHandle(ConnBlock)
	SetCoder(IProtocoler)
	//运行等待结束
	WaitFor()
}

//网络环境实例
type NetConn struct {
	gsys.Locked
	Conn      net.Conn
	handle    ConnBlock
	coder     IProtocoler
	buff      gsys.IAsynDispatcher
	closeFlag bool
}

func NewConn(conn interface{}) INetContext {
	this := new(NetConn)
	this.SetConn(conn)
	return this
}

func (this *NetConn) SetConn(conn interface{}) {
	this.Conn = conn.(net.Conn)
	this.buff = gsys.NewChannelSize(NET_CHAN_SIZE)
	this.coder = NewProtocoler()
}

func (this *NetConn) CloseWrite() {
	this.buff.AsynClose()
}

func (this *NetConn) CloseRead() {
	this.finally()
}

// interface INetConn
func (this *NetConn) Close() error {
	this.buff.Close() //直接关闭写入
	this.CloseRead()
	this.CloseWrite()
	//this.finally()
	return nil
}

func (this *NetConn) finally() {
	this.Lock()
	if !this.closeFlag {
		this.closeFlag = true
		if err := this.Conn.Close(); err != nil {
			fmt.Println("Close Err:", err)
		}
	}
	this.Unlock()
}

func (this *NetConn) Send(data interface{}) bool {
	if this.buff.Push(this.getCoder().Marshal(data)) { //写满了就关闭
		return true
	}
	this.CloseRead()
	return false
}

func (this *NetConn) Recv(data interface{}) bool {
	println("recv not used")
	return false
}

func (this *NetConn) Write(bits []byte) (int, error) {
	return this.Conn.Write(bits)
}

func (this *NetConn) Read(b []byte) (int, error) {
	return this.Conn.Read(b)
}

//阻塞写
func (this *NetConn) LoopWrite(size int, block func(interface{})) {
	this.buff.Loop(block)
}

//阻塞读
func (this *NetConn) LoopRead(size int, block func(interface{})) {
	bits := make([]byte, size)
	for {
		if ret, err := this.Read(bits); err == nil {
			block(bits[:ret])
		} else {
			break
		}
	}
}

func (this *NetConn) SetHandle(handle ConnBlock) {
	this.handle = handle
}

func (this *NetConn) SetCoder(val IProtocoler) {
	this.coder = val
}

//这里才是运行的根本
func (this *NetConn) WaitFor() {
	pack := this.getCoder()
	buffSize := pack.BuffSize()
	wgConn := new(sync.WaitGroup)
	wgConn.Add(1)
	go func() {
		defer wgConn.Done()
		defer this.CloseWrite()
		this.LoopRead(buffSize, func(data interface{}) {
			list := pack.Unmarshal(ToBytes(data))
			for i := range list {
				this.doMsgr(ToBytes(list[i]))
			}
		})
	}()
	//异步写
	wgConn.Add(1)
	go func() {
		defer wgConn.Done()
		defer this.CloseRead()
		this.LoopWrite(buffSize, func(data interface{}) {
			this.Write(ToBytes(data))
		})
	}()
	//等着结束
	wgConn.Wait()
	this.Close()
}

//释放消息
func (this *NetConn) doMsgr(b []byte) {
	if this.handle != nil {
		this.handle(b)
	} else {
		println("no handle")
	}
}

func (this *NetConn) getCoder() IProtocoler {
	return this.coder
}

func (this *NetConn) check_error() {
	if err := recover(); err != nil {
		fmt.Println("Conn Err:", err)
	}
}
