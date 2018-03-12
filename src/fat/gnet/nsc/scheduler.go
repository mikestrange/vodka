package nsc

import (
	"fat/gsys"
	"sync/atomic"
)

//Network remote scheduler manager(网络远程调度管理)
//Distributed network scheduler manager(分布式网络调度管理)

//远程调度
type IRemoteScheduler interface {
	ResetSession()
	UsedSessionID() uint64
	Data() IDataRoute
	ListenRouter(IRouter) bool
	Send(int, interface{})
	SendAll(interface{})
}

//远程调度器
type RemoteScheduler struct {
	gsys.Locked
	routes  map[int]IRouter
	data    IDataRoute
	session uint64
	handle  RemoteBlock
}

func NewRemoteScheduler(data IDataRoute, block RemoteBlock) IRemoteScheduler {
	this := new(RemoteScheduler)
	this.InitRemoteScheduler(data, block)
	return this
}

func (this *RemoteScheduler) InitRemoteScheduler(data IDataRoute, block RemoteBlock) {
	this.InitLocked()
	this.routes = make(map[int]IRouter)
	this.handle = block
	this.data = data
}

//自己的路由数据
func (this *RemoteScheduler) Data() IDataRoute {
	return this.data
}

//设置节点以及回调函数
func (this *RemoteScheduler) ListenRouter(route IRouter) bool {
	this.Lock()
	defer this.Unlock()
	idx := route.Data().RouteID()
	if _, ok := this.routes[idx]; ok {
		return false
	}
	//统一回调
	if this.handle != nil {
		route.SetHandle(this.handle)
	}
	this.routes[idx] = route
	return true
}

//通知(所有能够接收到的)
func (this *RemoteScheduler) Send(topic int, data interface{}) {
	var list []IRouter
	this.Lock()
	for i := range this.routes {
		if this.routes[i].Data().HasTopic(topic) {
			list = append(list, this.routes[i])
		}
	}
	this.Unlock()
	for i := range list {
		list[i].Push(data)
	}
}

func (this *RemoteScheduler) SendAll(data interface{}) {
	var list []IRouter
	this.Lock()
	for i := range this.routes {
		list = append(list, this.routes[i])
	}
	this.Unlock()
	for i := range list {
		list[i].Push(data)
	}
}

func (this *RemoteScheduler) UsedSessionID() uint64 {
	return atomic.AddUint64(&this.session, 1)
}

func (this *RemoteScheduler) ResetSession() {
	atomic.SwapUint64(&this.session, 0)
}
