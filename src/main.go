package main

import "fmt"

import "app"
import "app/server"

import "ants/gutil"
import "ants/glog"
import "ants/gnet"

import "ants/actor"

//import "ants/lib/gredis"
//import "ants/lib/gsql"

var ser gnet.INetServer

func test() {
	fmt.Println("test")
	//
	actor.Main.SetActor(actor.NewFunc(func(args ...interface{}) {
		fmt.Println(args...)
	}, nil))
	actor.RunWithActor(actor.Main)
	actor.Main.Router(112, 123)
}

func main() {
	glog.LogAndRunning(glog.LOG_DEBUG, 100000)
	glog.Debug("运行路径=%s", gutil.Pwd())
	//gutil.TraceData()
	if gutil.MatchSys(1, "cli") {
		//client_input()
	} else if gutil.MatchSys(1, "ser") {
		server.Launch(glog.Str(2, gutil.SysArgs(), "all"))
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
			//i := 0; i < 1000000; i++
			for i := 0; i < 10; i++ {
				gutil.Sleep(50)
				go func() {
					app.Test_send_all()
				}()
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
