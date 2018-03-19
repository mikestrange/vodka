package gnet

//被封装的网络环境
import (
	"ants/gsys"
	"fmt"
	"io"
	"net"
)

func Context(conn net.Conn, size int) INetContext {
	return new(context).Init(conn, size, false)
}

func Socket(addr string) (INetContext, bool) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", addr)
		return new(context).Init(conn, NET_CHAN_SIZE, true), true
	}
	fmt.Println("Socket Connect Err:", err)
	return nil, false
}

//网络环境实例
type context struct {
	gsys.Locked
	closed bool
	used   int
	client bool
	conn   net.Conn
	buffer gsys.IAsynDispatcher
	handle ContextBlock
	decode INetProcessor
}

func (this *context) Init(conn net.Conn, size int, client bool) INetContext {
	this.closed = false
	this.conn = conn
	this.client = client
	this.decode = NewSocketProcessor()
	this.buffer = gsys.NewChannelSize(size)
	this.thread()
	return this
}

func (this *context) thread() {
	//写通道
	go func(buff gsys.IAsynDispatcher) {
		buff.Loop(func(bits interface{}) {
			this.WriteBytes(ToBytes(bits))
		})
		this.Kill()
	}(this.buffer)
}

//interface INetContext
func (this *context) SetHandle(val ContextBlock) {
	this.handle = val
}

func (this *context) SetProcessor(val INetProcessor) {
	this.decode = val
}

func (this *context) CloseWrite() error {
	this.buffer.AsynClose()
	return nil
}

func (this *context) CloseRead() error {
	return nil
}

func (this *context) AsSocket() bool {
	return this.client
}

func (this *context) Send(args ...interface{}) {
	this.Flush(this.decode.Marshal(args...))
}

// interface INetConn
func (this *context) Kill() error {
	this.Lock()
	defer this.Unlock()
	if !this.closed {
		this.closed = true
		this.buffer.Close()
		return this.conn.Close()
	}
	return nil
}

func (this *context) WriteBytes(bits []byte) bool {
	if ret, err := this.conn.Write(bits); err != nil {
		fmt.Println("Write Err:", err, ret)
		return false
	}
	return true
}

func (this *context) ReadBytes(size int, block func([]byte)) int {
	bits := make([]byte, size)
	code := SIGN_CLOSE_OK
	for {
		ret, err := this.conn.Read(bits)
		if err == nil {
			block(bits[:ret])
		} else {
			fmt.Println("Close Ok:", err)
			if err == io.EOF {
				code = SIGN_CLOSE_DISTAL
			} else {
				code = SIGN_CLOSE_OK
			}
			break
		}
	}
	return code
}

func (this *context) check_error() {
	if err := recover(); err != nil {
		fmt.Println("Conn Err:", err)
	}
}

//interface INetProxy
func (this *context) Run() {
	//defer this.check_error()
	this.ReadBytes(NET_BUFF_NEW_SIZE, func(bits []byte) {
		list := this.decode.Unmarshal(bits)
		for i := range list {
			this.Done(EVENT_CONN_READ, ToBytes(list[i]))
		}
	})
}

func (this *context) OnClose() {
	//override public
}

//private closed
func (this *context) Done(event int, bits []byte) {
	if this.handle != nil {
		this.handle(event, bits)
	} else {
		fmt.Println("context not set handle")
	}
}

func (this *context) Flush(data interface{}) {
	if ok := this.buffer.Push(data); ok {
		return
	}
	this.Kill()
}
