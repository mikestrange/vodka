package glog

import (
	"bufio"
	"fmt"
	"os"
)

type InputBlock func(string)

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
	inputReader := bufio.NewReader(os.Stdin)
	for {
		if msg, err := inputReader.ReadString('\n'); err == nil {
			call(msg)
		} else {
			fmt.Println("Input Err:", err)
			break
		}
	}
	fmt.Println("[Exit Input]")
}

func GoAndRunningInput(call InputBlock) {
	go loop(call)
}
