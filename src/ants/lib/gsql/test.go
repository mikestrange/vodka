package gsql

import (
	"ants/gutil"
	"fmt"
)

func init() {

}

func Test() {
	//打开数据库
	//只需要链接一次，内部会建立新链接
	conn, ok := NewConnAddr(ToAddr("root", "127.0.0.1:3306", "123456", "user_info"), 20)
	if !ok {
		return
	}
	defer conn.Close()

	//关闭数据库，db会被多个goroutine共享，可以不调用
	//查询数据，指定字段名，返回sql.Rows结果集
	//	for i := 0; i < 1; i++ {
	//		go func() {
	//			ret, _ := conn.Exec("update account set status = ? where uid = ? and status = 2", 0, 1001)
	//			upd_nums, _ := ret.RowsAffected()
	//			fmt.Println("影响:", upd_nums)
	//		}()
	//	}
	gutil.Sleep(100)

	//查询一行数据
	form := conn.Form("select * from logon_log")
	if !form.Empty() {
		t := gutil.GetNano()
		fmt.Println("获取数据 form:", form.Item(0).Str("ip"), gutil.GetNano()-t)
	}
	gutil.Sleep(1000)
}
