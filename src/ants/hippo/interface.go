package hippo

const (
	//关闭信号
	CLOSE_SIGN_SELF   = 0 //自己关闭
	CLOSE_SIGN_ERROR  = 1 //内部错误
	CLOSE_SIGN_CLIENT = 2 //客户端关闭
	//事务处理
	EVENT_SEND    = 1
	EVENT_READ    = 2
	EVENT_CLOSED  = 3
	EVENT_TIMEOUT = 4
	//超时类型
	CODE_TIMEOUT_READ = 1
	CODE_TIMEOUT_SEND = 2
)

//网络工作
type IActor interface {
	Make(int) bool
	SetTimeout(int)            //读写超时
	SetCaller(interface{})     //设置调度者
	SetHandle(IHandle)         //设置处理器
	PushSend(interface{}) bool //推送写入
	SetSendTimeout(int)        //写超时(超时发事件)
	PushRead(interface{}) bool //推送读取
	SetReadTimeout(int)        //读超时
	PushEvent(IEvent) bool     //推送其他事务
	Exit(int) bool             //推送关闭
	LoopWithTimeout()          //超时包含
	LoopWithNothing()          //不包含超时
}

//处理
type IHandle interface {
	Handle(int, interface{}, interface{})
}

//协议
type ICoding interface {
	Decode([]byte) ([]interface{}, error)  //解码(发生错误直接抛出)
	Encode(...interface{}) ([]byte, error) //编码
	NewBits() []byte
}

//具体网络
type IConn interface {
	ReadBytes([]byte) (int, error)
	SendBytes([]byte) bool
	Local() string
	Remote() string
	Close() bool
}

//具体环境
type IContext interface {
	SetConn(IConn)                         //操作的
	Conn() IConn                           //链接
	SetCoding(ICoding)                     //编码格式
	ReadMsg([]byte) ([]interface{}, error) //读消息
	SendMsg(...interface{}) bool           //写消息
	CloseOf(...interface{})                //关闭消息
	CloseSign(int)                         //关闭信息
	Close()                                //直接关闭
	Loop()                                 //轮询
	SetActor(IActor)                       //工作
}

//套接字
type ISocket interface {
	SetAutoReconnect(bool) //是否自动重连
	Connect(string)        //连接(默认所有)
	SetLinkNum(int)        //默认一个
	Close() bool           //关闭所有
	Send([]byte) bool      //发送消息
}
