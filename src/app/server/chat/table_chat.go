package chat

import "fat/gnet"

type ChatTable struct {
	channel_id   int32
	channel_type int8
	users        map[int]*GameUser
}

func NewTable(cid int32, ctype int8) *ChatTable {
	this := new(ChatTable)
	this.InitTable(cid, ctype)
	return this
}

func (this *ChatTable) InitTable(cid int32, ctype int8) {
	this.users = make(map[int]*GameUser)
	this.channel_id = cid
	this.channel_type = ctype
}

//会直接替换
func (this *ChatTable) AddUser(player *GameUser) bool {
	ok := true
	uid := player.Player.UID()
	if _, ok2 := this.users[uid]; ok2 {
		ok = false
	}
	this.users[uid] = player
	return !ok
}

func (this *ChatTable) UpdateUser(player *GameUser) {

}

func (this *ChatTable) RemoveUser(uid int32) (*GameUser, bool) {
	player, ok := this.users[int(uid)]
	if ok {
		delete(this.users, int(uid))
	}
	return player, ok
}

func (this *ChatTable) GetUser(uid int32) (*GameUser, bool) {
	player, ok := this.users[int(uid)]
	return player, ok
}

func (this *ChatTable) GetUserList() []*GameUser {
	var list []*GameUser
	for _, v := range this.users {
		list = append(list, v)
	}
	return list
}

func (this *ChatTable) NoticeAllUser(block func(*GameUser) gnet.IBytes) {
	list := this.GetUserList()
	//通知所有(无锁)
	for _, player := range list {
		router.Send(player.SerID(), block(player))
	}
}
