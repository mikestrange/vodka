package world

import "ants/base"

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

func NewPlayer(pack base.IByteArray) *GamePlayer {
	this := new(GamePlayer)
	this.UnPack(pack)
	return this
}

func (this *GamePlayer) UnPack(pack base.IByteArray) {
	pack.ReadValue(&this.UserID, &this.GateID, &this.SessionID)
	this.RegTime = base.Timer()
}
