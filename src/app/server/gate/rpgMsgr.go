package gate

import (
	"app/config"
	"fat/gnet"
	//"fmt"
)

func on_move(tx gnet.INetContext, pack gnet.ISocketPacket) {
	if data, ok := check_logon(tx); ok {
		x, y, z := pack.ReadShort(), pack.ReadShort(), pack.ReadShort()
		data.Scene.X = x
		data.Scene.Y = y
		data.Scene.Z = z
		router.Send(config.TOPIC_WORLD, packet_world_notices(pack.Cmd(), data.UserID, pack.GetBody()))
	}
}

func on_action(tx gnet.INetContext, pack gnet.ISocketPacket) {
	if data, ok := check_logon(tx); ok {
		action := pack.ReadShort()
		data.Scene.Status = action
		router.Send(config.TOPIC_WORLD, packet_world_notices(pack.Cmd(), data.UserID, pack.GetBody()))
	}
}

func on_result(tx gnet.INetContext, pack gnet.ISocketPacket) {
	//结果
	if data, ok := check_logon(tx); ok {
		println("通知玩家结果:", data.UserID)
	}
}
