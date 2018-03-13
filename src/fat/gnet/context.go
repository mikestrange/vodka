package gnet

//不存在太多的异步操作，目前不考虑枷锁(阻塞)
//注意:(目前可能存在丢消息的情况,当关闭的时候,可能存在写消息的情况)
import (
	"fat/gsys"
	"fmt"
	"io"
	"net"
)

//接受回调
type ContextBlock func(INetContext, interface{})

//二进制接口数据
type IBytes interface {
	Bytes() []byte
}

//网络代理(只需要知道关闭的原因)
type INetProxy interface {
	Conn() net.Conn
	Send(interface{}) bool
	Read() interface{}
	Close() int
	CloseSign(int) int
}

//网络环境(断开后不要重复使用)
type INetContext interface {
	INetProxy
	//是否为客户端链接
	IsDie() bool
	AsSocket() bool
	Ping() bool
	LivePing()
	//i/o
	SendBytes([]byte) int
	ReadBytes(int, func([]byte)) int
	//可选
	SetClient(interface{}) interface{}
	Client() interface{}
}

//网络处理(目前只接收处理)
type ISocketHandler interface {
	BuffSize() int
	LoadBytes([]byte)
	Pack() (interface{}, bool)
}

//网络环境实例
type NetContext struct {
	gsys.Locked
	mtype  int
	live   bool
	heart  bool
	used   int
	conn   net.Conn
	client interface{}
}

//服务器链接
func NewConn(conn interface{}) INetContext {
	return NewContextWithType(conn, NET_TYPE_CONN)
}

//客户端链接
func NewSocket(addr string) (INetContext, bool) {
	conn, ok := DialConn(addr)
	if ok {
		return NewContextWithType(conn, NET_TYPE_SOCKET), true
	}
	return nil, false
}

//根据类型分配
func NewContextWithType(conn interface{}, mtype int) INetContext {
	this := new(NetContext)
	this.InitContext(conn, mtype)
	return this
}

//private begin
func (this *NetContext) InitContext(conn interface{}, mtype int) {
	this.InitLocked()
	this.SetConn(conn, mtype)
}

func (this *NetContext) SetConn(conn interface{}, mtype int) {
	this.used = 0
	this.heart = true
	this.live = true
	this.mtype = mtype
	this.conn = conn.(net.Conn)
}

//private end

func (this *NetContext) Conn() net.Conn {
	return this.conn
}

func (this *NetContext) Close() int {
	return this.close_asyn_sign(SIGN_CLOSE_OK)
}

func (this *NetContext) CloseSign(used int) int {
	return this.close_asyn_sign(used)
}

func (this *NetContext) close_asyn_sign(used int) int {
	this.Lock()
	if this.live {
		this.used = used
		this.live = false
		this.conn.Close()
	}
	this.Unlock()
	return this.used
}

func (this *NetContext) Send(data interface{}) bool {
	return this.SendBytes(ToBytes(data)) >= 0
}

func (this *NetContext) Read() interface{} {
	return nil
}

//public INetContex begin
func (this *NetContext) SendBytes(bits []byte) int {
	if bits == nil || len(bits) == 0 {
		return -1 //无效字节
	}
	if ret, err := this.conn.Write(bits); err == nil {
		return ret
	}
	return -2
}

func (this *NetContext) ReadBytes(size int, block func([]byte)) int {
	defer this.check_close()
	bits := make([]byte, size)
	code := SIGN_CLOSE_OK
	for {
		ret, err := this.conn.Read(bits)
		if err == nil {
			//1,如果在这里报错，那么会到defer不会向下走
			//2,如果里面处理了报错，那么就会往下面走
			block(bits[:ret])
		} else {
			if err == io.EOF {
				code = SIGN_CLOSE_DISTAL
			} else {
				code = SIGN_CLOSE_SELF
			}
			break
		}
	}
	//报错会走上面defer
	this.CloseSign(code)
	return code
}

func (this *NetContext) check_close() {
	if err := recover(); err != nil {
		fmt.Println("Conn Throw Err:", err)
		this.CloseSign(SIGN_CLOSE_ERROR)
	}
}

//sets
func (this *NetContext) SetClient(val interface{}) interface{} {
	this.Lock()
	old := this.client
	this.client = val
	this.Unlock()
	return old
}

func (this *NetContext) Client() interface{} {
	this.Lock()
	client := this.client
	this.Unlock()
	return client
}

func (this *NetContext) LivePing() {
	this.heart = true
}

func (this *NetContext) Ping() bool {
	if this.heart {
		this.heart = false
		return true
	}
	return false
}

func (this *NetContext) AsSocket() bool {
	return this.mtype == NET_TYPE_SOCKET
}

func (this *NetContext) IsDie() bool {
	return this.live == false
}

//static
//方式1:无任何处理
func LoopContext(tx INetContext, block func([]byte)) {
	tx.ReadBytes(BUFFER_SIZE, block)
}

//方式2:包处理
func LoopWithHandle(tx INetContext, block ContextBlock) {
	handle := NewSocketHandler()
	tx.ReadBytes(handle.BuffSize(), func(bits []byte) {
		handle.LoadBytes(bits)
		for {
			if data, ok := handle.Pack(); ok {
				block(tx, data)
			} else {
				break
			}
		}
	})
}

//方式3:包处理+心跳处理
func LoopConnWithPing(tx INetContext, block ContextBlock) {
	timer := gsys.Forever(PING_TIME, func() {
		if tx.Ping() {
			tx.Send(PacketWithHeartBeat)
		} else {
			tx.CloseSign(SIGN_CLOSE_PING)
		}
	})
	defer timer.Stop()
	LoopWithHandle(tx, block)
}
