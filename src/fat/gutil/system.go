package gutil

import "fmt"
import "time"

//系统启动的时间
var BeginTime int64

//时间格式
const go_day_string = "2006-01-02"
const go_time_string = "2006-01-02 15:04:05"

//时间参数
const TIME_DAY = 60 * 60 * 24
const TIME_HOUR = 60 * 60

//获取的毫秒
func GetTimer() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//获取运行时间(毫秒)
func GetSysTimer() int64 {
	return GetTimer() - BeginTime
}

//获取纳秒(测试用)
func GetNano() int64 {
	return time.Now().UnixNano() / int64(time.Nanosecond)
}

//所有格式时间
func FromtALL() string {
	return time.Now().Format(go_time_string)
}

//获取当前0点的时间 时间戳(秒)
func GetDayBeginTime() int64 {
	tstr := time.Now().Format(go_day_string)
	if val, err := time.Parse(go_day_string, tstr); CheckSucceed(err) {
		//值得一提的是t.Unix()这里获取的是当天早上8:00的时间
		//fmt.Println("获取的是当天早上8:00的时间秒:", val.Unix())
		return val.Unix() - TIME_HOUR*8
	}
	return 0
}

const (
	Nanosecond  int64 = 1                  //纳秒
	Microsecond       = 1000 * Nanosecond  //微妙
	Millisecond       = 1000 * Microsecond //豪妙
	Second            = 1000 * Millisecond //秒
	Minute            = 60 * Second        //分钟
	Hour              = 60 * Minute
)

//传入毫秒
func TimeStr(idx int64) string {
	return TimeNanoStr(idx / Millisecond)
}

//传入纳秒
func TimeNanoStr(idx int64) string {
	if idx < Microsecond {
		return fmt.Sprintf("%d纳秒", idx)
	} else if idx < Millisecond {
		t1, t2 := IntTofloat(idx, Microsecond)
		return fmt.Sprintf("%d.%d微秒", t1, t2)
	} else if idx < Second {
		t1, t2 := IntTofloat(idx, Millisecond)
		return fmt.Sprintf("%d.%d毫秒", t1, t2)
	} else if idx < Minute {
		t1, t2 := IntTofloat(idx, Second)
		return fmt.Sprintf("%d.%d秒", t1, t2)
	} else if idx < Hour {
		t1, t2 := IntTofloat(idx, Minute)
		return fmt.Sprintf("%d.%d分钟", t1, t2)
	}
	t1, t2 := IntTofloat(idx, Hour)
	return fmt.Sprintf("%d.%d小时", t1, t2)
}

//精确度2位
func IntTofloat(val int64, sub int64) (int, int) {
	top := val / sub
	pot := (val - top*sub) / (sub / 100)
	return int(top), int(pot)
}
