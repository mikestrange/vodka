package conf

import (
	"fmt"
)

const host string = "120.77.149.74"

//const host string = "127.0.0.1"

type RouteVo struct {
	Port  int
	Addr  string
	Name  string
	Topic int //消息派送 actor ID
	Type  int
}

func toAddr(port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

//ports
const (
	PORT_GATE  = 8081
	PORT_LOGIN = 7100
	PORT_WORLD = 7200
	PORT_GAME  = 7300
	PORT_HALL  = 7400
	PORT_CHAT  = 7500
	PORT_DATA  = 7600
)

//topics
const (
	TOPIC_SELF   = 0
	TOPIC_LOGON  = 1
	TOPIC_WORLD  = 2
	TOPIC_GAME   = 3
	TOPIC_HALL   = 4
	TOPIC_CHAT   = 5
	TOPIC_DATAS  = 6
	TOPIC_CLIENT = 10 //派送给用户的
)

//gate topics
const (
	TOPIC_GATE = 0 //动态的
)

//分布式服务器
var serMap map[int]*RouteVo

func init() {
	serMap = map[int]*RouteVo{
		PORT_GATE:  &RouteVo{PORT_GATE, toAddr(PORT_GATE), "[网关服务器]", TOPIC_GATE, 0},
		PORT_LOGIN: &RouteVo{PORT_LOGIN, toAddr(PORT_LOGIN), "[登录服务器]", TOPIC_LOGON, 0},
		PORT_WORLD: &RouteVo{PORT_WORLD, toAddr(PORT_WORLD), "[世界服务器]", TOPIC_WORLD, 0},
		PORT_GAME:  &RouteVo{PORT_GAME, toAddr(PORT_GAME), "[游戏服务器]", TOPIC_GAME, 0},
		PORT_HALL:  &RouteVo{PORT_HALL, toAddr(PORT_HALL), "[大厅服务器]", TOPIC_HALL, 0},
		PORT_CHAT:  &RouteVo{PORT_CHAT, toAddr(PORT_CHAT), "[聊天服务器]", TOPIC_CHAT, 0},
		PORT_DATA:  &RouteVo{PORT_DATA, toAddr(PORT_DATA), "[数据服务器]", TOPIC_DATAS, 0},
	}
}

func GetRouter(port int) *RouteVo {
	v, ok := serMap[port]
	if ok {
		return v
	}
	return nil
}

func EachVo(block func(*RouteVo)) {
	for _, v := range serMap {
		block(v)
	}
}
