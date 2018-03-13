package model

func _init_game_data() {
	//game data
	game := new(GameDataModel)
	game.UpdateExp(1008, 12)
	game.UpdateMoney(1008, 1)
	game.UpdateGold(1008, 23)
	game.CleanUser(1008)
}

//游戏数据
type GameDataModel struct {
	BaseModel
}

func (this *GameDataModel) UpdateMoney(uid int32, money int64) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("UPDATE game_info SET money=money+? WHERE uid=?", money, uid)
}

func (this *GameDataModel) UpdateExp(uid int32, exp int) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("UPDATE game_info SET exp=exp+? WHERE uid=?", exp, uid)
}

func (this *GameDataModel) UpdateGold(uid int32, gold int32) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("UPDATE game_info SET gold=gold+? WHERE uid=?", gold, uid)
}

//清空用户数据（必须保存日志）
func (this *GameDataModel) CleanUser(uid int32) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	tx.WriteQuery("UPDATE game_info SET gold=0, exp=0, money=0 WHERE uid=?", uid)
}
