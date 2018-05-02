package gnet

import "ants/core"

var global_network core.IBox

func init() {
	global_network = core.RunAndThrowBox(core.NewBase(OnMessage, "网络节点主程序"), nil)
}

//网络节点
func Host() core.IBox {
	return global_network
}

func OnMessage(event interface{}) {
	//do it
}

//内部是一个socket
type Cell struct {
}
