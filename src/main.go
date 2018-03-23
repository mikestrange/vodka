package main

import "app"
import "app/server"

import "ants/gutil"
import "ants/glog"

//import "ants/gnet"

import "ants/actor"

//import "ants/lib/gredis"
//import "ants/lib/gsql"

func test() {
	actor.Main.Added(actor.NewTestActor())
	for {
		actor.Main.Send(1, 1, 3, 4, 5, 6, 7)
		gutil.Sleep(1000)
	}
}

func main() {
	glog.LogAndRunning(glog.LOG_DEBUG, 100000)
	glog.Debug("运行路径=%s", gutil.Pwd())
	gutil.TraceData()
	if gutil.MatchSys(1, "cli") {
		//client_input()
	} else if gutil.MatchSys(1, "ser") {
		server.Launch(glog.Str(2, gutil.SysArgs(), "gate"))
	} else if gutil.MatchSys(1, "test") {
		go app.Test(glog.Int(2, gutil.SysArgs(), glog.Int(3, gutil.SysArgs(), 1)))
	} else {
		//server.Launch("") //启动服务器
		//go gredis.Test()
		//go gsql.Test()
	}
	//var arr interface{} = []interface{}{1, 2, 3, "12312312"}
	//fmt.Println(string(gutil.JsonEncode(arr)))
	test()
	client_input()
	gutil.Add(1)
	gutil.Wait()
}

func client_input() {
	//println(gutil.GetTimer(), int(gutil.GetTimer()-1521172200000))
	glog.Input(func(str []string) {
		switch str[0] {
		case "exit":
			gutil.ExitSystem(1)
		case "in":
			go app.Test_login_send(glog.Int(1, str, 1), glog.Str(2, str, ""))
		case "out":
			go app.Test_remove_player(glog.Int(1, str, 1))
		case "on":
			go app.Test_get_online()
		case "all":
			for i := 0; i < 1; i++ {
				go func(idx int) {
					gutil.Sleep(idx * 10)
					app.Test_send_all()
				}(i)
			}
		case "test":
			go app.Test(glog.Int(1, str, 1))
		case "test2":
			go app.Test_max_login(glog.Int(1, str, 1))
		case "ser":
			server.Launch(glog.Str(2, str, "all"))
		}
	})
}
