package cluster

import (
	"ants/gnet"
	"sync/atomic"
)

const ALL_TOPIC = -1

//分布式基础
type INodeRouter interface {
	INode
	//其他节点
	ResetSession()
	UsedSessionID() uint64
	//推送子节点
	Send(int, interface{})
	SendAll(interface{})
	//自身推送
	Push(interface{}) bool
	Pull() (interface{}, bool)
	//路由数据
	Data() IDataRoute
}

//远程调度器
type NodeRouter struct {
	Node
	session uint64
	data    IDataRoute
	tx      gnet.INetContext
}

func NewRouter(data IDataRoute) INodeRouter {
	return NewMainRouter(data, nil)
}

func NewRouterPort(port int) INodeRouter {
	return NewRouter(NewData(port))
}

func NewMainRouter(data IDataRoute, handle NodeBlock) INodeRouter {
	this := new(NodeRouter)
	this.InitRemoteScheduler(data, handle)
	return this
}

func (this *NodeRouter) InitRemoteScheduler(data IDataRoute, handle NodeBlock) {
	this.init()
	this.SetTaker(this)
	this.SetHandle(handle)
	this.data = data
}

func (this *NodeRouter) Send(topic int, data interface{}) {
	list := this.getChildrens(topic)
	for i := range list {
		list[i].Push(data)
	}
}

func (this *NodeRouter) SendAll(data interface{}) {
	list := this.getChildrens(ALL_TOPIC)
	for i := range list {
		list[i].Push(data)
	}
}

func (this *NodeRouter) UsedSessionID() uint64 {
	return atomic.AddUint64(&this.session, 1)
}

func (this *NodeRouter) ResetSession() {
	atomic.SwapUint64(&this.session, 0)
}

//interface INodeRouter(处理方式就在这里)
func (this *NodeRouter) Push(data interface{}) bool {
	this.Lock()
	ok := this.socket_proxy(data)
	this.Unlock()
	return ok
}

func (this *NodeRouter) Pull() (interface{}, bool) {
	return nil, false
}

func (this *NodeRouter) Data() IDataRoute {
	return this.data
}

//privates(自行判断)
func (this *NodeRouter) getChildrens(topic int) []INodeRouter {
	var nodes []INodeRouter
	this.Lock()
	for k := range this.Nodes {
		node := k.(INodeRouter)
		if topic == ALL_TOPIC || node.Data().HasTopic(topic) {
			nodes = append(nodes, node)
		}
	}
	this.Unlock()
	return nodes
}

//成功的方法
func (this *NodeRouter) socket_proxy(data interface{}) bool {
	if this.tx == nil {
		if tx, ok := gnet.Socket(this.Data().Addr()); ok {
			this.tx = tx
			tx.SetHandle(func(b []byte) {
				this.Done(gnet.NewPackBytes(b))
			})
			go func() {
				tx.WaitFor()
				this.Lock()
				this.tx = nil
				this.Unlock()
			}()
		} else {
			return false
		}
	}
	this.tx.Send(data)
	return true
}
