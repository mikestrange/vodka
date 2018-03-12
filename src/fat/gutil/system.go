package gutil

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
