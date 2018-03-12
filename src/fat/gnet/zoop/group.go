package gnet

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type NetNodeBlock func(INetNode, interface{})

//服务器组 publish
type IZookeeper interface {
	//生成会话ID
	ResetSession()
	UsedSessionID() uint64
	//本身节点
	Target() INetNode
	//添加节点
	AddNode(INetNode) INetNode
	//移除某个节点
	RemoveNode(int) INetNode
	//获取某个节点
	GetNode(int) INetNode
	//获取某条消息的所有监听节点
	GetNodeByTopic(int) NetNodeList
	//清理所有节点
	CleanNodes() NetNodeList
	HasNode(idx int) bool
	//推送消息
	DoneTopic(int, interface{}) int
	//推送到某个节点
	DoneNode(int, interface{}) int
	//处理对象
	SetHandle(NetNodeBlock)
}

//集群管理
type Zookeeper struct {
	NetNode
	mutex   *sync.Mutex
	session uint64
	group   map[int]INetNode
	handle  NetNodeBlock
}

func NewGroup() IZookeeper {
	this := new(Zookeeper)
	this.InitZookeeper()
	return this
}

func NewGroupWithHandle(block NetNodeBlock) IZookeeper {
	this := NewGroup()
	this.SetHandle(block)
	return this
}

func (this *Zookeeper) InitZookeeper() {
	this.mutex = new(sync.Mutex)
	this.group = make(map[int]INetNode)
}

func (this *Zookeeper) Target() INetNode {
	return this
}

func (this *Zookeeper) Lock() {
	this.mutex.Lock()
}

func (this *Zookeeper) Unlock() {
	this.mutex.Unlock()
}

func (this *Zookeeper) AddNode(node INetNode) INetNode {
	this.Lock()
	idx := node.NodeID()
	old, ok := this.group[idx]
	this.group[idx] = node
	this.Unlock()
	if ok {
		return old
	}
	return nil
}

func (this *Zookeeper) GetNode(idx int) INetNode {
	this.Lock()
	node, ok := this.group[idx]
	this.Unlock()
	if ok {
		return node
	}
	return nil
}

func (this *Zookeeper) GetNodeByTopic(topic int) NetNodeList {
	var nodes NetNodeList
	this.Lock()
	for _, node := range this.group {
		if node.HasTopic(topic) {
			nodes = append(nodes, node)
		}
	}
	this.Unlock()
	return nodes
}

func (this *Zookeeper) CleanNodes() NetNodeList {
	var nodes NetNodeList
	this.Lock()
	for _, node := range this.group {
		nodes = append(nodes, node)
	}
	this.group = make(map[int]INetNode)
	this.Unlock()
	return nodes
}

func (this *Zookeeper) RemoveNode(idx int) INetNode {
	this.Lock()
	if val, ok := this.group[idx]; ok {
		delete(this.group, idx)
		this.Unlock()
		return val
	}
	this.Unlock()
	return nil
}

func (this *Zookeeper) HasNode(idx int) bool {
	return this.GetNode(idx) != nil
}

func (this *Zookeeper) UsedSessionID() uint64 {
	return atomic.AddUint64(&this.session, 1)
}

func (this *Zookeeper) ResetSession() {
	atomic.SwapUint64(&this.session, 0)
}

//推送消息(非定向推送, 消息的广泛推送)
func (this *Zookeeper) DoneTopic(topic int, data interface{}) int {
	if nodes := this.GetNodeByTopic(topic); len(nodes) > 0 {
		for _, node := range nodes {
			this.Done(node, data)
		}
		return len(nodes)
	}
	fmt.Println(this.Name(), "集群没有协议处理:", topic)
	return -1
}

//推送到节点(定向推送)
func (this *Zookeeper) DoneNode(idx int, data interface{}) int {
	if node := this.GetNode(idx); node != nil {
		this.Done(node, data)
		return 0
	}
	fmt.Println(this.Name(), "集群没有节点:", idx)
	return -1
}

func (this *Zookeeper) SetHandle(val NetNodeBlock) {
	this.handle = val
}

//做事
func (this *Zookeeper) Done(node INetNode, data interface{}) {
	this.handle(node, data)
}
