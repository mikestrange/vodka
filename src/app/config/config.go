package config

import (
	"fat/gnet/nsc"
	"fmt"
)

const host string = "120.77.149.74"

//const host string = "127.0.0.1"

//服务器端口
const (
	GATE_PORT  = 8081
	LOGIN_PORT = 7100
	WORLD_PORT = 7200
	GAME_PORT  = 7300
	HALL_PORT  = 7400
	CHAT_PORT  = 7500
	DATA_PORT  = 7600
)

//topic/服务id
const (
	TOPIC_GATE   = 0
	TOPIC_LOGON  = 1
	TOPIC_WORLD  = 2
	TOPIC_GAME   = 3
	TOPIC_HALL   = 4
	TOPIC_CHAT   = 5
	TOPIC_DATAS  = 6
	TOPIC_CLIENT = 10 //派送给用户的
)

//分布式服务器
var serMap map[int]nsc.IDataRoute

func init() {
	serMap = map[int]nsc.IDataRoute{
		GATE_PORT:  nsc.NewDataRouteWithArgs(0, fmt.Sprintf("%s:%d", host, GATE_PORT), "[网关服务器]", TOPIC_GATE),
		LOGIN_PORT: nsc.NewDataRouteWithArgs(1, fmt.Sprintf("%s:%d", host, LOGIN_PORT), "[登陆服务器]", TOPIC_LOGON),
		WORLD_PORT: nsc.NewDataRouteWithArgs(2, fmt.Sprintf("%s:%d", host, WORLD_PORT), "[世界服务器]", TOPIC_WORLD),
		GAME_PORT:  nsc.NewDataRouteWithArgs(3, fmt.Sprintf("%s:%d", host, GAME_PORT), "[游戏服务器]", TOPIC_GAME),
		HALL_PORT:  nsc.NewDataRouteWithArgs(4, fmt.Sprintf("%s:%d", host, HALL_PORT), "[大厅服务器]", TOPIC_HALL),
		CHAT_PORT:  nsc.NewDataRouteWithArgs(5, fmt.Sprintf("%s:%d", host, CHAT_PORT), "[聊天服务器]", TOPIC_CHAT),
		DATA_PORT:  nsc.NewDataRouteWithArgs(6, fmt.Sprintf("%s:%d", host, DATA_PORT), "[数据服务器]", TOPIC_DATAS),
	}
}

func GetDataRouter(port int) nsc.IDataRoute {
	if val, ok := serMap[port]; ok {
		return val
	}
	return nil
}

//注册所有分布式的服务器
func SetServerLists(remote nsc.IRemoteScheduler) {
	for i := range serMap {
		remote.ListenRouter(nsc.NewRouter(serMap[i]))
	}
}
