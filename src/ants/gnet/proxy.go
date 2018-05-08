package gnet

//基础的代理
type BaseProxy struct {
	NetContext
}

func NewProxy(conn interface{}) *BaseProxy {
	return new(BaseProxy)
}
