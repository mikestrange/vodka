package hall

import "ants/gnet"
import "ants/conf"
import "app/command"

func pack_change_name(code int16, uid int32, session uint64, name string) gnet.IBytes {
	return gnet.NewPackTopic(command.CLIENT_CHANGE_NAME, conf.TOPIC_CLIENT, uid, session, code, name)
}
