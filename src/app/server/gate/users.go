package gate

import "ants/base"

type LogonMode struct {
	open_list base.IArrayObject    //等待列表
	users     map[int]*GateSession //成功的列表
}

func NewLogonMode() *LogonMode {
	this := new(LogonMode)
	this.init()
	return this
}

func (this *LogonMode) init() {
	this.open_list = base.NewArray()
	this.users = make(map[int]*GateSession)
}

//提交登陆
func (this *LogonMode) CommitLogon(data interface{}) {
	this.open_list.Push(data)
}

//完成登陆 获取登陆链接
func (this *LogonMode) CompleteLogon(uid int, session uint64) (*GateSession, bool) {
	index, data := this.open_list.SeachValue(func(val interface{}) bool {
		player := val.(*GateSession).Player
		return player.UserID == uid && player.SessionID == session
	})
	if index != base.NOT_VALUE {
		this.open_list.DelIndex(index)
		return data.(*GateSession), true
	}
	return nil, false
}

//绑定登陆数据用户
func (this *LogonMode) AddUser(uid int, data *GateSession) (*GateSession, bool) {
	val, ok := this.users[uid]
	this.users[uid] = data
	if ok {
		return val, true
	}
	return nil, false
}

func (this *LogonMode) GetUser(uid int) (*GateSession, bool) {
	if val, ok := this.users[uid]; ok {
		return val, true
	}
	return nil, false
}

func (this *LogonMode) GetUserBySession(uid int, session uint64) (*GateSession, bool) {
	if GetUser, ok := this.users[uid]; ok {
		return GetUser, GetUser.Player.SessionID == session
	}
	return nil, false
}

func (this *LogonMode) RemoveUser(uid int) (*GateSession, bool) {
	if player, ok := this.users[uid]; ok {
		delete(this.users, uid)
		return player, true
	}
	return nil, false
}

//合法的移除
func (this *LogonMode) RemoveUserWithSession(uid int, session uint64) (*GateSession, bool) {
	if player, ok := this.users[uid]; ok {
		if player.Player.SessionID == session {
			delete(this.users, uid)
			return player, true
		}
	}
	return nil, false
}
