package actor

//具体被响应(event)
type IBoxActor interface {
	//监听消息
	OnMessage(...interface{})
	//退出的时候调度
	OnDie()
}

//盒子推送，不具备其他方法
type IBoxPusher interface {
	Router(...interface{}) bool //摆渡
}

//盒子代理(run)
type IBoxRef interface {
	IBoxPusher
	SetActor(IBoxActor)    //监听对象
	Make(interface{}) bool //制造chan
	OnReady()              //运行时准备
	PerformRunning()       //执行运行
	OnRelease()            //调度释放
	Die()                  //自身关闭
	setFather(IBoxSystem)  //设置上级(不能继承)
	Father() IBoxSystem    //父运行
	Main() IBoxSystem      //顶峰
}

//盒子集合(node)
type IBoxSystem interface {
	//独立的环境
	IBoxRef
	//节点控制
	ActorOf(interface{}, IBoxRef) bool     //添加并且运行
	UnRef(IBoxRef) bool                    //移除，未close
	CloseRef(interface{}) (IBoxRef, bool)  //移除并且close
	FindRef(interface{}) (IBoxRef, bool)   //查找
	Send(interface{}, ...interface{}) bool //发送
	Broadcast(...interface{})              //通知所有
	CloseAll()                             //关闭移除所有(禁止子类调用)
}
