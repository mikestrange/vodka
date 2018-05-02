package logon

import "ants/base"
import "ants/lib/gredis"
import "ants/lib/gsql"
import "fmt"

var redis gredis.IConn
var mysql gsql.IConn

//只用来获取用户的数据，并发处理
func init_dber() {
	//redis = gredis.NewConn()
	mysql, _ = gsql.NewConnAddr(gsql.ToAddr("root", "120.77.149.74:3306", "123456", "game_master"), 20)
}

func check_user(uid int, PassWord string) int {
	return 0
	err_code := 0
	pwd, ok := gredis.Str(redis, gredis.ToUser(uid, "pwd"))
	if ok {
		fmt.Println("redis 获取密码:", pwd, ",用户=", PassWord)
		if pwd != PassWord {
			err_code = 1
		}
	} else {
		row := mysql.QueryRow("select pwd from account where uid = ?", uid)
		err := row.Scan(&pwd)
		if err == nil {
			fmt.Println("mysql 获取密码:", pwd, ",用户=", PassWord)
			if pwd != PassWord {
				err_code = 1
			}
			//写入redis
			redis.SetUser(uid, "pwd", pwd, 0)
		} else {
			err_code = 1
			fmt.Println("Scan Err:", err)
		}
	}
	return err_code
}

//获取用户数据
func get_user_info(uid int) []byte {
	return []byte{}
	data, ok := gredis.Bytes(redis, gredis.ToUser(uid, "player.info"))
	if ok {
		//		pack := base.NewByteArrayWithBytes(data)
		//		pack.SetBegin()
		//		name := pack.ReadString()
		//		exp := pack.ReadInt()
		//		money := pack.ReadInt64()
		//		vipexp := pack.ReadInt()
		//		viptype := pack.ReadInt()
		//		pion := pack.ReadInt()
		//		fmt.Println("Redis Get Info:", name, exp, money, vipexp, viptype, pion)
		return data
	} else {
		row := mysql.QueryRow("select name,exp,money,vipexp,viptype,pion from player where uid = ?", uid)
		var name string
		var money int64
		var exp, vipexp, viptype, pion int
		err := row.Scan(&name, &exp, &money, &vipexp, &viptype, &pion)
		if err == nil {
			pack := base.NewByteArrayWithVals(&name, &exp, &money, &vipexp, &viptype, &pion)
			redis.SetUser(uid, "player.info", pack.Bytes(), 0)
			return pack.Bytes()
		} else {
			println("获取用户数据失败")
		}
	}
	return []byte{}
}
