package chat

import "ants/glog"
import "ants/core"
import "ants/gcode"
import "ants/gnet"

type ChatTable struct {
	channel_id   int
	channel_type int
	users        map[int]*GameUser
}

func NewTable(cid int, ctype int) *ChatTable {
	this := new(ChatTable)
	this.InitTable(cid, ctype)
	return this
}

func (this *ChatTable) InitTable(cid int, ctype int) {
	this.users = make(map[int]*GameUser)
	this.channel_id = cid
	this.channel_type = ctype
}

//会直接替换
func (this *ChatTable) AddUser(player *GameUser) (*GameUser, bool) {
	uid := player.Player.UserID
	old, ok := this.users[uid]
	this.users[uid] = player
	if ok {
		return old, false
	}
	return nil, true
}

func (this *ChatTable) GetUser(uid int) (*GameUser, bool) {
	if user, ok := this.users[uid]; ok {
		return user, true
	}
	return nil, false
}

func (this *ChatTable) RemoveUser(uid int) (*GameUser, bool) {
	player, ok := this.users[uid]
	if ok {
		delete(this.users, uid)
	}
	return player, ok
}

func (this *ChatTable) GetUserList() []*GameUser {
	var list []*GameUser
	for _, v := range this.users {
		list = append(list, v)
	}
	return list
}

func (this *ChatTable) NoticeAllUser(block func(*GameUser) interface{}) {
	list := this.GetUserList()
	//通知所有(无锁)
	for _, player := range list {
		core.Main().Send(player.SerID(), gnet.NewPack(nil, block(player)))
	}
}

func (this *ChatTable) Close() {

}

//handle
func (this *ChatTable) on_join_channel(header *GameHeader, pack gcode.ISocketPacket) {
	player := &GameUser{Player: header, UserName: pack.ReadString()}
	old, ok := this.AddUser(player)
	if ok {
		glog.Debug("Enter Chat ok: uid=%d gate=%d with cid=%d", player.Player.UserID, player.SerID(), this.channel_id)
		this.NoticeAllUser(func(data *GameUser) interface{} {
			return pack_join_table(data.Player.UserID, data.Player.SessionID, this.channel_id, player.Player.UserID)
		})
	} else {
		glog.Debug("Enter Chat Err[kill user]: uid=%d gate=%d with cid=%d", header.UserID, old.SerID(), this.channel_id)
	}
}

func (this *ChatTable) on_quit_channel(header *GameHeader, pack gcode.ISocketPacket) {
	if player, ok := this.RemoveUser(header.UserID); ok {
		glog.Debug("Exit Chat Ok# uid=%d with cid=%d", header.UserID, this.channel_id)
		//通知>所有
		this.NoticeAllUser(func(data *GameUser) interface{} {
			return pack_exit_table(data.Player.UserID, data.Player.SessionID, this.channel_id, player.Player.UserID)
		})
		//通知>自己
		psend := pack_exit_table(player.Player.UserID, player.Player.SessionID, this.channel_id, player.Player.UserID)
		core.Main().Send(player.SerID(), gnet.NewPack(nil, psend))
	} else {
		glog.Debug("Exit Chat Err[not user]# uid=%d with cid=%d", header.UserID, this.channel_id)
	}
}

func (this *ChatTable) on_message(header *GameHeader, pack gcode.ISocketPacket) {
	mtype, message := pack.ReadShort(), pack.ReadString()
	if player, ok := this.GetUser(header.UserID); ok {
		//通知>所有
		glog.Debug("Msg Chat Ok: uid=%d len=%d with cid=%d", header.UserID, len(message), this.channel_id)
		this.NoticeAllUser(func(data *GameUser) interface{} {
			return pack_message(data.Player.UserID, data.Player.SessionID,
				this.channel_id, player.Player.UserID, mtype, message)
		})
	} else {
		glog.Debug("Msg Chat Err[not user]: uid=%d len=%d with cid=%d", header.UserID, this.channel_id, len(message))
	}
}
