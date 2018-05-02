package server

import "app/conf"

import "app/server/gate"
import "app/server/logon"
import "app/server/world"
import "app/server/hall"
import "app/server/chat"
import "app/server/game"

func init() {
	//Launch("")
}

func Launch(name string) {
	switch name {
	case "gate":
		gate.ServerLaunch(conf.PORT_GATE, conf.TOPIC_GATE)
	case "login":
		logon.ServerLaunch(conf.PORT_LOGIN)
	case "game":
		game.ServerLaunch(conf.PORT_GAME)
	case "world":
		world.ServerLaunch(conf.PORT_WORLD)
	case "chat":
		chat.ServerLaunch(conf.PORT_CHAT)
	case "hall":
		hall.ServerLaunch(conf.PORT_HALL)
	default:
		Launch("gate")
		Launch("login")
		Launch("world")
		Launch("chat")
		Launch("game")
		Launch("hall")
	}
}
