package base

import (
	"fmt"
	"time"
)

func check_error() {
	if err := recover(); err != nil {
		fmt.Println("Ticker Timeout Err :", err)
	}
}

/*
1,一个计时器和回调
2,传入的以毫米作为计算单位
3,当传入小于最低延迟时，设置默认20毫秒，不然报错
*/
const MIN_DELAY = 20

func tdelay(delay int) time.Duration {
	if delay < MIN_DELAY {
		return time.Duration(MIN_DELAY)
	}
	return time.Duration(delay)
}

//一次性计时器
func After(delay int, block func()) *time.Timer {
	return time.AfterFunc(tdelay(delay)*time.Millisecond, block)
}

//永久性
func Forever(delay int, block func()) *time.Ticker {
	tm := time.NewTicker(tdelay(delay) * time.Millisecond)
	TryGo(func() {
		for {
			<-tm.C
			block()
		}
	}, func(ok bool) {
		tm.Stop()
	})
	return tm
}

//次数限制
func SetTimeout(delay int, count int, block func()) *time.Ticker {
	tm := time.NewTicker(tdelay(delay) * time.Millisecond)
	TryGo(func() {
		for {
			<-tm.C
			if count <= 0 {
				block()
			} else {
				count = count - 1
				block()
				if count == 0 {
					break
				}
			}
		}
	}, func(ok bool) {
		tm.Stop()
	})
	return tm
}
