package game

//房间和用户的管理
type TablePlayer struct {
	UserID    int32  //玩家ID
	SessionID uint64 //世界唯一的会话id,用于却别同一用户不同连接
	GateID    int32  //连接的网关ID
	Status    int32  //游戏状态(具体的游戏状态)
	RegTime   int64  //登陆的时间（秒）
	RoomID    int    //游戏房间
	//
	Tables map[int]interface{}
}

func (this *TablePlayer) ServerID() int {
	return int(this.GateID)
}

func (this *TablePlayer) UID() int {
	return int(this.UserID)
}

//个人房间操作
func (this *TablePlayer) HasTable(roomid int) bool {
	_, ok := this.Tables[roomid]
	return ok
}

func (this *TablePlayer) SetTable(roomid int) bool {
	this.Tables[roomid] = true
	return true
}

func (this *TablePlayer) UnSetTable(roomid int) bool {
	delete(this.Tables, roomid)
	return true
}

//用户管理可能几百万人
var users map[int]*TablePlayer = make(map[int]*TablePlayer)
