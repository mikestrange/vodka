package game

const (
	STATE_ONLINE = 0
	STATE_DROPS  = 1
)

//房间和用户的管理
type DataPlayer struct {
	UserID    int    //玩家ID
	SessionID uint64 //世界唯一的会话id,用于却别同一用户不同连接
	GateID    int    //连接的网关ID
	Status    int    //状态(具体的游戏状态) 0在线，1离线
	RoomID    int    //游戏房间
	//进入房间列表
	rooms map[int]int64 //0空闲，1锁定，2进入
}

func (this *DataPlayer) init() {
	this.Status = STATE_ONLINE
	this.rooms = make(map[int]int64)
}

func (this *DataPlayer) ServerID() int {
	return int(this.GateID)
}

func (this *DataPlayer) UID() int {
	return int(this.UserID)
}

func (this *DataPlayer) Enter(roomid int) bool {
	//记录一个时间（超时能重新进入）
	_, ok := this.rooms[roomid]
	if ok {
		return false
	}
	this.rooms[roomid] = 1
	return true
}

func (this *DataPlayer) Exit(roomid int) bool {
	_, ok := this.rooms[roomid]
	if ok {
		delete(this.rooms, roomid)
		return true
	}
	return false
}

func (this *DataPlayer) Check(gid int, sid uint64) bool {
	return this.GateID == gid && this.SessionID == sid
}

//离线
func (this *DataPlayer) Die() bool {
	return this.Status == STATE_DROPS
}

func (this *DataPlayer) Empty() bool {
	return len(this.rooms) == 0
}

func (this *DataPlayer) HomeList() []int {
	var list []int
	for i := range this.rooms {
		list = append(list, i)
	}
	return list
}
