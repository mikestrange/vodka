package model

//用户基本数据
import "soy/db/link"
import "soy/vat"
import "math/rand"

type PlayerModel struct {
	BaseModel
}

func _init_player() {
	player := new(PlayerModel)
	var str string = "未命名" + vat.IntToString(rand.Int()%1000)
	println("设置名字:", str)
	player.ChangeName(str, 1008)
	ret := player.GetPlayerInfo(1008)
	println("获取用户:", ret.String("name"))
}

//获取用户数据
func (this *PlayerModel) GetPlayerInfo(uid int32) link.IRowItem {
	tx := this.Sql().Begin()
	defer tx.Commit()
	ret := tx.ReadQuery("SELECT * FROM player WHERE uid=?", uid)
	if ret == nil || ret.Empty() {
		return nil
	}
	return ret.Row(0)
}

//改名
func (this *PlayerModel) ChangeName(name string, uid int32) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("UPDATE player SET name=? WHERE uid=?", name, uid)
}

//移除一个用户数据(本地必须保存日志)
func (this *PlayerModel) RemovePlayerInfo(uid int32) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("DELETE FROM player WHERE uid=?", uid)
}

//插入一个用户(由账号注册那里生成)
//func (this *PlayerModel) InsertPlayer(uid int32) int {
//	tx := this.Sql().Begin()
//	defer tx.Commit()
//	return tx.WriteQuery("INSERT INTO player (uid) VALUES (?)", uid)
//}
