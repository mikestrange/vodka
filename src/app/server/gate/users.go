package gate

import "fat/gutil"
import "fat/gsys"

type LogonMode struct {
	gsys.Locked
	open_list gutil.IArrayObject //等待列表
	users     gutil.IHashMap     //成功的列表
}

func NewLogonMode() *LogonMode {
	this := new(LogonMode)
	this.InitLogonMode()
	return this
}

func (this *LogonMode) InitLogonMode() {
	this.InitLocked()
	this.open_list = gutil.NewArray()
	this.users = gutil.NewHashMap()
}

func (this *LogonMode) CommitLogon(data interface{}) {
	this.Lock()
	this.open_list.Push(data)
	this.Unlock()
}

func (this *LogonMode) CompleteLogon(uid int32, session uint64) (*GamePlayer, bool) {
	this.Lock()
	defer this.Unlock()
	index, data := this.open_list.SeachValue(func(val interface{}) bool {
		player := val.(*GamePlayer).Player
		return player.UserID == uid && player.SessionID == session
	})
	if index != gutil.NOT_VALUE {
		this.open_list.DelIndex(index)
		return data.(*GamePlayer), true
	}
	return nil, false
}

//这里是登录后的用户
func (this *LogonMode) AddUser(uid int32, data interface{}) (*GamePlayer, bool) {
	this.Lock()
	defer this.Unlock()
	old := this.users.Set(uid, data)
	if old == nil {
		return nil, false
	}
	return old.(*GamePlayer), true
}

func (this *LogonMode) GetUser(uid int32) (*GamePlayer, bool) {
	this.Lock()
	defer this.Unlock()
	val := this.users.Val(uid)
	if val == nil {
		return nil, false
	}
	return val.(*GamePlayer), true
}

func (this *LogonMode) GetUserBySession(uid int32, session uint64) (*GamePlayer, bool) {
	if player, ok := this.GetUser(uid); ok {
		return player, player.Player.SessionID == session
	}
	return nil, false
}

func (this *LogonMode) RemoveUser(uid int32) (*GamePlayer, bool) {
	this.Lock()
	defer this.Unlock()
	data := this.users.Del(uid)
	if data == nil {
		return nil, false
	}
	return data.(*GamePlayer), true
}
