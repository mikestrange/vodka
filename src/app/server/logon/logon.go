package logon

import "fmt"
import "app/command"

//非网关模块通用
import "ants/conf"
import "ants/gnet"
import "ants/actor"

//弱连接服务器，不用管心跳
func ServerLaunch(port int) {
	//数据服务器链接
	init_dber()
	//模块调度
	refLogic := actor.RunAndThrowBox(new(LogicActor), nil)
	//服务器快速启动
	gnet.ListenAndRunServer(port, func(session gnet.IBaseProxy) {
		session.SetHandle(func(b []byte) {
			refLogic.Router(gnet.NewPackBytes(b))
		})
	})
}

//逻辑块
type LogicActor struct {
	actor.BaseBox
}

func (this *LogicActor) OnReady() {
	this.SetActor(this)
}

func (this *LogicActor) OnDie() {

}

func (this *LogicActor) PerformRunning() {
	this.Worker().ReadRound(this, 1000)
}

func (this *LogicActor) OnMessage(args ...interface{}) {
	pack := args[0].(gnet.ISocketPacket)
	switch pack.Cmd() {
	case command.CLIENT_LOGON:
		this.on_logon(pack)
	default:
		println("login no handle:", pack.Cmd())
	}
}

//message
func (this *LogicActor) on_logon(pack gnet.ISocketPacket) {
	//header
	UserID, PassWord, SerID, SessionID := pack.ReadInt(), pack.ReadString(), pack.ReadInt(), pack.ReadUInt64()
	//other
	fmt.Println(fmt.Sprintf("Logon Info# uid=%d, session=%v gateid=%d", UserID, SessionID, SerID))
	err_code := check_user(UserID, PassWord)
	fmt.Println("Seach Result Code:", err_code, UserID, PassWord, SerID, SessionID)
	var body []byte = []byte{}
	if err_code == 0 {
		body = get_user_info(UserID)
	}
	//错误直接返回
	if err_code != 0 {
		this.Main().Send(conf.TOPIC_WORLD, pack_logon_result(err_code, UserID, SerID, SessionID, body))
	} else {
		this.Main().Send(conf.TOPIC_WORLD, pack_logon_result(err_code, UserID, SerID, SessionID, body))
	}
}

//send world(通知登录结果)
func pack_logon_result(code int, uid int, gate int, session uint64, body []byte) gnet.IBytes {
	return gnet.NewPackArgs(command.SERVER_WORLD_ADD_PLAYER, int16(code), uid, gate, session, body)
}
