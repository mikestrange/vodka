package world

//用户管理(每一段时间(10分钟)需要校正网关的登录情况，防止僵尸玩家)
import "fat/gutil"
import "fat/gnet/nsc"
import "fat/gnet"
import "app/command"
import "app/config"
import "fat/gsys"

var router nsc.IRemoteScheduler
var mode gutil.IModeAccessor

func init() {
	//基础设施(单线程)
	router, mode = command.SetRouter(config.WORLD_PORT, events, gsys.NewChannel())
}

//服务器的启动
func ServerLaunch(port int) {
	gnet.ListenAndRunServer(port, func(conn interface{}) {
		gnet.LoopConnWithPing(gnet.NewConn(conn), func(tx gnet.INetContext, data interface{}) {
			mode.Done(data.(gnet.ISocketPacket).Cmd(), tx, data)
		})
	})
}
