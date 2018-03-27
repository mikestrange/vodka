package glog

//日志目录
import (
	"fmt"
)

const (
	LOG_DEBUG = 1
	LOG_INFO  = 2
	LOG_WARN  = 3
	LOG_ERROR = 4
)

var m_lv int = 0
var channel chan string

//第一个引用
func LogAndRunning(lv int, size int) {
	if channel != nil {
		return
	}
	m_lv = lv
	channel = make(chan string, size)
	go func() {
		for {
			str := <-channel
			fmt.Println(str)
		}
	}()
}

func output(str string, args ...interface{}) {
	channel <- fmt.Sprintf(str, args...)
}

//本地日志
func Debug(str string, args ...interface{}) {
	if LOG_DEBUG >= m_lv {
		output("[DEBUG]"+str, args...)
	}
}

func Info(str string, args ...interface{}) {
	if LOG_INFO >= m_lv {
		output("[INFO]"+str, args...)
	}
}

func Warn(str string, args ...interface{}) {
	if LOG_WARN >= m_lv {
		output("[WARN]"+str, args...)
	}
}

func Error(str string, args ...interface{}) {
	if LOG_ERROR >= m_lv {
		output("[ERROR]"+str, args...)
	}
}

//报错
func Assert(ok bool, str string, args ...interface{}) {
	if ok {
		Error(str, args...)
		panic(fmt.Sprintf(str, args...))
	}
}
