package world

import "ants/gnet"
import "ants/gutil"
import "ants/actor"

//世界角色
type GamePlayer struct {
	UserID    int    //玩家ID
	GateID    int    //连接的网关ID
	SessionID uint64 //世界唯一的会话id,用于却别同一用户不同连接
	//其他条件
	Status     int    //游戏状态
	AppID      int    //登入的平台
	RegTime    int64  //登陆的时间（秒）
	UpdateTime uint64 //刷新时间
}

func NewPlayer(pack gnet.IByteArray) *GamePlayer {
	this := new(GamePlayer)
	this.UnPack(pack)
	return this
}

func (this *GamePlayer) UnPack(pack gnet.IByteArray) {
	pack.ReadValue(&this.UserID, &this.GateID, &this.SessionID)
	this.RegTime = gutil.GetTimer()
}

//世界玩家列表
var players map[int]*GamePlayer = make(map[int]*GamePlayer)

//管理
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
		actor.Main().Send(player.GateID, block(player))
	}
}
