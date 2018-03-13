package model

func _init_logon_log() {
	logon := new(LogonLogModel)
	logon.Logon(10028, "127.0.0.1", 1, "xx-xx-xx-xx", "IOS 10.0.1", "test")
}

//登陆日志
type LogonLogModel struct {
	BaseModel
}

//登陆日志
func (this *LogonLogModel) Logon(uid int32, ip string, appid int, mac string, name string, desc string) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("INSERT INTO logon_log (uid,ip,appid,dev_mac,dev_name,desc_t) VALUES (?,?,?,?,?,?)",
		uid, ip, appid, mac, name, desc)
}

//删除一个用户登陆日志
func (this *LogonLogModel) RemoveLogon(uid int32) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("DELETE FROM logon_log WHERE uid=?", uid)
}

//删除所有登陆日志
func (this *LogonLogModel) RemoveAll() {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("DELETE FROM logon_log")
}
