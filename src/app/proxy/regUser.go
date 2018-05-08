package proxy

import "fmt"
import "ants/base"
import "ants/glog"
import "ants/lib/gsql"

//随机生成用户 开始位置和长度
func RandUids(conn gsql.IConn) {
	conn.Exec("DELETE FROM uids.rand1_uids")
	conn.Exec("truncate uids.rand1_uids") //主健清0

	uids := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		uids[i] = i
	}
	for i := 0; i < 1000; i++ {
		r := base.Random(1000)
		temp := uids[r]
		uids[r] = uids[i]
		uids[i] = temp
	}
	//fmt.Println(uids)
	for i := 0; i < 1000; i++ {
		uid := 1000 + uids[i]
		fmt.Println(uid)
		conn.Exec("INSERT INTO uids.rand1_uids (uid,name,atom,md5) VALUES(?,?,0,'abc123')", uid, base.Format("专业用户%d", uid))
	}
}

func regUser(conn gsql.IConn, md5 string) int {
	tx, ok := conn.Begin()
	if !ok {
		return 1
	}
	//锁表(MD5为唯一的用户识别)
	ret1, err1 := tx.Exec("UPDATE uids.rand1_uids SET atom=1, md5=? WHERE atom=0 limit 1", md5)
	if err1 != nil {
		glog.Debug("%v", err1)
		tx.Rollback()
		return 2
	}
	//注册中(UID和md5必须不重复)
	sql := "INSERT INTO texas_game.user (uid,md5,name)"
	sql += " SELECT uid,md5,name FROM uids.rand1_uids WHERE md5=? AND atom=1"
	sql += " AND not exists(SELECT * FROM texas_game.user WHERE texas_game.user.md5=?)"
	sql += " limit 1"
	ret2, err2 := tx.Exec(sql, md5, md5)
	if err2 != nil {
		glog.Debug("%v", err2)
		tx.Rollback()
		return 3
	}
	//注册完成刷新状态
	ret3, err3 := tx.Exec("UPDATE uids.rand1_uids SET atom=2 WHERE atom=1 and md5=? limit 1", md5)
	if err3 != nil {
		glog.Debug("%v", err3)
		tx.Rollback()
		return 4
	}
	//检查3次是否发生作用
	c1, _ := ret1.RowsAffected()
	c2, _ := ret2.RowsAffected()
	c3, _ := ret3.RowsAffected()
	if c1 == 0 || c2 == 0 || c3 == 0 {
		glog.Debug("reg err:md5=%s", md5)
		tx.Rollback()
		return 5
	}
	glog.Debug("reg ok:md5=%s", md5)
	tx.Commit()
	return 0
}
