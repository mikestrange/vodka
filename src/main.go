package main

import "app"
import "app/server"

import "ants/gutil"
import "ants/glog"

func main() {
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
		//server.Launch("") //启动服务器
		client_input()
	}
	//test()
	gutil.Add(1)
	gutil.Wait()
}

func client_input() {
	glog.Input(func(str []string) {
		switch str[0] {
		case "exit":
			gutil.ExitSystem(1)
		case "in":
			//test_send()
			go app.Test_login_send(glog.Int(1, str, 1))
		case "out":
			go app.Test_remove_player(glog.Int(1, str, 1))
		case "on":
			go app.Test_get_online()
		case "all":
			go app.Test_send_all()
		case "test":
			go app.Test(glog.Int(1, str, 1))
		case "test2":
			go app.Test_max_login(glog.Int(1, str, 1))
		case "ser":
			server.Launch(glog.Str(2, str, "all"))
		}
	})
}
