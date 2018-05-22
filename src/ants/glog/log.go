package glog

//系统日志
import (
	"ants/base"
	"ants/core"
	"fmt"
)

const (
	LOG_DEBUG = 1
	LOG_INFO  = 2
	LOG_WARN  = 3
	LOG_ERROR = 4
)

var logMap map[int]string = map[int]string{
	LOG_DEBUG: "DEBUG",
	LOG_INFO:  "INFO",
	LOG_WARN:  "WARN",
	LOG_ERROR: "ERROR",
}

type logItem struct {
	lv   int
	str  string
	args []interface{}
}

var str_rep = "%s [%s] %s"
var box LogBox

type LogBox struct {
	core.Box
	//打印等级
	Lv   int
	size int
}

func (this *LogBox) OnReady() {
	this.SetName("打印日志")
	this.SetAgent(this)
}

func (this *LogBox) Handle(event interface{}) {
	item := event.(*logItem)
	str := base.Format(item.str, item.args...)
	str = fmt.Sprintf(str_rep, base.FromtALL(), logMap[item.lv], str)
	fmt.Println(str)
}

func (this *LogBox) Output(lv int, str string, args []interface{}) {
	if this.Lv > lv {
		return
	}
	this.Push(&logItem{lv, str, args})
}

//日志运行
func LogAndRunning(lv int, size int) {
	core.RunAndThrowBox(&box, size)
	box.Lv = lv
}

//本地日志
func Debug(str string, args ...interface{}) {
	box.Output(LOG_DEBUG, str, args)
}

func Info(str string, args ...interface{}) {
	box.Output(LOG_INFO, str, args)
}

func Warn(str string, args ...interface{}) {
	box.Output(LOG_WARN, str, args)
}

func Error(str string, args ...interface{}) {
	box.Output(LOG_ERROR, str, args)
}
