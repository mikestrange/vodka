package main

import "ants/base"
import "ants/glog"

import _ "ants/core"
import _ "ants/ghttp"
import _ "ants/gnet"
import _ "ants/gcode"

import "app"
import "app/server"

func init() {
	if base.OsCompare(1, "cli") {
		//client_input()
	} else if base.OsCompare(1, "ser") {
		server.Launch(base.Str(2, base.OsArgs(), "all"))
	} else if base.OsCompare(1, "test") {
		app.Test(base.Int(2, base.OsArgs(), 1))
	} else {
		//go gredis.Test()
		//go gsql.Test()
	}

}

func main() {
	glog.LogAndRunning(glog.LOG_DEBUG, 0)
	vim()
}

func vim() {
	glog.Input(func(str []string) {
		switch str[0] {
		case "in":
			app.Test_login_send(base.Int(1, str, 1), base.Str(2, str, ""))
		case "out":
			app.Test_remove_player(base.Int(1, str, 1))
		case "on":
			app.Test_get_online()
		case "all":
			app.Test_send_all()
		case "test":
			app.Test(base.Int(1, str, 1))
		case "test2":
			app.Test_max_login(base.Int(1, str, 1))
		case "ser":
			server.Launch(base.Str(2, str, "all"))
		}
	})
}
