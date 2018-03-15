package world

import "ants/gnet"
import "ants/gutil"

//世界角色
type GamePlayer struct {
	UserID    int32  //玩家ID
	GateID    int32  //连接的网关ID
	SessionID uint64 //世界唯一的会话id,用于却别同一用户不同连接
	//其他条件
	Status     int32  //游戏状态
	AppID      int32  //登入的平台
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

func (this *GamePlayer) SerID() int {
	return int(this.GateID)
}

func (this *GamePlayer) UID() int {
	return int(this.UserID)
}

//世界玩家列表
var players map[int]*GamePlayer = make(map[int]*GamePlayer)

//管理
func GetUser(uid int32) (*GamePlayer, bool) {
	player, ok := players[int(uid)]
	return player, ok
}

func SetUser(player *GamePlayer) (*GamePlayer, bool) {
	uid := player.UID()
	old, ok := players[uid]
	players[uid] = player
	return old, ok
}

func RemoveUser(uid int32) (*GamePlayer, bool) {
	player, ok := players[int(uid)]
	if ok {
		delete(players, int(uid))
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
		router.Send(player.SerID(), block(player))
	}
}
