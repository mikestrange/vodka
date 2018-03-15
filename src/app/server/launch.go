package server

import "ants/conf"

import "app/server/gate"
import "app/server/logon"
import "app/server/world"
import "app/server/chat"
import "app/server/game"

func init() {

}

func Launch(name string) {
	switch name {
	case "gate":
		go gate.ServerLaunch(conf.PORT_GATE)
	case "login":
		go logon.ServerLaunch(conf.PORT_LOGIN)
	case "game":
		go game.ServerLaunch(conf.PORT_GAME)
	case "world":
		go world.ServerLaunch(conf.PORT_WORLD)
	case "chat":
		go chat.ServerLaunch(conf.PORT_CHAT)
	default:
		Launch("gate")
		Launch("login")
		Launch("world")
		Launch("chat")
		Launch("game")
	}
}
