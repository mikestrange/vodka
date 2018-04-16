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
	//Read([]byte) (int, error)
	//Write([]byte) (int, error)
	//读消息
	ReadMsg(func([]byte))
	//循环写
	WriteLoop()
	//写消息
	WriteMsg(...interface{}) error
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

//编码
func (this *NetConn) Decode(args ...interface{}) ([]byte, error) {
	return this.coder.Marshal(args...)
}

//解码
func (this *NetConn) UnDecode(bits []byte) ([][]byte, error) {
	return this.coder.Unmarshal(bits)
}

//缓冲
func (this *NetConn) MakeBuffer() []byte {
	if this.coder == nil {
		return make([]byte, NET_BUFF_NEW_SIZE)
	}
	return make([]byte, this.coder.BuffSize())
}

//处理者(决定上三个)
func (this *NetConn) SetProcesser(val IProtocoler) {
	this.coder = val
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

func (this *NetConn) Write(b []byte) (int, error) {
	ret, err := this.Conn.Write(b)
	return ret, err
}

func (this *NetConn) Read(b []byte) (int, error) {
	ret, err := this.Conn.Read(b)
	return ret, err
}

func (this *NetConn) WriteMsg(args ...interface{}) error {
	ret, err := this.Decode(args...)
	if err == nil {
		this.Write(ret)
	}
	return err
}

// interface INetConn
func (this *NetConn) Close() error {
	this.CloseWrite()
	return this.CloseRead()
}

func (this *NetConn) CloseOf(args ...interface{}) {
	this.WriteMsg(args...)
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

func (this *NetConn) ReadMsg(block func([]byte)) {
	bit := this.MakeBuffer()
	//异步读取
	for {
		if ret, err := this.Read(bit); err == nil {
			if list, uerr := this.UnDecode(bit[:ret]); uerr == nil {
				for i := range list {
					//Client错误就关闭
					block(list[i])
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	this.CloseWrite()
}

func (this *NetConn) WriteLoop() {
	//写入消息(一般不存在错误处理)
	for args := range this.work {
		this.WriteMsg(args...)
	}
	this.CloseRead()
}

func (this *NetConn) Join(block func([]byte)) {
	gsys.Wraps(func() {
		this.ReadMsg(block)
	}, func() {
		this.WriteLoop()
	})
}
