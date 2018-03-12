package logon

//import "app/model"
import "app/config"
import "app/command"
import "fat/gnet"
import "fmt"

var events map[int]interface{} = map[int]interface{}{
	command.CLIENT_LOGON: on_logon,
}

func on_logon(packet gnet.ISocketPacket) {
	UserID, PassWord, SerID, SessionID := packet.ReadInt(), packet.ReadString(), packet.ReadInt(), packet.ReadUInt64()
	fmt.Println(fmt.Sprintf("Logon Info# uid=%d, session=%v gateid=%d", UserID, SessionID, SerID))
	err_code := 0
	if false {
		//      accout := new(model.AccountModel)
		//		if accout.CheckLogon(UserID, PassWord) {
		//			m_player := new(model.PlayerModel)
		//			ret := m_player.GetPlayerInfo(UserID)
		//			if ret == nil {
		//				err_code = 2
		//			} else {
		//				//写入用户数据
		//			}
		//		} else {
		//			err_code = 1
		//		}
	}
	fmt.Println("Seach Result Code:", err_code, UserID, PassWord, SerID, SessionID)
	var body []byte = []byte{}
	//错误直接返回
	if err_code != 0 {
		router.Send(config.TOPIC_WORLD, pack_logon_result(int16(err_code), UserID, SerID, SessionID, body))
	} else {
		router.Send(config.TOPIC_WORLD, pack_logon_result(int16(err_code), UserID, SerID, SessionID, body))
	}
}

//通知登录结果
func pack_logon_result(code int16, uid int32, gateid int32, session uint64, body []byte) gnet.IBytes {
	psend := gnet.NewPacket()
	psend.WriteBegin(command.SERVER_WORLD_ADD_PLAYER)
	psend.WriteValue(code, uid, gateid, session, body)
	//成功写入资料
	psend.WriteEnd()
	return psend
}
