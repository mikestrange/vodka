package gate

//用户状态
const (
	LOGON_NULL  = 0
	LOGON_WAIT  = 1
	LOGON_OK    = 2 //登陆成功
	LOGON_KICK  = 3 //被踢
	LOGON_OUT   = 4 //用户退出
	LOGON_ERROR = 5 //登陆错误
	USER_LIMIT  = 6 //被封
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
}

func (this *UserData) ServerID() int {
	return int(this.GateID)
}
