package hippo

//具体环境
type IContext interface {
	//SetConn(IConn)                  //操作的
	Conn() IConn                    //具体链接
	SetCoding(ICoding)              //编码格式
	SendMsg(...interface{}) bool    //写消息
	ReadMsg() ([]interface{}, bool) //读消息
	Close() bool                    //自身关闭

}

//加入协议部分
type Context struct {
	conn IConn
	code ICoding
}

func (this *Context) SetConn(conn IConn) {
	this.conn = conn
	//this.code = NewBigCoding() //默认
}

func (this *Context) Conn() IConn {
	return this.conn
}

func (this *Context) SetCoding(val ICoding) {
	this.code = val
}

func (this *Context) SendMsg(args ...interface{}) bool {
	for i := range args {
		if b, ok := this.code.Encode(args[i]); ok {
			if this.conn.Write(b) {
				return true
			}
		}
	}
	return false
}

func (this *Context) ReadMsg() ([]interface{}, bool) {
	if b, ok := this.conn.Read(this.code.Buffer()); ok {
		return this.code.Decode(b)
	}
	return nil, false
}

func (this *Context) Close() bool {
	return this.conn.Close()
}
