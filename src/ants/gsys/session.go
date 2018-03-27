package gsys

import "sync/atomic"

type Session struct {
	session uint64
}

func (this *Session) UsedSessionID() uint64 {
	return atomic.AddUint64(&this.session, 1)
}

func (this *Session) ResetSession() {
	atomic.SwapUint64(&this.session, 0)
}

//主要的
var MainSession = new(Session)
