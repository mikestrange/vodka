package glog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//输入
func Input(call func([]string)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[Exit Input err] ", err)
			} else {
				fmt.Println("[Exit Input ok]")
			}
		}()
		fmt.Println("[Join Input ok]")
		inputReader := bufio.NewReader(os.Stdin)
		for {
			if str, err := inputReader.ReadString('\n'); err == nil {
				str = strings.TrimSpace(str)
				if str == "exit" {
					os.Exit(6)
				} else {
					call(strings.Split(str, " "))
				}
			} else {
				fmt.Println("Input Err:", err.Error())
				break
			}
		}
	}()
}
