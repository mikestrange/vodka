package gnet

//网络的具体事务
const (
	EVENT_CONN_CONNECT   = 0
	EVENT_CONN_READ      = 1
	EVENT_CONN_SEND      = 2
	EVENT_CONN_CLOSE     = 3
	EVENT_CONN_SIGN      = 4 //信号
	EVENT_CONN_HEARTBEAT = 5 //心跳
)

//其他
const (
	//服务器容量
	NET_SERVER_SIZE = 10000
	//默认网络通道容量
	NET_CHAN_SIZE = 1000
	//默认心跳时间
	PING_TIME = 1000 * 60 * 5 //5分钟
)

//关闭信号
const (
	CLOSE_SIGN_CLIENT = 1 //客户端关闭
	CLOSE_SIGN_ERROR  = 2 //发生错误
	CLOSE_SIGN_SELF   = 3 //自己关闭
)

//发送通道大小
func check_send_size(size int) int {
	if size < NET_CHAN_SIZE {
		return NET_CHAN_SIZE
	}
	return size
}

func check_conn_size(size int) int {
	if size < NET_SERVER_SIZE {
		return NET_SERVER_SIZE
	}
	return size
}

//网络代理
type IAgent interface {
	Run()   //阻塞
	Wait()  //异步等待
	OnDie() //结束
}

//tcp链接
func Socket(addr string) (Context, bool) {
	this := new(TcpSocket)
	return this, this.Connect(addr)
}

//抛出代理运行
func RunAndThrowAgent(val IAgent, args ...func()) {
	go func() {
		val.Run()
		val.OnDie()
		for i := range args {
			args[i]()
		}
	}()
}
