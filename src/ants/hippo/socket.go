package hippo

//套接字
type ISocket interface {
	Connect(string) bool //连接
	Close() bool         //关闭
}

type Socket struct {
	connected bool
	tx        IContext
}

func (this *Socket) Connect(addr string) bool {
	if !this.connected {
		if tx, ok := Dial(addr); ok {
			this.connected = true
			this.tx = tx
			return true
		}
	}
	return false
}

func (this *Socket) Close() bool {
	this.connected = false
	if this.tx != nil {
		return this.tx.Close()
	}
	return false
}

func (this *Socket) Run() {
	for {
		if msg, ok := this.tx.ReadMsg(); ok {
			for i := range msg {
				println(msg[i])
			}
		} else {
			break
		}
	}
}
