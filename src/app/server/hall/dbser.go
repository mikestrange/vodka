package hall

//import "ants/gnet"
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

func change_name(uid int, name string) bool {
	//redis需要更新？
	ret, ok := mysql.Exec("update player set name = ? where uid = ?", name, uid)
	if ok {
		num, err := ret.RowsAffected()
		if err == nil && num > 0 {
			fmt.Println("修改名称成功")
			return true
		} else {
			fmt.Println("修改名称失败:", err)
		}
	} else {
		fmt.Println("修改名称失败,没有对象")
	}
	return false
}
