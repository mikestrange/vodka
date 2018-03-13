package texas

import "fat/gnet"

type PlayerVo struct {
	UserID   int32
	GateID   int32
	Session  uint64
	UserName string
	//用户信息
	PlayerInfo interface{}
}

func NewPlayerVo(uid int32, gid int32, sid uint64, pack gnet.ISocketPacket) *PlayerVo {
	this := new(PlayerVo)
	this.UserID = uid
	this.GateID = gid
	this.Session = sid
	//获取一些用户数据
	return this
}

func (this *PlayerVo) SerID() int {
	return int(this.GateID)
}

func (this *PlayerVo) UID() int {
	return int(this.UserID)
}

func (this *PlayerVo) CheckSession(session uint64) bool {
	return this.Session == session
}
