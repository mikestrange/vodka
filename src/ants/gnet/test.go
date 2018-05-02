package gnet

import "ants/core"
import "ants/base"

//test
func _init() {
	test()
}

func test() {
	//服务进程
	ref := core.RunAndThrowBox(core.NewBox(new(BoxRef), "net box"), nil, func() {
		println("ref close")
		base.Sleep(1000)
		//test()
	})
	//服务器进程
	tcp, ok := RunAndThrowServer(new(TCPServer), 8081, func(conn interface{}) IAgent {
		tx := NewContext(conn)
		println("todo")
		tx.SetReceiver(func(b []byte) {
			println("msg")
		})
		return tx
	}, func() {
		println("server close")
		//通知其他服务器删除
		ref.Die()
	})
	//>>
	if ok {
		//通知其他服务器注册
	}
	base.Sleep(1000)
	tcp.Die()
}

type BoxRef struct {
}

func (this *BoxRef) Handle(event interface{}) {

}

func (this *BoxRef) OnEnd() {

}
