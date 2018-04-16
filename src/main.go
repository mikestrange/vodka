package main

import "fmt"

import "app"
import "app/server"

import "ants/gutil"
import "ants/gnet"
import "ants/glog"
import "ants/actor"
import "ants/kernel"

//import _ "ants/kernel"

//import "ants/lib/gredis"
//import "ants/lib/gsql"

var ser gnet.INetServer

func test() {
	//抛出盒子运行
	actor.RunAndThrowActor(actor.Main(), actor.NewFunc(func(args ...interface{}) {
		fmt.Println(args...)
	}, nil), nil)

	actor.Main().Router("13", "13")
	//
	kernel.Go(func() {
		println("执行")
		//panic("no")
	}).Catch(func(err interface{}) {
		//输出错误
		fmt.Println("处理错误2:", err)
	}).Die(func() {
		println("最终素材")
	}).Done()

	kernel.Try(func() {
		panic("123")
	}).Catch(func(err interface{}) {
		fmt.Println("处理错误1")
	}).Die(func() {

	}).Done()

	ena := []byte{1, 3, 5, 7}
	//
	m := []byte{1, 3, 4, 5}
	for i := range m {
		m[i] = (m[i] ^ ena[i]) + 1
	}

	for i := range m {
		m[i] = (m[i] ^ ena[i]) - 1
	}

	fmt.Println(m)
}

func main() {
	glog.Debug("运行路径=%s", gutil.Pwd())
	//gutil.TraceData()
	if gutil.Match(1, "cli") {
		//client_input()
	} else if gutil.Match(1, "ser") {
		server.Launch(gutil.Str(2, gutil.OsArgs(), "all"))
	} else if gutil.Match(1, "test") {
		go app.Test(gutil.Int(2, gutil.OsArgs(), 1))
	} else {
		//go gredis.Test()
		//go gsql.Test()
	}
	test()
	client_input()
	gutil.Add()
	gutil.Wait()
}

func client_input() {
	glog.Input(func(str []string) {
		switch str[0] {
		case "in":
			go app.Test_login_send(gutil.Int(1, str, 1), gutil.Str(2, str, ""))
		case "out":
			go app.Test_remove_player(gutil.Int(1, str, 1))
		case "on":
			go app.Test_get_online()
		case "all":
			for i := 0; i < 10; i++ {
				gutil.Sleep(50)
				go func() {
					app.Test_send_all()
				}()
			}
		case "test":
			go app.Test(gutil.Int(1, str, 1))
		case "test2":
			go app.Test_max_login(gutil.Int(1, str, 1))
		case "ser":
			server.Launch(gutil.Str(2, str, "all"))
		}
	})
}
