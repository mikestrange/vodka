# goser
----------begin [一些你能用到的实用的lib]
go get golang.org/x/sys/unix

go get github.com/shirou/gopsutil

go get github.com/garyburd/redigo

go get github.com/go-sql-driver/mysql
----------end

直接sh make.sh能运行服务器

//直接进行客户端测试，每次链接200个(mac一个终端最大这么多)
go run main.go cli 1 


//查看世界多少用户在线
go run main.go online 


//给世界所有用户发送消息
go run main.go all