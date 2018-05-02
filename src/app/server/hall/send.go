package hall

import "ants/gcode"
import "app/conf"
import "app/command"

func pack_change_name(code int16, uid int, session uint64, name string) interface{} {
	return gcode.NewPackTopic(command.CLIENT_CHANGE_NAME, conf.TOPIC_CLIENT, uid, session, code, name)
}
