package gnet

type INetContext interface {
	IConn
}

type context struct {
	NetConn
}

func NewConn(conn interface{}) INetContext {
	this := new(context)
	this.SetConn(conn)
	return this
}
