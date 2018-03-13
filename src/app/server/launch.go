package server

import "app/config"

//
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
		go gate.ServerLaunch(config.GATE_PORT)
	case "login":
		go logon.ServerLaunch(config.LOGIN_PORT)
	case "game":
		go game.ServerLaunch(config.GAME_PORT)
	case "world":
		go world.ServerLaunch(config.WORLD_PORT)
	case "chat":
		go chat.ServerLaunch(config.CHAT_PORT)
	default:
		Launch("gate")
		Launch("login")
		Launch("world")
		Launch("chat")
		Launch("login")
	}
}
