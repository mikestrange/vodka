package gate

import "ants/gnet"
import "app/command"

type GateSession struct {
	gnet.BaseProxy
	Player *UserData
}

func NewSession(conn interface{}) *GateSession {
	this := new(GateSession)
	this.SetConn(conn)
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

func (this *GateSession) IsBegin() bool {
	return this.Player.Status > LOGON_NULL
}

func (this *GateSession) IsLogout() bool {
	return this.Player.Status > LOGON_OK
}

func (this *GateSession) IsLive() bool {
	return this.Player.Status > LOGON_NULL && this.Player.Status < LOGON_OUT
}

func (this *GateSession) OnClose() {
	refLogic.Router(this, gnet.NewPackArgs(command.CLIENT_LOGOUT))
}
