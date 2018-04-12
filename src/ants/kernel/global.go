package kernel

import (
	"fmt"
	"sync/atomic"
	"time"
)

//###########################闹钟
type clockBlock func(uint64, interface{})

/*
生成一个闹钟，用来跨线程回调
*/
var private_time_idx uint64 = 0

func clockHandler(delay int, count int, data interface{}, callback clockBlock) (uint64, *time.Ticker) {
	timeid := atomic.AddUint64(&private_time_idx, 1)
	timer := SetTimeout(delay, count, func() {
		callback(timeid, data)
	})
	return timeid, timer
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

func After(delay int, block func()) *time.Timer {
	return time.AfterFunc(tdelay(delay)*time.Millisecond, block)
}

func Forever(delay int, block func()) *time.Ticker {
	tm := time.NewTicker(tdelay(delay) * time.Millisecond)
	go func() {
		defer tm.Stop()
		defer check_error()
		for {
			<-tm.C
			block()
		}
	}()
	return tm
}

func SetTimeout(delay int, count int, block func()) *time.Ticker {
	tm := time.NewTicker(tdelay(delay) * time.Millisecond)
	go func() {
		defer tm.Stop()
		defer check_error()
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
	}()
	return tm
}

func check_error() {
	if err := recover(); err != nil {
		fmt.Println("Ticker Timeout Err :", err)
	}
}
