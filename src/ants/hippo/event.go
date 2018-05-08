package hippo

type IEvent interface {
	Perform()
}

//发送事件
type SendEvent struct {
	tx    IContext
	conn  IConn
	bytes []byte
}

func newEvent(tx IContext, conn IConn, b []byte) IEvent {
	return &SendEvent{tx, conn, b}
}

func (this *SendEvent) Perform() {
	if !this.conn.SendBytes(this.bytes) {
		this.tx.CloseSign(CLOSE_SIGN_ERROR)
	}
}
