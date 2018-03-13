package model

import "soy/vat"

func _init_game_log() {
	game_log := new(GameLogModel)
	game_log.CreateLogTable(100, 1)
	game_log.InsertGameLog(100, 1, "我敢")
	game_log.InsertMainLog(1, 100, 1, 0, "127.0.0.1", "test")
}

//游戏数据
type GameLogModel struct {
	BaseModel
}

//可能重复？
func (this *GameLogModel) CreateLogTable(roomid int, idx int) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	//建立一个牌局(可能重复？)
	var str string = "CREATE TABLE game_log_" + vat.IntToString(roomid) + "_" + vat.IntToString(idx)
	str += "(`id` INT NOT NULL AUTO_INCREMENT,"
	str += "`record_json` TEXT(1024) NULL,"
	str += " `time` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,"
	str += "PRIMARY KEY (`id`))"
	tx.WriteQuery(str)
	//插入到主表
}

//插入到主表里面?InsertMainLog(1001, 1, 0, "127.0.0.1","test")
func (this *GameLogModel) InsertMainLog(gameid int, roomid int, idx int, rtype int, ip string, desc string) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	sql := "INSERT INTO game_log (game_id,room_id,room_idx,room_type,room_ip,log_addr,desc_t) VALUES (?,?,?,?,?,?,?)"
	addr := "game_log_" + vat.IntToString(roomid) + "_" + vat.IntToString(idx)
	tx.WriteQuery(sql, gameid, roomid, idx, rtype, ip, addr, desc)
}

func (this *GameLogModel) InsertGameLog(roomid int, idx int, json string) {
	tx := this.Sql().Begin()
	defer tx.Commit()
	var str string = "INSERT INTO game_log_" + vat.IntToString(roomid) + "_" + vat.IntToString(idx)
	str += "(record_json) VALUES (?)"
	tx.WriteQuery(str, json)
}
