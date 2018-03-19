package logon

import "ants/lib/gredis"
import "ants/lib/gsql"

var redis gredis.IConn
var mysql gsql.IConn

//只用来获取用户的数据，并发处理
func init_dber() {
	redis = gredis.NewConn()
	mysql = gsql.NewConn()
}
