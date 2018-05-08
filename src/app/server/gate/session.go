package gate

import (
	"ants/gnet"
)

type GateSession struct {
	gnet.NetContext
	Player *UserData
}

func NewSession() *GateSession {
	this := new(GateSession)
	this.Player = new(UserData)
	return this
}

func (this *GateSession) LoginOk() {
	this.Player.Status = LOGON_OK
}

func (this *GateSession) LoginOut() {
	this.Player.Status = LOGON_OUT
}

func (this *GateSession) KickOut() {
	this.Player.Status = LOGON_OUT
}

func (this *GateSession) LoginWait() {
	this.Player.Status = LOGON_WAIT
}

func (this *GateSession) IsLogin() bool {
	return this.Player.Status == LOGON_OK
}

func (this *GateSession) Timeout() bool {
	return false
}

func (this *GateSession) IsBegin() bool {
	return this.Player.Status > LOGON_NULL
}

func (this *GateSession) IsLogout() bool {
	return this.Player.Status > LOGON_OK
}

func (this *GateSession) IsLive() bool {
	return this.Player.Status > LOGON_NULL && this.Player.Status < LOGON_OUT
}
