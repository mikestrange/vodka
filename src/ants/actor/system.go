package actor

import "ants/gsys"

type ActorSet map[int]IActorRef
type ActorList []IActorRef
type ActorRefFunc func(int, IActor) IActorRef

//主要的分布式系统(为分布式提供基础)
var Main IActorSystem = NewSystem()

//节点
type IActorNode interface {
	ActorOf(int, IActor) (IActorRef, bool)
	Remove(int) (IActor, bool)
	Send(int, ...interface{}) bool
	Broadcast(...interface{})
	Shutdown()
	//工厂
	SetRefHandle(ActorRefFunc)
}

//系统
type IActorSystem interface {
	IActorNode
	IActorRef
}

type ActorSystem struct {
	gsys.Locked
	//自己监听
	ActorRef
	//等待结束
	gsys.WaitGroup
	//生成器
	handle ActorRefFunc
	//列表
	actors ActorSet
}

func NewSystem() IActorSystem {
	this := new(ActorSystem)
	this.Init()
	return this
}

//继承必须实现
func (this *ActorSystem) Init() {
	this.actors = make(ActorSet)
	//默认处理
	this.handle = func(idx int, ator IActor) IActorRef {
		ref := NewRef()
		ref.SetActor(ator)
		return ref
	}
}

func (this *ActorSystem) SetRefHandle(hanle ActorRefFunc) {
	this.handle = hanle
}

func (this *ActorSystem) ActorOf(idx int, ator IActor) (IActorRef, bool) {
	ref := this.handle(idx, ator)
	this.Lock()
	defer this.Unlock()
	//失败后不公开
	_, ok := this.actors[idx]
	if ok {
		return ref, false
	}
	this.actors[idx] = ref
	//处理机智
	this.HandleActor(idx, ator, ref)
	return ref, true
}

func (this *ActorSystem) Remove(aid int) (IActor, bool) {
	this.Lock()
	val, ok := this.actors[aid]
	if ok {
		delete(this.actors, aid)
		this.Unlock()
		val.Close()
		return val.ActorObj(), true
	}
	this.Unlock()
	return nil, false
}

func (this *ActorSystem) actorRef(aid int) (IActorRef, bool) {
	this.Lock()
	defer this.Unlock()
	if val, ok := this.actors[aid]; ok {
		return val, true
	}
	return nil, false
}

func (this *ActorSystem) Send(aid int, args ...interface{}) bool {
	if val, ok := this.actorRef(aid); ok {
		return this.callActor(val, args)
	}
	return false
}

func (this *ActorSystem) Broadcast(args ...interface{}) {
	arr := this.getActors(false)
	//调度所有
	for i := range arr {
		this.callActor(arr[i], args)
	}
}

func (this *ActorSystem) getActors(clean bool) ActorList {
	var arr ActorList
	this.Lock()
	for _, val := range this.actors {
		arr = append(arr, val)
	}
	if clean {
		this.actors = make(ActorSet)
	}
	this.Unlock()
	return arr
}

func (this *ActorSystem) callActor(ref IActorRef, datas []interface{}) bool {
	return ref.Router(datas...)
}

func (this *ActorSystem) Shutdown() {
	arr := this.getActors(true)
	//关闭所有
	for i := range arr {
		arr[i].Close()
	}
	this.Wait()
}

//处理(子类可以处理)
func (this *ActorSystem) HandleActor(idx int, ator IActor, ref IActorRef) {
	this.Add()
	ator.OnReady(ref)
	go func() {
		//运行
		ref.Run()
		//正常删除
		this.remActor(idx, ator)
		//最后关闭
		ator.OnClose()
		//最后的温柔
		this.Done()
	}()
}

//匹配删除
func (this *ActorSystem) remActor(idx int, ator IActor) bool {
	this.Lock()
	defer this.Unlock()
	val, ok := this.actors[idx]
	if ok && val.ActorObj() == ator {
		delete(this.actors, idx)
		return true
	}
	return false
}
