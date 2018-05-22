package hippo

const (
	CLOSE_CODE_SELF   = 1 //自己关闭
	CLOSE_CODE_ERROR  = 2 //内部发送错误
	CLOSE_CODE_CLIENT = 3 //客户端关闭
)

const (
	EVENT_TYPE_SEND    = 1
	EVENT_TYPE_READ    = 2
	EVENT_TYPE_CLOSE   = 3
	EVENT_TYPE_TIMEOUT = 4
)

//派送对象
type IPusher interface {
	SetCaller(interface{}) //调度者
	Push(IEvent) bool      //推送
}
