package chat

import "app/command"
import "app/config"
import "fat/gnet"

//头
//touid int32, session uint64
func pack_join_table(touid int32, session uint64, cid int32, uid int32) gnet.IBytes {
	return gnet.NewPacketWithTopic(command.CLIENT_JOIN_CHANNEL, config.TOPIC_CLIENT, touid, session, cid, uid)
}

func pack_exit_table(touid int32, session uint64, cid int32, uid int32) gnet.IBytes {
	return gnet.NewPacketWithTopic(command.CLIENT_QUIT_CHANNEL, config.TOPIC_CLIENT, touid, session, cid, uid)
}

func pack_message(touid int32, session uint64, cid int32, uid int32, ctype int16, msg string) gnet.IBytes {
	return gnet.NewPacketWithTopic(command.CLIENT_NOTICE_CHANNEL, config.TOPIC_CLIENT, touid, session, cid, uid, ctype, msg)
}
