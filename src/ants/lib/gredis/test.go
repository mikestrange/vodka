package gredis

import (
	"fmt"
)

func init() {

}

func Test() {
	//注意redis不能异步set
	conn := NewConn()
	defer conn.Close()
	//
	conn.SetUser(1, "pwd", "1234321", 0)
	ret, ok := Str(conn, ToUser(1, "pwd"))
	if ok {
		fmt.Println("redis=", ret)
	} else {
		println("什么:", ret)
	}
}
