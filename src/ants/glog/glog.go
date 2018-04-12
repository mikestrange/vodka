package glog

//日志目录
import (
	"ants/gutil"
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

var m_lv int = 0
var log_size int = 0
var channel chan *logItem

//第一个引用
func LogAndRunning(lv int, size int) {
	if channel != nil {
		return
	}
	println("<<log init>>")
	m_lv = lv
	channel = make(chan *logItem, size)
	go func() {
		for {
			item, ok := <-channel
			if ok {
				handleItem(item)
			} else {
				break
			}
		}
	}()
}

func begin() {
	//没有设置情况下
	LogAndRunning(LOG_DEBUG, 10000)
}

func handleItem(item *logItem) {
	str := fmt.Sprintf(item.str, item.args...)
	str = fmt.Sprintf("%s [%s] %s", gutil.FromtALL(), logMap[item.lv], str)
	fmt.Println(str)
}

//始于
func output(lv int, str string, args []interface{}) {
	begin()
	channel <- &logItem{lv, str, args}
}

//本地日志
func Debug(str string, args ...interface{}) {
	if LOG_DEBUG >= m_lv {
		output(LOG_DEBUG, str, args)
	}
}

func Info(str string, args ...interface{}) {
	if LOG_INFO >= m_lv {
		output(LOG_INFO, str, args)
	}
}

func Warn(str string, args ...interface{}) {
	if LOG_WARN >= m_lv {
		output(LOG_WARN, str, args)
	}
}

func Error(str string, args ...interface{}) {
	if LOG_ERROR >= m_lv {
		output(LOG_ERROR, str, args)
	}
}
