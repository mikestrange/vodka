package chat

import "app/command"
import "ants/gcode"
import "app/conf"

//头
//touid int32, session uint64
func pack_join_table(touid int, session uint64, cid int, uid int) interface{} {
	return gcode.NewPackTopic(command.CLIENT_JOIN_CHANNEL, conf.TOPIC_CLIENT, touid, session, cid, uid)
}

func pack_exit_table(touid int, session uint64, cid int, uid int) interface{} {
	return gcode.NewPackTopic(command.CLIENT_QUIT_CHANNEL, conf.TOPIC_CLIENT, touid, session, cid, uid)
}

func pack_message(touid int, session uint64, cid int, uid int, ctype int, msg string) interface{} {
	return gcode.NewPackTopic(command.CLIENT_NOTICE_CHANNEL, conf.TOPIC_CLIENT, touid, session, cid, uid, int16(ctype), msg)
}
