package gnet

import "net"
import "ants/gsys"

//关闭命令
type workChan chan []interface{}

//公开它
type Context interface {
	//直接关闭
	Close() error
	//关闭之前告白
	CloseOf(...interface{})
	//异步写
	Send(...interface{}) bool
	//异步读(内部知道就好了)
	Recv([]byte) ([]interface{}, error)
	//协议处理
	SetProcesser(IProtocoler)
	//等待结束
	Join(func([]byte))
}

//网络环境实例
type NetConn struct {
	L        gsys.Locked
	Conn     net.Conn
	coder    IProtocoler
	work     workChan
	openFlag bool
}

func NewConn(conn interface{}) Context {
	this := new(NetConn)
	this.SetConn(conn)
	return this
}

//protected 继承能获得
func (this *NetConn) SetConn(conn interface{}) {
	this.openFlag = true
	this.Conn = conn.(net.Conn)
	this.coder = NewProtocoler()
	this.work = make(workChan, NET_CHAN_SIZE)
}

func (this *NetConn) CloseWrite() {
	this.L.Lock()
	if this.openFlag {
		this.openFlag = false
		close(this.work)
	}
	this.L.Unlock()
}

func (this *NetConn) CloseRead() error {
	return this.Conn.Close()
}

func (this *NetConn) Write(bits []byte) (int, error) {
	return this.Conn.Write(bits)
}

func (this *NetConn) Read(b []byte) (int, error) {
	return this.Conn.Read(b)
}

// interface INetConn
func (this *NetConn) Close() error {
	this.CloseWrite()
	return this.CloseRead()
}

func (this *NetConn) CloseOf(args ...interface{}) {
	this.Write(this.coder.Marshal(args...))
	this.CloseWrite()
}

func (this *NetConn) Send(args ...interface{}) bool {
	this.L.Lock()
	if this.openFlag {
		this.work <- args
	}
	this.L.Unlock()
	return true
}

func (this *NetConn) Recv(b []byte) ([]interface{}, error) {
	ret, err := this.Read(b)
	if err == nil {
		return this.coder.Unmarshal(b[:ret]), nil
	}
	return nil, err
}

func (this *NetConn) SetProcesser(val IProtocoler) {
	this.coder = val
}

func (this *NetConn) Join(block func([]byte)) {
	buffSize := this.coder.BuffSize()
	wg := gsys.NewGroup()
	//异步read
	wg.Wrap(func() {
		b := make([]byte, buffSize)
		for {
			if ls, err := this.Recv(b); err == nil {
				for i := range ls {
					block(ToBytes(ls[i]))
				}
			} else {
				break
			}
		}
		this.CloseWrite()
	})
	//异步send
	wg.Wrap(func() {
		for {
			args, ok := <-this.work
			if ok {
				this.Write(this.coder.Marshal(args...))
			} else {
				break
			}
		}
		this.CloseRead()
	})
	//等待结束
	wg.Wait()
}
