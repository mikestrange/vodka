package core

//分级处理错误
type IAgent interface {
	Handle(interface{})
}

//静态函数
type funcActor struct {
	handle func(interface{})
}

//静态处理
func NewAgent(f func(interface{})) IAgent {
	return &funcActor{handle: f}
}

func (this *funcActor) Handle(event interface{}) {
	this.handle(event)
}
