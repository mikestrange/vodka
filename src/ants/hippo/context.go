package hippo

//环境
type Context struct {
	conn  IConn
	code  ICoding
	actor IActor
}

func NewContext(conn IConn) IContext {
	this := new(Context)
	this.SetConn(conn)
	return this
}

//必须设置
func (this *Context) SetConn(conn IConn) {
	this.conn = conn
}

//必须设置
func (this *Context) SetCoding(code ICoding) {
	this.code = code
}

//必须设置
func (this *Context) SetActor(ator IActor) {
	this.actor = ator
}

func (this *Context) ReadMsg(b []byte) ([]interface{}, error) {
	ret, err := this.conn.ReadBytes(b)
	if err == nil {
		list, derr := this.code.Decode(b[:ret])
		if derr == nil {
			return list, nil
		} else {
			return nil, derr
		}
	}
	return nil, err
}

func (this *Context) CloseOf(args ...interface{}) {
	this.SendMsg(args...)
	this.Close()
}

func (this *Context) Close() {
	this.CloseSign(CLOSE_SIGN_SELF)
}

func (this *Context) Conn() IConn {
	return this.conn
}

//actor
func (this *Context) SendMsg(args ...interface{}) bool {
	if b, err := this.code.Encode(args...); err == nil {
		this.actor.PushEvent(newEvent(this, this.conn, b))
		return true
	}
	return false
}

func (this *Context) CloseSign(code int) {
	this.actor.Exit(code)
}

func (this *Context) Loop() {
	b := this.code.NewBits()
	for {
		if args, err := this.ReadMsg(b); err == nil {
			for i := range args {
				this.actor.PushRead(args[i])
			}
		} else {
			this.CloseSign(CLOSE_SIGN_ERROR)
			break
		}
	}
	this.conn.Close()
}
