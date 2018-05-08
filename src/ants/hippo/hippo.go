package hippo

//河马分布式引擎
import "ants/core"

var global_network CellActor

func init() {
	core.RunAndThrowBox(&global_network, nil)
}

//网络节点
func Host() core.IBox {
	return &global_network
}

//内部为一个链接
type CellActor struct {
	core.BaseBox
}

func (this *CellActor) OnReady() {
	this.SetAgent(this)
	this.SetBlock(this.OnMessage)
	this.SetName("网络节点主程序")
}

func (this *CellActor) OnMessage(event interface{}) {

}
