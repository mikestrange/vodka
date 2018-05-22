package base

import "fmt"
import "time"

//系统启动的时间
var b_sys_time int64

//时间格式
const go_day_string = "2006-01-02"
const go_time_string = "2006-01-02 15:04:05.0"

func init() {
	b_sys_time = Timer()
}

func Nano() int64 {
	return time.Now().UnixNano()
}

//获取的毫秒
func Timer() int64 {
	return Nano() / Millisecond
}

//获取运行时间(毫秒)
func SysTimer() int64 {
	return Timer() - b_sys_time
}

//所有格式时间
func FromtALL() string {
	return time.Now().Format(go_time_string)
}

func FromtDay() string {
	return time.Now().Format(go_day_string)
}

/* 获取当天0点的时间>时间戳(秒)
1,值得一提的是t.Unix()这里获取的是当天早上8:00的时间
2,fmt.Println("获取的是当天早上8:00的时间秒:", val.Unix())
*/
func DayTime() int64 {
	if val, err := time.Parse(go_day_string, FromtDay()); err == nil {
		return val.Unix() - 60*60*8
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

//测试函数时间
func TryFun(f func()) string {
	t := Nano()
	f()
	str := NanoStr(Nano() - t)
	//println(str)
	return str
}

//传入纳秒
func NanoStr(idx int64) string {
	if idx < Microsecond {
		return fmt.Sprintf("%d纳秒", idx)
	} else if idx < Millisecond {
		t1, t2 := to2Float(idx, Microsecond)
		return fmt.Sprintf("%d.%d微秒", t1, t2)
	} else if idx < Second {
		t1, t2 := to2Float(idx, Millisecond)
		return fmt.Sprintf("%d.%d毫秒", t1, t2)
	} else if idx < Minute {
		t1, t2 := to2Float(idx, Second)
		return fmt.Sprintf("%d.%d秒", t1, t2)
	} else if idx < Hour {
		t1, t2 := to2Float(idx, Minute)
		return fmt.Sprintf("%d.%d分钟", t1, t2)
	}
	t1, t2 := to2Float(idx, Hour)
	return fmt.Sprintf("%d.%d小时", t1, t2)
}

//传入毫秒
func MicStr(idx int64) string {
	return NanoStr(idx * Millisecond)
}

//精确度2位
func to2Float(val int64, sub int64) (int, int) {
	top := val / sub
	pot := (val - top*sub) / (sub / 100)
	return int(top), int(pot)
}
