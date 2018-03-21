package cluster

import "ants/gsys"

type NodeSet map[INode]int
type NodeBlock func(interface{}, interface{})

type INode interface {
	//接受者对象
	Taker() interface{}
	UnTaker() interface{}
	SetTaker(interface{}) bool
	//回调方法
	SetHandle(NodeBlock)
	//操作
	AddSet(INode)
	RemSet(INode)
	Parent() INode
	OnRemove()
	OnAdded()
	//处理（顶级）
	Done(interface{})
	//私有，避免别人调度
	setParent(INode)
	doEvent(interface{}, interface{})
}

type Node struct {
	gsys.Locked
	Nodes  NodeSet //公开给继承者
	parent INode
	handle NodeBlock
	client interface{}
}

func (this *Node) init() {
	this.Nodes = make(NodeSet)
}

//handle
func (this *Node) SetHandle(h NodeBlock) {
	this.handle = h
}

//parent
func (this *Node) Parent() INode {
	v := this.parent
	return v
}

func (this *Node) setParent(val INode) {
	this.parent = val
}

//add
func (this *Node) AddSet(n INode) {
	this.Lock()
	if _, ok := this.Nodes[n]; !ok {
		this.Nodes[n] = len(this.Nodes)
		n.setParent(this)
		this.Unlock()
		n.OnAdded()
	} else {
		this.Unlock()
	}
}

func (this *Node) OnAdded() {

}

//rm
func (this *Node) RemSet(n INode) {
	this.Lock()
	if _, ok := this.Nodes[n]; ok {
		delete(this.Nodes, n)
		n.setParent(nil)
		this.Unlock()
		n.OnRemove()
	} else {
		this.Unlock()
	}
}

func (this *Node) OnRemove() {

}

func (this *Node) Taker() interface{} {
	return this.client
}

func (this *Node) UnTaker() interface{} {
	this.Lock()
	v := this.client
	this.client = nil
	this.Unlock()
	return v
}

func (this *Node) SetTaker(v interface{}) bool {
	this.Lock()
	defer this.Unlock()
	if this.client == nil {
		this.client = v
		return true
	}
	return false
}

//do
func (this *Node) Done(data interface{}) {
	this.doEvent(this.client, data)
}

func (this *Node) doEvent(client interface{}, data interface{}) {
	p := this.parent
	if p != nil {
		p.doEvent(client, data)
	} else if this.handle != nil {
		this.handle(client, data)
	} else {
		println("node not handle")
	}
}
