package model

import "soy/db/link"

func _init_account() {
	account := new(AccountModel)
	account.RegUser(0, 101, "openid", "123", "127.0.0.1", "测试啊")
	account.BindUserID(1009, 1002)
}

//用户基本数据
type AccountModel struct {
	BaseModel
}

func (this *AccountModel) ChangePwd(oldpwd string, uid int32, pwd string) bool {
	tx := this.Sql().Begin()
	defer tx.Commit()
	if tx.WriteQuery("UPDATE account SET pwd=? WHERE uid=? and pwd=?", pwd, uid, oldpwd) > 0 {
		return true
	}
	return false
}

func (this *AccountModel) CheckLogon(uid int32, pwd string) bool {
	tx := this.Sql().Begin()
	defer tx.Commit()
	if tx.ReadQuery("SELECT * FROM account WHERE uid=? and pwd=?", uid, pwd).Size() > 0 {
		return true
	}
	return false
}

func (this *AccountModel) GetUserInfo(uid int32) link.IDataResult {
	tx := this.Sql().Begin()
	defer tx.Commit()
	return tx.ReadQuery("SELECT * FROM account WHERE uid=?", uid)
}

/*
 *	账号注册 appid = 1, openid = 账号名称 appname
	appStore appid=100, openid = mac
	安卓 appid = 200, openid =mac
	//---
	uid int32, appid int32, openid string, ip string, desc string, appname string
*/
func (this *AccountModel) RegUser(uid int, appid int, openid string, pwd string, regip string, desc string) int {
	if uid == 0 {
		println("获取注册id失败")
		return -1
	}
	tx := this.Sql().Begin()
	//插入账号 [uid, 101, "openid", "pwd", "regip", "desc_t"]
	ret0 := tx.WriteQuery("INSERT INTO account (uid,appid,openid,pwd,regip,desc_t) VALUES (?,?,?,?,?,?)",
		uid, appid, openid, pwd, regip, desc)
	//生成账号数据
	ret1 := tx.WriteQuery("INSERT INTO player (uid) VALUES (?)", uid)
	//生成游戏数据
	ret2 := tx.WriteQuery("INSERT INTO game_info (uid) VALUES (?)", uid)
	//判断是否都成功
	if ret0 > 0 && ret1 > 0 && ret2 > 0 {
		tx.Commit()
		return 0
	}
	tx.Rollback()
	return -1
}

//绑定账号(也就是绑定uid)(当前的财产转移)
func (this *AccountModel) BindUserID(uid int, touid int) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	//绑定失败无所谓
	tx.WriteQuery("UPDATE account SET binduid=? WHERE uid=? and binduid=0", touid, uid)
}
