package game

import "ants/base"

//会话头
type GameHeader struct {
	UserID    int
	GateID    int
	SessionID uint64
	TableID   int
}

func NewHeader(pack base.IByteArray) *GameHeader {
	this := new(GameHeader)
	this.UnPack(pack)
	return this
}

func (this *GameHeader) UnPack(pack base.IByteArray) {
	pack.ReadValue(&this.UserID, &this.GateID, &this.SessionID)
}
