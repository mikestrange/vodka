package gate

import "ants/gnet"
import "ants/conf"

type GateSession struct {
	gnet.BaseProxy
	Player *UserData
}

func NewSession(tx gnet.INetContext) *GateSession {
	this := new(GateSession)
	this.SetContext(tx)
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
	//发生了登录(成功是否不知道)
	if this.IsLive() {
		//等待列表删除
		logon.CompleteLogon(this.Player.UserID, this.Player.SessionID)
		//登录成功后的删除
		logon.RemoveUserWithSession(this.Player.UserID, this.Player.SessionID)
		//通知世界
		router.Send(conf.TOPIC_WORLD, packet_world_delplayer(this.Player.UserID,
			this.Player.GateID, this.Player.SessionID))
	}
}
