package gnet

import "fmt"
import "net"
import "sync"
import "ants/gsys"

const (
	NET_OK = 0
	NET_NO = 1
)

//网络接口
type IConn interface {
	Close() bool            //直接关闭
	WriteBytes([]byte) bool //无协议写入
	//异步 (写入缓冲)
	Send(interface{}) bool
	CloseRead() bool
	CloseWrite() bool
}

type IConnDelegate interface {
	OnRead([]byte)
	OnSend([]byte)
	OnClose()
	BuffSize() int
}

//网络环境实例
type NetConn struct {
	gsys.Locked
	Conn      net.Conn
	sendFlag  bool
	readFlag  bool
	closeFlag bool
	buff      chan []byte
	wgConn    sync.WaitGroup
}

func (this *NetConn) SetConn(conn net.Conn) {
	this.Conn = conn
	this.buff = make(chan []byte, 1024)
}

func (this *NetConn) CloseWrite() bool {
	this.Lock()
	defer this.Unlock()
	if !this.sendFlag {
		this.sendFlag = true
		this.buff <- nil
		return true
	}
	return false
}

func (this *NetConn) CloseRead() bool {
	this.Lock()
	defer this.Unlock()
	if !this.readFlag {
		this.readFlag = true
		return true
	}
	return false
}

// interface INetConn
func (this *NetConn) Close() bool {
	if this.doClose() {
		this.OnFinally()
		return true
	}
	return false
}

func (this *NetConn) doClose() bool {
	this.Lock()
	defer this.Unlock()
	if this.closeFlag {
		return false
	}
	this.closeFlag = true
	return true
}

func (this *NetConn) OnFinally() {
	this.CloseRead()
	this.CloseWrite()
	if err := this.Conn.Close(); err != nil {
		fmt.Println("Close Err:", err)
	}
}

func (this *NetConn) Send(data interface{}) bool {
	this.Lock()
	defer this.Unlock()
	if this.sendFlag {
		fmt.Println("close Flag = ", this.closeFlag)
		return false
	}
	this.buff <- ToBytes(data)
	return true
}

func (this *NetConn) WriteBytes(bits []byte) bool {
	if ret, err := this.Conn.Write(bits); err != nil {
		fmt.Println("Write Err:", err, ", ret=", ret)
		return false
	}
	return true
}

func (this *NetConn) ReadDelegate(delegate IConnDelegate) {
	this.wgConn.Add(1)
	go func() {
		defer this.wgConn.Done()
		defer this.CloseWrite()
		bits := make([]byte, delegate.BuffSize())
		for {
			if ret, err := this.Conn.Read(bits); err == nil {
				delegate.OnRead(bits[:ret])
			} else {
				break
			}
		}
		this.Close()
	}()
	//异步写
	this.wgConn.Add(1)
	go func(buff chan []byte) {
		defer this.wgConn.Done()
		defer this.CloseRead()
		defer close(buff)
		for {
			if v, ok := <-buff; ok {
				if v == nil {
					continue
				}
				this.WriteBytes(v)
				this.OnSend(v)
			} else {
				break
			}
			if this.sendFlag {
				break
			}
		}
	}(this.buff)
	//等着结束
	this.wgConn.Wait()
	delegate.OnClose()
}

//通用接口
func (this *NetConn) OnRead(b []byte) {

}

func (this *NetConn) OnSend(b []byte) {

}

func (this *NetConn) OnClose() {

}

func (this *NetConn) BuffSize() int {
	return 1024
}
