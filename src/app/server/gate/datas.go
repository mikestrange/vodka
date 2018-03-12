package gate

import "fat/gnet"

//游戏绑定玩家
type GamePlayer struct {
	Conn   gnet.INetContext
	Player *UserData
}

const (
	LOGON_WAIT  = 0
	LOGON_OK    = 1 //登陆成功
	LOGON_ERROR = 2 //登陆错误
	USER_LIMIT  = 3 //被封
)

//玩家信息
type UserData struct {
	UserID    int32  //玩家ID
	PassWord  string //登录密码
	UserName  string //玩家名称
	SessionID uint64 //世界唯一的会话id,用于却别同一用户不同连接
	GateID    int32  //服务器ID(一般指游戏服务器)
	Status    int32  //游戏状态
	AppID     int32  //登入的平台
	RegTime   int64  //登陆的时间（秒）
	//RPG信息
	Scene SceneData
}

func (this *UserData) ServerID() int {
	return int(this.GateID)
}

//RPG
type SceneData struct {
	Map    int32
	Status int16
	X      int16
	Y      int16
	Z      int16
}
