package logon

//
import "app/command"
import "ants/gnet"
import "ants/conf"
import "fmt"

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_LOGON: on_logon,
}

func on_logon(packet gnet.ISocketPacket) {
	UserID, PassWord, SerID, SessionID := packet.ReadInt(), packet.ReadString(), packet.ReadInt(), packet.ReadUInt64()
	fmt.Println(fmt.Sprintf("Logon Info# uid=%d, session=%v gateid=%d", UserID, SessionID, SerID))
	err_code := check_user(UserID, PassWord)
	fmt.Println("Seach Result Code:", err_code, UserID, PassWord, SerID, SessionID)
	var body []byte = []byte{}
	if err_code == 0 {
		body = get_user_info(int(UserID))
	}
	//错误直接返回
	if err_code != 0 {
		router.Send(conf.TOPIC_WORLD, pack_logon_result(int16(err_code), UserID, SerID, SessionID, body))
	} else {
		router.Send(conf.TOPIC_WORLD, pack_logon_result(int16(err_code), UserID, SerID, SessionID, body))
	}
}

//通知登录结果
func pack_logon_result(code int16, uid int32, gateid int32, session uint64, body []byte) gnet.IBytes {
	return gnet.NewPackArgs(command.SERVER_WORLD_ADD_PLAYER, code, uid, gateid, session, body)
}
