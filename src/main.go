package main

import "fmt"

import "app"
import "app/server"

import "fat/gutil"

//import "fat/gsys"
import "fat/glog"

//import "fat/gnet"

func main() {
	t := gutil.GetNano()

	glog.LogAndRunning(glog.LOG_DEBUG, 100000)
	glog.Debug("运行路径=%s", gutil.Pwd())
	gutil.TraceData()
	if gutil.MatchSys(1, "cli") {
		client_input()
	} else if gutil.MatchSys(1, "ser") {
		server.Launch(glog.Str(2, gutil.SysArgs(), "gate"))
	} else if gutil.MatchSys(1, "test") {
		app.Test(glog.Int(2, gutil.SysArgs(), 0))
	} else {
		//server.Launch() //启动服务器
		client_input()
	}
	val := gutil.GetNano() - t
	t1, t2 := gutil.IntTofloat(val, 1000)
	println(t1, t2, val)

	fmt.Println(gutil.GetNano()-t, gutil.TimeNanoStr(gutil.GetNano()-t))

	gutil.Add(1)
	gutil.Wait()
}

func client_input() {
	glog.GoAndRunningInput(func(str []string) {
		switch str[0] {
		case "exit":
			gutil.ExitSystem(1)
		case "in":
			app.Test_login_send(glog.Int(1, str, 1))
		case "out":
			app.Test_remove_player(glog.Int(1, str, 1))
		case "on":
			app.Test_get_online()
		case "all":
			go app.Test_send_all()
		case "test":
			app.Test(glog.Int(1, str, 1))
		case "test2":
			app.Test_max_login(glog.Int(1, str, 1))
		}
	})
}
