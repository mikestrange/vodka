package base

import "sync/atomic"

//全局
var b_session uint64

func UsedSessionID() uint64 {
	return atomic.AddUint64(&b_session, 1)
}

func ResetSession() {
	atomic.SwapUint64(&b_session, 0)
}
