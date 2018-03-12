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

func Launch() {
	go gate.ServerLaunch(config.GATE_PORT)
	go logon.ServerLaunch(config.LOGIN_PORT)
	go world.ServerLaunch(config.WORLD_PORT)
	go chat.ServerLaunch(config.CHAT_PORT)
	go game.ServerLaunch(config.GAME_PORT)
}
