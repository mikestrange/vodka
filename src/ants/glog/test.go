package glog

import "ants/base"
import "ants/core"
import "os"

func _init() {
	LogAndRunning(LOG_DEBUG, 100)
	Debug("运行路径=%s", base.Pwd())
	//递归建立目录
	err := os.MkdirAll("./temp/log", 0777)
	if err != nil {
		Debug("%s", err)
	} else {
		Debug("Create Directory OK!")
	}
	//
	store := new(LogStore) //NewStore("./temp/log/debug.txt")
	if store.Open("./temp/log/debug.txt") {
		core.Main().Join(11, store)
	}
	store.SetEnd()
	store.Write("12312")
	store.Write("我是谁2")
	store.Write("我是谁3")

}
