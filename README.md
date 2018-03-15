# ants(框架) 没有包含任何第三方框架 能直接运行
----------begin [一些你能用到的实用的lib]

go get golang.org/x/sys/unix

go get github.com/shirou/gopsutil

go get github.com/garyburd/redigo

go get github.com/go-sql-driver/mysql

----------end

###测试
可以直接运行sh make.sh来启动服务器(包括app下面的几个模块，game模块请自行注销 app/server/launch.go)

##客户端测试
1，打开终端进入该目录，输入; go run main.go test 任意数字
2，go run main.go cli 下输入 in(链接一个用户) on(查看在线人数) out(踢出一个链接) all(给所有用户发消息)