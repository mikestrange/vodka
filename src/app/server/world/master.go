package world

import "ants/core"

//世界玩家列表
var players map[int]*GamePlayer

func init() {
	players = make(map[int]*GamePlayer)
}

func GetUser(uid int) (*GamePlayer, bool) {
	player, ok := players[uid]
	return player, ok
}

func SetUser(player *GamePlayer) (*GamePlayer, bool) {
	uid := player.UserID
	old, ok := players[uid]
	players[uid] = player
	return old, ok
}

func RemoveUser(uid int) (*GamePlayer, bool) {
	player, ok := players[uid]
	if ok {
		delete(players, uid)
	}
	return player, ok
}

func GetUserList() []*GamePlayer {
	var list []*GamePlayer
	for _, v := range players {
		list = append(list, v)
	}
	return list
}

func NoticeAllUser(block func(*GamePlayer) interface{}) {
	list := GetUserList()
	for _, player := range list {
		core.Main().Send(player.GateID, block(player))
	}
}
