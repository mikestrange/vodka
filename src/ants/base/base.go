package base

import (
	"fmt"
	"os"
	"time"
)

func Exit(code int) {
	os.Exit(code)
}

/*
暂停毫秒(毫秒计算)
*/
func Sleep(val int) {
	time.Sleep(time.Millisecond * time.Duration(val))
}

/*
断言
*/
func Assert(ok bool, str string, args ...interface{}) {
	if ok {
		panic(Format(str, args...))
	}
}

/*
抛出错误
*/
func Throw(str string, args ...interface{}) {
	panic(Format(str, args...))
}

/*
捕获函数
*/
func Try(do func(), ok func(bool)) {
	//defer Catch(ok)
	defer ok(true)
	do()
}

/*
捕获进程并且返回
*/
func TryGo(do func(), ok func(bool)) {
	go func() {
		Try(do, ok)
	}()
}

//捕获错误
func Catch(ok func(bool)) {
	if err := recover(); err != nil {
		fmt.Println("TryGo err##", err)
		ok(false)
	} else {
		ok(true)
	}
}
