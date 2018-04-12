package gate

import "ants/gutil"

type LogonMode struct {
	open_list gutil.IArrayObject     //等待列表
	users     map[int32]*GateSession //成功的列表
}

func NewLogonMode() *LogonMode {
	this := new(LogonMode)
	this.InitLogonMode()
	return this
}

func (this *LogonMode) InitLogonMode() {
	this.open_list = gutil.NewArray()
	this.users = make(map[int32]*GateSession)
}

//提交登陆
func (this *LogonMode) CommitLogon(data interface{}) {
	this.open_list.Push(data)
}

//完成登陆 获取登陆链接
func (this *LogonMode) CompleteLogon(uid int32, session uint64) (*GateSession, bool) {
	index, data := this.open_list.SeachValue(func(val interface{}) bool {
		player := val.(*GateSession).Player
		return player.UserID == uid && player.SessionID == session
	})
	if index != gutil.NOT_VALUE {
		this.open_list.DelIndex(index)
		return data.(*GateSession), true
	}
	return nil, false
}

//绑定登陆数据用户
func (this *LogonMode) AddUser(uid int32, data *GateSession) (*GateSession, bool) {
	val, ok := this.users[uid]
	this.users[uid] = data
	if ok {
		return val, true
	}
	return nil, false
}

func (this *LogonMode) GetUser(uid int32) (*GateSession, bool) {
	if val, ok := this.users[uid]; ok {
		return val, true
	}
	return nil, false
}

func (this *LogonMode) GetUserBySession(uid int32, session uint64) (*GateSession, bool) {
	if GetUser, ok := this.users[uid]; ok {
		return GetUser, GetUser.Player.SessionID == session
	}
	return nil, false
}

func (this *LogonMode) RemoveUser(uid int32) (*GateSession, bool) {
	if player, ok := this.users[uid]; ok {
		delete(this.users, uid)
		return player, true
	}
	return nil, false
}

//合法的移除
func (this *LogonMode) RemoveUserWithSession(uid int32, session uint64) (*GateSession, bool) {
	if player, ok := this.users[uid]; ok {
		if player.Player.SessionID == session {
			delete(this.users, uid)
			return player, true
		}
	}
	return nil, false
}
