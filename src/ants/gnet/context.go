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
	csend  INetChan
	cread  INetChan
	handle ContextBlock
	decode INetProcessor
}

func (this *context) Init(conn net.Conn, size int, client bool) INetContext {
	this.closed = false
	this.conn = conn
	this.client = client
	this.used = SIGN_CLOSE_NULL
	this.decode = NewSocketProcessor()
	this.csend = newChan(size)
	//this.cread = newChan(size)
	this.thread()
	return this
}

func (this *context) thread() {
	//写通道
	go func(buff INetChan) {
		buff.Loop(func(bits []byte) {
			this.WriteBytes(bits)
		})
		this.Shutdown(SIGN_CLOSE_ERROR_SEND)
	}(this.csend)
	//	//读通道
	//	go func(buff INetChan) {
	//		buff.Loop(func(bits []byte) {
	//			t := gutil.GetNano()
	//			this.Done(EVENT_CONN_READ, bits)
	//			println("Msg ##:", gutil.NanoStr(gutil.GetNano()-t), this.AsSocket())
	//		})
	//		this.Shutdown(SIGN_CLOSE_ERROR_READ)
	//	}(this.cread)
}

//interface INetContext
func (this *context) SetHandle(val ContextBlock) {
	this.handle = val
}

func (this *context) SetProcessor(val INetProcessor) {
	this.decode = val
}

func (this *context) Close() int {
	return this.CloseSign(SIGN_CLOSE_OK)
}

func (this *context) CloseSign(used int) int {
	code := this.setUsed(used)
	this.csend.AsynClose()
	return code
}

func (this *context) AsSocket() bool {
	return this.client
}

func (this *context) Send(args ...interface{}) {
	this.csend.Push(this.decode.Marshal(args...))
}

// interface INetConn
func (this *context) Shutdown(used int) int {
	this.Lock()
	if !this.closed {
		this.closed = true
		this.conn.Close()
		//关闭读写
		this.csend.Close()
		//this.cread.Close()
	}
	this.Unlock()
	return this.setUsed(used)
}

func (this *context) WriteBytes(bits []byte) int {
	//	size := len(bits)
	//	buffSize := NET_BUFF_NEW_SIZE
	//	if buffSize > size {
	//		buffSize = size
	//	}
	//	pos := 0
	//	for {
	//		b := bits[pos : pos+buffSize]
	//		pos = pos + buffSize
	//		if sub := size - pos; sub < buffSize {
	//			buffSize = sub
	//		}
	//		if pos == size {
	//			break
	//		}
	//		this.conn.Write(b)
	//	}
	if ret, err := this.conn.Write(bits); err == nil {
		//println("write:", ret)
		return ret
	}
	return SIGN_SEND_ERROR
}

func (this *context) ReadBytes(size int, block func([]byte)) {
	if this.closed {
		println("this conn is closed")
		return
	}
	bits := make([]byte, size)
	code := SIGN_CLOSE_OK
	//	defer func() {
	//		if err := recover(); err != nil {
	//			this.Shutdown(SIGN_CLOSE_ERROR_READ)
	//		} else {
	//			this.Shutdown(code)
	//		}
	//	}()
	for {
		ret, err := this.conn.Read(bits)
		if err == nil {
			//println("read:", ret)
			block(bits[:ret])
		} else {
			if err == io.EOF {
				code = SIGN_CLOSE_DISTAL
			} else {
				code = SIGN_CLOSE_OK
			}
			break
		}
	}
	this.Shutdown(code)
}

//interface INetProxy
func (this *context) Run() {
	this.ReadBytes(NET_BUFF_NEW_SIZE, func(bits []byte) {
		list := this.decode.Unmarshal(bits)
		for i := range list {
			//this.AsynRead(list[i])
			this.Done(EVENT_CONN_READ, ToBytes(list[i]))
		}
	})
}

func (this *context) OnClose(code int) {
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

func (this *context) setUsed(used int) int {
	this.Lock()
	if this.used == SIGN_CLOSE_NULL {
		this.used = used
	}
	this.Unlock()
	return this.used
}

func (this *context) AsynSend(data interface{}) {
	if ok := this.csend.Push(ToBytes(data)); ok {
		return
	}
	this.Shutdown(SIGN_CLOSE_ERROR_SEND)
}

//目前没使用
func (this *context) AsynRead(data interface{}) {
	if ok := this.cread.Push(ToBytes(data)); ok {
		return
	}
	this.Shutdown(SIGN_CLOSE_ERROR_READ)
}
