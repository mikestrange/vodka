package cluster

import (
	"ants/gnet"
	"ants/gsys"
	"ants/gutil"
	"sync/atomic"
)

const ALL_TOPIC = -1

//回调方式
type ClusterBlock func(INodeRouter, interface{})

//推送方法()
type INodeRouter interface {
	OnAdd()
	OnRemove()
	Data() IDataRoute
	Push(...interface{}) bool
	Pull() (interface{}, bool)
	Parent() INetCluster
	DoSelf(interface{})
	//私有
	setParent(INetCluster)
}

//分布式基础 (IEventCluster)
type INetCluster interface {
	//自身就是一个节点
	INodeRouter
	//其他节点
	ResetSession()
	UsedSessionID() uint64
	//当前
	AddRouter(INodeRouter)
	RemoveRouter(INodeRouter)
	IndexOf(INodeRouter) int
	//推送子节点
	Send(int, ...interface{})
	SendAll(...interface{})
	//私有
	doEvent(INodeRouter, interface{})
}

//远程调度器
type NetCluster struct {
	gsys.Locked
	session uint64
	data    IDataRoute
	routes  gutil.IArrayObject
	handle  ClusterBlock
	parent  INetCluster
	client  gnet.INetContext
}

func NewCluster(data IDataRoute) INetCluster {
	return NewMainCluster(data, nil)
}

func NewClusterPort(port int) INetCluster {
	return NewCluster(NewDataRoute(port))
}

func NewMainCluster(data IDataRoute, handle ClusterBlock) INetCluster {
	this := new(NetCluster)
	this.InitRemoteScheduler(data, handle)
	return this
}

func (this *NetCluster) InitRemoteScheduler(data IDataRoute, handle ClusterBlock) {
	this.routes = gutil.NewArray()
	this.handle = handle
	this.data = data
}

//设置节点以及回调函数
func (this *NetCluster) AddRouter(val INodeRouter) {
	this.Lock()
	if this.routes.IndexOf(val) == gutil.NOT_VALUE {
		val.setParent(this)
		this.routes.Push(val)
		this.Unlock()
		val.OnAdd()
	} else {
		this.Unlock()
	}
}

func (this *NetCluster) RemoveRouter(val INodeRouter) {
	this.Lock()
	if this.routes.DelVal(val) {
		val.setParent(nil)
		this.Unlock()
		val.OnRemove()
	} else {
		this.Unlock()
	}
}

func (this *NetCluster) IndexOf(val INodeRouter) int {
	return this.routes.IndexOf(val)
}

func (this *NetCluster) Send(topic int, args ...interface{}) {
	list := this.getChildrens(topic)
	for i := range list {
		list[i].Push(args...)
	}
}

func (this *NetCluster) SendAll(args ...interface{}) {
	list := this.getChildrens(ALL_TOPIC)
	for i := range list {
		list[i].Push(args...)
	}
}

func (this *NetCluster) setParent(parent INetCluster) {
	this.parent = parent
}

//处理(交给顶级处理)
func (this *NetCluster) doEvent(node INodeRouter, data interface{}) {
	if this.parent != nil {
		this.parent.doEvent(this, data)
	} else if this.handle != nil {
		this.handle(node, data)
	} else {
		println("no handle rotuer:", this.data.Name())
	}
}

func (this *NetCluster) UsedSessionID() uint64 {
	return atomic.AddUint64(&this.session, 1)
}

func (this *NetCluster) ResetSession() {
	atomic.SwapUint64(&this.session, 0)
}

//interface INodeRouter(处理方式就在这里)
func (this *NetCluster) Push(args ...interface{}) bool {
	this.Lock()
	ok := this.socket_proxy(args...)
	this.Unlock()
	return ok
}

func (this *NetCluster) OnAdd() {

}

func (this *NetCluster) OnRemove() {

}

func (this *NetCluster) DoSelf(data interface{}) {
	this.doEvent(this, data)
}

func (this *NetCluster) Parent() INetCluster {
	return this.parent
}

func (this *NetCluster) Pull() (interface{}, bool) {
	return nil, false
}

func (this *NetCluster) Data() IDataRoute {
	return this.data
}

//privates
func (this *NetCluster) getChildrens(topic int) []INodeRouter {
	var nodes []INodeRouter
	this.Lock()
	list := this.routes.Elements()
	for i := range list {
		node := list[i].(INodeRouter)
		if topic == ALL_TOPIC || node.Data().HasTopic(topic) {
			nodes = append(nodes, node)
		}
	}
	this.Unlock()
	return nodes
}

//成功的方法
func (this *NetCluster) socket_proxy(args ...interface{}) bool {
	if this.client == nil {
		if tx, ok := gnet.Socket(this.Data().Addr()); ok {
			this.client = tx
			tx.SetHandle(func(event int, bits []byte) {
				this.DoSelf(gnet.NewPackBytes(bits))
			})
			go func() {
				tx.Run()
				tx.OnClose(tx.Close())
				this.Lock()
				this.client = nil
				this.Unlock()
			}()
		} else {
			return false
		}
	}
	this.client.Send(args...)
	return true
}
