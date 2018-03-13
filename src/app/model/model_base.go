package model

import "soy/db/link"

func Launch() {
	println("连接Mysql服务器")
}

type BaseModel struct {
}

func (this *BaseModel) Sql() link.IMysqlClient {
	return link.Mysql()
}
