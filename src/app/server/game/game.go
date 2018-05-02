package game

import "app/command"

//非网关模块通用(状态服务器)
import "ants/gnet"
import "ants/core"
import "ants/gcode"
import "ants/glog"
import "app/conf"

//game scene
import _ "app/server/game/texas"
import _ "app/server/game/taurus"

//服务器的启动(快速启动)
func ServerLaunch(port int) {
	//模块调度
	ref := new(LogicActor)
	if conf.LOCAL_TEST {
		core.Main().Join(conf.TOPIC_GAME, ref)
	} else {
		core.RunAndThrowBox(ref, nil, func() {
			//重启
		})
		run_service(port, ref)
	}
}

func run_service(port int, ref core.IBox) {
	gnet.RunAndThrowServer(new(gnet.TCPServer), port, func(conn interface{}) gnet.IAgent {
		session := gnet.NewProxy(conn)
		session.SetReceiver(func(b []byte) {
			ref.Push(gnet.NewBytes(session, b))
		})
		return session
	}, func() {
		ref.Die()
	})
}

//逻辑模块
type LogicActor struct {
	core.BaseBox
	//房间信息管理
	tables map[int]interface{}
	//用户管理
	users map[int]*DataPlayer
}

func (this *LogicActor) OnReady() {
	this.SetName("游戏服务器")
	this.SetAgent(this)
	this.SetBlock(this.OnMessage)
	this.tables = make(map[int]interface{})
	this.users = make(map[int]*DataPlayer)
	//this.ActorOf(100, NewTexasLogic(100, nil))
	//this.ActorOf(101, taurus.NewGame(101, nil))
}

func (this *LogicActor) OnMessage(event interface{}) {
	data := event.(*gnet.SocketEvent)
	pack := data.BeginPack()
	switch pack.Cmd() {
	case command.SERVER_ADD_PLAYER: //世界传来
		this.on_add_player(pack)
	case command.SERVER_DEL_PLAYER: //世界传来
		this.on_del_player(pack)
	case command.SERVER_OPEN_TABLE: //其他地方传来
		//开放一个房间游戏
	case command.SERVER_STOP_TABLE: //其他地方传来
		//停止一个房间游戏
	case command.CLIENT_GAME_ENTER: //用户传来
		this.on_enter(pack)
	case command.CLIENT_GAME_DROPS: //房间传来
		this.on_exit(pack)
	default: //用户传来
		this.on_table_notice(pack)
	}
}

//message
func (this *LogicActor) on_add_player(pack gcode.ISocketPacket) {
	header := NewHeader(pack)
	if oplayer, ok := this.users[header.UserID]; ok {
		//刷新
		oplayer.GateID = header.GateID
		oplayer.SessionID = header.SessionID
		oplayer.Status = STATE_ONLINE
		glog.Debug("game reconnect uid=%d gate=%d session=%v", header.UserID, header.GateID, header.SessionID)
		this.reconnect_table(oplayer)
	} else {
		player := new(DataPlayer)
		player.init()
		player.UserID = header.UserID
		player.GateID = header.GateID
		player.SessionID = header.SessionID
		this.users[header.UserID] = player
		glog.Debug("game connect uid=%d gate=%d session=%v", header.UserID, header.GateID, header.SessionID)
	}
}

//重连
func (this *LogicActor) reconnect_table(player *DataPlayer) {
	list := player.HomeList()
	for i := range list {
		psend := gcode.NewPackArgs(command.CLIENT_GAME_RECONNENT, player.UserID, player.GateID, player.SessionID)
		this.Send(list[i], psend)
	}
}

func (this *LogicActor) on_del_player(pack gcode.ISocketPacket) {
	uid := int(pack.ReadInt())
	if player, ok := this.users[uid]; ok {
		player.Status = STATE_DROPS
		this.drops_player(player)
		this.check_player_out(player)
	} else {
		println("del err: no user ", uid)
	}
}

//
func (this *LogicActor) drops_player(player *DataPlayer) {
	list := player.HomeList()
	for i := range list {
		this.Send(list[i], gcode.NewPackArgs(command.CLIENT_GAME_DROPS, player.UserID))
	}
}

//room message
func (this *LogicActor) on_enter(pack gcode.ISocketPacket) {
	header := NewHeader(pack)
	pack.ReadValue(&header.TableID)
	if player, ok := this.users[header.UserID]; ok {
		if player.Enter(header.TableID) {
			//发送失败表示无房间
			if this.Send(header.TableID, gcode.NewPackArgs(pack.Cmd(), header.UserID, header.GateID, header.SessionID)) {
				println("enter ok :", header.TableID)
			} else {
				//获取用户数据
				player.Exit(header.TableID)
				println("enter err, not table:", header.TableID)
			}
		} else {
			println("enter err, state is lock:", header.UserID)
		}
	} else {
		println("enter err not uid:", header.UserID)
	}
}

//出来需要>设置状态(不用匹配)
func (this *LogicActor) on_exit(pack gcode.ISocketPacket) {
	uid, roomid := int(pack.ReadInt()), int(pack.ReadInt())
	if player, ok := this.users[uid]; ok {
		player.Exit(roomid)
		println("exit room ok:", uid, roomid)
		this.check_player_out(player)
	}
}

//交给其他的处理(不检测玩家是否存在)
func (this *LogicActor) on_table_notice(pack gcode.ISocketPacket) {
	header := NewHeader(pack)
	pack.ReadValue(&header.TableID)
	body := pack.ReadRemaining()
	//直接推送(必须进入房间过，不然不会有结果)
	this.Send(header.TableID, gcode.NewPackArgs(pack.Cmd(), header.UserID, body))
}

//是否离线且没有在房间
func (this *LogicActor) check_player_out(player *DataPlayer) {
	if player.Die() && player.Empty() {
		println("退出ALL房间,退出游戏:", player.UID())
		delete(this.users, player.UID())
	}
}

//房间退出最后的通知(处理一些数据问题)
func (this *LogicActor) on_remove_table(pack gcode.ISocketPacket) {

}
