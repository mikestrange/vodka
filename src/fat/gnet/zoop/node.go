package gnet

import (
	"sync"
)

type NodeBlock func() (interface{}, bool)

//节点列表
type NetNodeList []INetNode

//网络节点
type INetNode interface {
	//get
	Name() string  //节点名称
	NodeID() int   //节点id = logic_id | logic_type
	Topics() []int //关心的事务
	Addr() string  //服务器地址
	Tag() int
	//set
	SetTag(int)
	SetAddr(string)
	SetName(string)
	SetNodeID(int)
	SetTopics(...int)
	SetLocalData(int, string, string, ...int)
	HasTopic(int) bool
	//代理客户端
	Client(NodeBlock) (interface{}, bool)
	UnClient() interface{}
	//
	//Lock()
	//Unlock()
}

//每一个节点只做一个事务
type NetNode struct {
	nodeid int
	tag    int
	name   string
	addr   string
	topics []int
	mutex  sync.Mutex
	client interface{}
}

func NewNode() INetNode {
	this := new(NetNode)
	return this
}

func NewNodeWithArgs(idx int, name string, addr string, topics ...int) INetNode {
	this := new(NetNode)
	this.SetLocalData(idx, name, addr, topics...)
	return this
}

//gets
func (this *NetNode) Name() string {
	return this.name
}

func (this *NetNode) NodeID() int {
	return this.nodeid
}

func (this *NetNode) Addr() string {
	return this.addr
}

func (this *NetNode) Topics() []int {
	return this.topics
}

//sets
func (this *NetNode) SetAddr(addr string) {
	this.addr = addr
}

func (this *NetNode) SetName(name string) {
	this.name = name
}

func (this *NetNode) SetNodeID(idx int) {
	this.nodeid = idx
}

func (this *NetNode) SetTopics(topics ...int) {
	this.topics = topics
}

func (this *NetNode) SetLocalData(idx int, name string, addr string, topics ...int) {
	this.SetNodeID(idx)
	this.SetName(name)
	this.SetAddr(addr)
	this.SetTopics(topics...)
}

//others
func (this *NetNode) HasTopic(topic int) bool {
	for _, val := range this.topics {
		if val == topic {
			return true
		}
	}
	return false
}

//不要出现在 Client回调中
func (this *NetNode) UnClient() interface{} {
	this.Lock()
	client := this.client
	this.client = nil
	this.Unlock()
	return client
}

//不允许多次设置(注意回调中不要出现UnClient)
func (this *NetNode) Client(block NodeBlock) (interface{}, bool) {
	this.Lock()
	defer this.Unlock()
	if this.client == nil {
		if client, ok := block(); ok {
			this.client = client
			return client, true
		}
		println("设置节点回调失败")
		return nil, false
	}
	return this.client, true
}

func (this *NetNode) SetTag(val int) {
	this.tag = val
}

func (this *NetNode) Tag() int {
	return this.tag
}

func (this *NetNode) Lock() {
	this.mutex.Lock()
}

func (this *NetNode) Unlock() {
	this.mutex.Unlock()
}
