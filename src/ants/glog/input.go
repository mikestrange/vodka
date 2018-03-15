package glog

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InputBlock func([]string)

//系统输出
func system_output(msg string) bool {
	switch msg {
	case "exit\n":
		os.Exit(1)
		return true
	case "mem\n":
		//vat.TraceMemStat()
		return true
	case "disk\n":
		//vat.TraceDiskUsage()
		return true
	}
	return false
}

//回车输入
func loop(call InputBlock) {
	println("[godark input open]")
	inputReader := bufio.NewReader(os.Stdin)
	for {
		if str, err := inputReader.ReadString('\n'); err == nil {
			str = strings.TrimSpace(str)
			call(strings.Split(str, " "))
		} else {
			println("Input Err:", err.Error())
			break
		}
	}
	println("[Exit Input]")
}

func Input(call InputBlock) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("input err over:", err)
			}
		}()
		loop(call)
	}()
}

//获取参数
func Str(idx int, str []string, def string) string {
	if idx >= len(str) {
		return def
	}
	return str[idx]
}

func Int(idx int, str []string, def int) int {
	if idx >= len(str) {
		return def
	}
	if val, err := strconv.Atoi(str[idx]); err == nil {
		return val
	}
	return def
}
