package chat

import "app/command"
import "ants/gcode"
import "ants/glog"

var events map[int]interface{} = map[int]interface{}{
	command.SERVER_BUILD_CHANNEL:  on_build_channel,
	command.SERVER_REMOVE_CHANNEL: on_remove_channel,
	command.CLIENT_JOIN_CHANNEL:   on_join_channel,
	command.CLIENT_QUIT_CHANNEL:   on_quit_channel,
	command.CLIENT_NOTICE_CHANNEL: on_message,
}

var tables map[int]*ChatTable

func init() {
	tables = make(map[int]*ChatTable)
}

func on_build_channel(pack gcode.ISocketPacket) {
	tableid, ttype := pack.ReadInt(), pack.ReadByte()
	if _, ok := tables[tableid]; ok {
		glog.Debug("new table err[have table] cid=%d", tableid)
		return
	}
	tables[tableid] = NewTable(tableid, ttype)
	glog.Debug("new table ok: cid=%d type=%d", tableid, ttype)
}

func on_remove_channel(pack gcode.ISocketPacket) {
	tableid := pack.ReadInt()
	if table, ok := tables[tableid]; ok {
		delete(tables, tableid)
		table.Close()
		glog.Debug("del table ok: cid=%d", tableid)
	} else {
		glog.Debug("del table err[not table]: cid=%d ", tableid)
	}
}

//message

func on_join_channel(pack gcode.ISocketPacket) {

}

func on_quit_channel(pack gcode.ISocketPacket) {

}

func on_message(pack gcode.ISocketPacket) {

}

func get_header(pack gcode.ISocketPacket) *GameHeader {
	this := new(GameHeader)
	pack.ReadValue(&this.UserID, &this.GateID, &this.SessionID)
	this.TableID = pack.ReadInt()
	return this
}
