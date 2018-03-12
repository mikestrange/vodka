package logon

//只用来获取用户的数据，并发处理

//import "soy/db/link"
//import "soy/db/stock"
import "fat/gnet"
import "fat/gutil"
import "fat/gnet/nsc"
import "app/command"
import "app/config"

//var redis stock.IRedis
var router nsc.IRemoteScheduler
var mode gutil.IModeAccessor

func init() {
	//init_redis()
	init_mysql()
	//基础设施(基本属于单线程)
	router, mode = command.SetRouter(config.LOGIN_PORT, events, nil)
}

func init_redis() {
	//stock.Redis().Connect("localhost", 6379, "")
	//redis = stock.Redis()
}

func init_mysql() {
	//client := link.Mysql()
	//client.LinkAddr("127.0.0.1:3306", "root", "123456", "user_info")
	//client.Connect(20)
}

//服务器的启动
func ServerLaunch(port int) {
	gnet.ListenAndRunServer(port, func(conn interface{}) {
		gnet.LoopConnWithPing(gnet.NewConn(conn), func(tx gnet.INetContext, data interface{}) {
			mode.Done(data.(gnet.ISocketPacket).Cmd(), tx, data)
		})
	})
}
