package main

import "app"

//import "fmt"
import "app/server"

import "ants/gutil"
import "ants/glog"
import "ants/gnet"

import "ants/lib/gredis"
import "ants/lib/gsql"

var send_size = 10000 //10MB

func test() {
	gutil.Sleep(1000)
	go gnet.ListenAndRunServer(8080, func(session gnet.IBaseProxy) {
		session.Context().SetHandle(func(code int, bits []byte) {
			println("收到:", len(bits))
			//session.Send(bits)
		})
	})

	gutil.Sleep(1)
	if tx, ok := gnet.Socket("127.0.0.1:8080"); ok {
		message := ""
		for i := 0; i < send_size; i++ {
			message += "1234567"
		}
		tx.SetHandle(func(code int, bits []byte) {
			pack := gnet.NewPackBytes(bits)
			println(pack.Cmd(), "消息：", gutil.MemStr(int64(pack.Length())), "字节>", gutil.NanoStr(gutil.GetNano()-pack.ReadInt64()))
		})
		t11 := gutil.GetNano()
		tx.Send(gnet.NewPackArgs(101, gutil.GetNano(), message))
		println("发送:", gutil.NanoStr(gutil.GetNano()-t11))
		tx.Run()
	}
}

func main() {
	glog.LogAndRunning(glog.LOG_DEBUG, 100000)
	glog.Debug("运行路径=%s", gutil.Pwd())
	gutil.TraceData()
	if gutil.MatchSys(1, "cli") {
		client_input()
	} else if gutil.MatchSys(1, "ser") {
		server.Launch(glog.Str(2, gutil.SysArgs(), "gate"))
	} else if gutil.MatchSys(1, "test") {
		app.Test(glog.Int(2, gutil.SysArgs(), int(gutil.GetTimer()-1521172200000)))
	} else {
		//server.Launch("") //启动服务器
		go gredis.Test()
		go gsql.Test()
	}
	test()
	client_input()
	//test()
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
