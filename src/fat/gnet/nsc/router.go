package nsc

import "fat/gsys"
import "fmt"
import "fat/gnet"

//回调方式
type RemoteBlock func(IRouter, interface{})

//路由器(推送)
type IRouter interface {
	SetHandle(RemoteBlock)
	Data() IDataRoute
	Push(interface{}) bool
	Pull(interface{}) (interface{}, bool) //这里基本是发送给dbser
}

//路由数据不允许被更换
type Router struct {
	gsys.Locked
	data    IDataRoute
	handle  RemoteBlock
	context gnet.INetContext
}

//class Router
func NewRouter(data IDataRoute) IRouter {
	this := new(Router)
	this.InitRouter(data, nil)
	return this
}

func (this *Router) InitRouter(data IDataRoute, block RemoteBlock) {
	this.InitLocked()
	this.data = data
	this.handle = block
}

func (this *Router) Data() IDataRoute {
	return this.data
}

func (this *Router) Push(data interface{}) bool {
	this.Lock()
	defer this.Unlock()
	if this.context == nil {
		if tx, ok := gnet.NewSocket(this.data.Addr()); ok {
			this.context = tx
			go func() {
				defer this.uncontext()
				gnet.LoopWithHandle(tx, this.thread)
			}()
		} else {
			fmt.Println("连接", this.data.Name(), "服务器失败!")
			return false
		}
	}
	return this.context.Send(data)
}

func (this *Router) SetHandle(handle RemoteBlock) {
	this.handle = handle
}

//这里会直接返回给推送者的socket
func (this *Router) Pull(data interface{}) (interface{}, bool) {
	return nil, false
}

//private function
func (this *Router) thread(tx gnet.INetContext, data interface{}) {
	if this.handle == nil {
		fmt.Println("节点回调未处理:", this.data.Name())
	} else {
		this.handle(this, data)
	}
}

func (this *Router) uncontext() {
	fmt.Println("节点关闭:", this.Data().Name())
	this.Lock()
	this.context = nil
	this.Unlock()
}
