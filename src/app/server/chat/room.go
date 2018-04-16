package chat

import "ants/actor"
import "ants/gnet"
import "app/command"
import "fmt"

type ChatTable struct {
	actor.BaseBox
	//args
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
		this.Main().Send(player.SerID(), block(player))
	}
}

func (this *ChatTable) OnClose() {

}

func (this *ChatTable) OnMessage(args ...interface{}) {
	header := args[0].(*GameHeader)
	pack := args[1].(gnet.ISocketPacket)
	switch pack.Cmd() {
	case command.CLIENT_JOIN_CHANNEL:
		this.on_join_channel(header, pack)
	case command.CLIENT_QUIT_CHANNEL:
		this.on_quit_channel(header, pack)
	case command.CLIENT_NOTICE_CHANNEL:
		this.on_message(header, pack)
	}
}

//message
func (this *ChatTable) on_join_channel(header *GameHeader, packet gnet.ISocketPacket) {
	cid := this.channel_id
	player := &GameUser{Player: header, UserName: packet.ReadString()}
	old, ok := this.AddUser(player)
	if ok {
		fmt.Println("Enter Chat ok: uid=", player.Player.UserID, ",cid=", cid, ", gate=", player.SerID())
		this.NoticeAllUser(func(data *GameUser) interface{} {
			return pack_join_table(data.Player.UserID, data.Player.SessionID, cid, player.Player.UserID)
		})
	} else {
		//不用管已经被踢掉的用户
		fmt.Println("Enter Chat kill# user = ", header.UserID, ", cid=", cid, ",session=", old.Player.SessionID)
	}
}

func (this *ChatTable) on_quit_channel(header *GameHeader, packet gnet.ISocketPacket) {
	cid := this.channel_id
	if player, ok := this.RemoveUser(header.UserID); ok {
		fmt.Println("Exit Chat Ok# user=", header.UserID, ", cid=", cid)
		//通知>所有
		this.NoticeAllUser(func(data *GameUser) interface{} {
			return pack_exit_table(data.Player.UserID, data.Player.SessionID, cid, player.Player.UserID)
		})
		//通知>自己
		psend := pack_exit_table(player.Player.UserID, player.Player.SessionID, cid, player.Player.UserID)
		this.Main().Send(player.SerID(), psend)
	} else {
		fmt.Println("Exit Chat Err# no user [", header.UserID, "]in cid=", cid)
	}
}

func (this *ChatTable) on_message(header *GameHeader, packet gnet.ISocketPacket) {
	cid := this.channel_id
	mtype, message := packet.ReadShort(), packet.ReadString()
	if player, ok := this.GetUser(header.UserID); ok {
		//通知>所有
		fmt.Println("Notice Chat Ok: user=", header.UserID, ", cid=", cid, ", size=", len(message))
		this.NoticeAllUser(func(data *GameUser) interface{} {
			return pack_message(data.Player.UserID, data.Player.SessionID, cid, player.Player.UserID, mtype, message)
		})
	} else {
		fmt.Println("Notice Chat Err: no user [", header.UserID, "] in cid=", cid, ", size=", len(message))
	}
}
