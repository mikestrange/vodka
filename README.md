# ants(轻量级框架) 
----------begin-----------

 [一些你能用到的实用的lib]

go get golang.org/x/sys/unix

go get github.com/shirou/gopsutil

go get github.com/garyburd/redigo

go get github.com/go-sql-driver/mysql

----------end--------------


###测试服务器
1，直接运行sh make.sh启动服务器(gate, chat, logon, world, game)

2，game模块在开发阶段，请自行注销(在目录app/server/launch.go里面)


##客户端测试(i表示任意数字 +表示空格)

1，打开终端进入该目录，输入; go run main.go test+i

2，go run main.go cli 启动后输入:

in+i(链接一个用户) 

on(查看在线人数) 

out+i(踢出一个链接) 

all(给所有用户发消息)