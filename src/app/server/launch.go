package server

import "ants/conf"

import "app/server/gate"
import "app/server/logon"
import "app/server/world"
import "app/server/chat"
import "app/server/game"
import "app/server/hall"

//
import "ants/actor"
import "ants/cluster"

func init() {
	//路由分配
	conf.EachVo(func(vo *conf.RouteVo) {
		actor.Main.ActorOf(vo.Topic, cluster.NewRouterPort(vo.Port))
	})
}

func Launch(name string) {
	switch name {
	case "gate":
		go gate.ServerLaunch(conf.PORT_GATE, conf.TOPIC_GATE)
	case "login":
		go logon.ServerLaunch(conf.PORT_LOGIN)
	case "game":
		go game.ServerLaunch(conf.PORT_GAME)
	case "world":
		go world.ServerLaunch(conf.PORT_WORLD)
	case "chat":
		go chat.ServerLaunch(conf.PORT_CHAT)
	case "hall":
		go hall.ServerLaunch(conf.PORT_HALL)
	default:
		Launch("gate")
		Launch("login")
		Launch("world")
		Launch("chat")
		Launch("game")
		Launch("hall")
	}
}
