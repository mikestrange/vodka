package gnet

//服务器基本接口
type ServerBlock func(INetContext) INetProxy

type INetServer interface {
	Start()
	Close()
	ConnSize() int
}

//网络代理接口
type INetProxy interface {
	Run()
	OnClose(int)
}

//网络接口
type INetConn interface {
	Shutdown(int) int            //直接关闭
	WriteBytes([]byte) int       //无协议写入
	ReadBytes(int, func([]byte)) //无协议读取
}

//网络环境接口 :双通道读写(内部使用)
type ContextBlock func(int, []byte)

type INetContext interface {
	INetProxy
	INetConn //(隐式)
	//关闭(异步些)
	Close() int
	CloseSign(int) int
	//是否客户端链接
	AsSocket() bool
	//回调机制
	SetHandle(ContextBlock)
	//默认
	SetProcessor(INetProcessor)
	//发送多个包(直接发送)
	Send(...interface{})
}

//网络通道
type INetChan interface {
	Push([]byte) bool
	Pull() ([]byte, bool)
	Loop(func([]byte))
	AsynClose()
	Close()
}

//解码-编码
type INetProcessor interface {
	//解码直接获得消息
	Unmarshal([]byte) []interface{}
	//Message() []interface{}
	//编码
	Marshal(...interface{}) []byte
	//Commit() []byte
}

//获得字节接口
type IBytes interface {
	Bytes() []byte
}

//快速启动服务器
func ListenAndRunServer(port int, block func(IBaseProxy)) INetServer {
	ser := NewTcpServer(port, func(tx INetContext) INetProxy {
		session := NewBaseProxy(tx)
		block(session)
		return session
	})
	ser.Start()
	return ser
}
