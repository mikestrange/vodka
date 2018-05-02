package gnet

//基础的代理
type BaseProxy struct {
	NetContext
}

func NewProxy(conn interface{}) *BaseProxy {
	this := new(BaseProxy)
	this.SetConn(conn)
	return this
}
