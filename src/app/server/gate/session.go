package gate

import (
	"ants/base"
	"ants/gcode"
	"ants/gnet"
)

type GateSession struct {
	gnet.NetContext
	Player *UserData
}

func NewSession(conn interface{}) *GateSession {
	this := new(GateSession)
	this.SetConn(conn)
	this.Listen(this, 1024, 3)
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

func (this *GateSession) OnMessage(code int, data interface{}) {
	if code == gnet.EVENT_CONN_READ {
		this.DoHandle(base.ToBytes(data))
	} else if code == gnet.EVENT_CONN_HEARTBEAT {
		if this.IsLogin() {
			this.Send(gcode.NewPackArgs(gnet.EVENT_CONN_HEARTBEAT))
		} else {
			println("关闭啊")
			this.Close()
		}
	} else if code == gnet.EVENT_CONN_SEND {
		this.Conn().WriteBytes(base.ToBytes(data))
	}
}
