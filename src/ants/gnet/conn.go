package gnet

import "fmt"
import "net"
import "ants/gsys"

const (
	NET_OK = 0
	NET_NO = 1
)

//关闭命令
type ConnBlock func([]byte)
type sendChan chan []interface{}

//网络接口
type IConn interface {
	Close() error
	//最后的通话
	CloseOf(...interface{})
	//关闭读写
	CloseRead()
	CloseWrite()
	//读写
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	//异步 (写入缓冲)
	Send(...interface{}) bool
	Recv(interface{}) bool
	//轮询写入
	LoopWrite(int, func(interface{}))
	LoopRead(int, func(interface{}))
	//处理
	SetCoder(IProtocoler)
	SetHandle(ConnBlock)
	//运行等待结束
	WaitFor()
}

//网络环境实例
type NetConn struct {
	gsys.Locked
	openFlag bool
	Conn     net.Conn
	coder    IProtocoler
	handle   ConnBlock
	chans    sendChan
}

func NewConn(conn interface{}) IConn {
	this := new(NetConn)
	this.SetConn(conn)
	return this
}

func (this *NetConn) SetConn(conn interface{}) {
	this.openFlag = true
	this.Conn = conn.(net.Conn)
	this.coder = NewProtocoler()
	this.chans = make(sendChan, NET_CHAN_SIZE)
}

func (this *NetConn) CloseWrite() {
	this.Lock()
	if this.openFlag {
		this.openFlag = false
		close(this.chans)
	}
	this.Unlock()

}

func (this *NetConn) CloseRead() {
	if err := this.Conn.Close(); err != nil {
		fmt.Println("Close Err:", err)
	}
}

// interface INetConn
func (this *NetConn) Close() error {
	this.CloseRead()
	this.CloseWrite()
	return nil
}

func (this *NetConn) CloseOf(args ...interface{}) {
	this.Write(this.getCoder().Marshal(args...))
	this.CloseWrite()
}

func (this *NetConn) Send(args ...interface{}) bool {
	this.Lock()
	if this.openFlag {
		this.chans <- args
	}
	this.Unlock()
	return true
}

func (this *NetConn) Recv(data interface{}) bool {
	//println("recv not used")
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
	for {
		args, ok := <-this.chans
		if ok {
			block(args)
		} else {
			break
		}
	}
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

func (this *NetConn) SetCoder(val IProtocoler) {
	this.coder = val
}

func (this *NetConn) SetHandle(handle ConnBlock) {
	this.handle = handle
}

//这里才是运行的根本
func (this *NetConn) WaitFor() {
	pack := this.getCoder()
	buffSize := pack.BuffSize()
	wg := gsys.NewWgGroup()
	//异步读消息
	wg.Add()
	go func() {
		this.LoopRead(buffSize, func(data interface{}) {
			list := pack.Unmarshal(ToBytes(data))
			for i := range list {
				this.doMsgr(ToBytes(list[i]))
			}
		})
		this.CloseWrite()
		wg.Done()
	}()
	//异步写消息
	wg.Add()
	go func() {
		this.LoopWrite(buffSize, func(args interface{}) {
			this.Write(this.getCoder().Marshal(args.([]interface{})...))
		})
		this.CloseRead()
		wg.Done()
	}()
	//等着结束
	wg.Wait()
}

func (this *NetConn) getCoder() IProtocoler {
	return this.coder
}

func (this *NetConn) doMsgr(b []byte) {
	if this.handle != nil {
		this.handle(b)
	} else {
		fmt.Println("no handle conn")
	}
}

func (this *NetConn) check_error() {
	if err := recover(); err != nil {
		fmt.Println("Conn Err:", err)
	}
}
