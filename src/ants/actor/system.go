package actor

import "sync"

type ActorSet map[int]IActor
type ActorList []IActor

//主要的分布式系统(为分布式提供基础)
var Main IActorSystem = NewSystem()

type IActorSystem interface {
	Shutdown()
	Added(IActor) (IActor, bool)
	Remove(int) (IActor, bool)
	Actor(int) (IActor, bool)
	Send(int, ...interface{}) bool
	Broadcast(...interface{})
}

type ActorSystem struct {
	actors  ActorSet
	mutex   sync.Mutex
	wgAtors sync.WaitGroup
}

func NewSystem() IActorSystem {
	this := new(ActorSystem)
	this.init()
	return this
}

func (this *ActorSystem) init() {
	this.actors = make(ActorSet)
}

func (this *ActorSystem) put(ator IActor) (IActor, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	_, ok := this.actors[ator.Idx()]
	if ok {
		return ator, false
	}
	this.actors[ator.Idx()] = ator
	return ator, true
}

func (this *ActorSystem) remove(aid int) (IActor, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	val, ok := this.actors[aid]
	if ok {
		delete(this.actors, aid)
		return val, true
	}
	return nil, false
}

func (this *ActorSystem) getActor(aid int) (IActor, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if val, ok := this.actors[aid]; ok {
		return val, true
	}
	return nil, false
}

func (this *ActorSystem) getActors() ActorList {
	var arr ActorList
	this.mutex.Lock()
	for _, val := range this.actors {
		arr = append(arr, val)
	}
	this.mutex.Unlock()
	return arr
}

func (this *ActorSystem) removeAll() {
	var arr ActorList
	this.mutex.Lock()
	for _, val := range this.actors {
		arr = append(arr, val)
	}
	this.actors = make(ActorSet)
	this.mutex.Unlock()
	//关闭所有
	for i := range arr {
		arr[i].Close()
	}
}

//完全匹配移除
func (this *ActorSystem) removeAtor(ator IActor) bool {
	this.mutex.Lock()
	val, ok := this.actors[ator.Idx()]
	if ok && val == ator {
		delete(this.actors, ator.Idx())
		this.mutex.Unlock()
		ator.Close()
		return true
	}
	this.mutex.Unlock()
	return false
}

//interfaces
func (this *ActorSystem) Added(ator IActor) (IActor, bool) {
	val, ok := this.put(ator)
	if ok {
		this.handleActor(ator)
		return ator, true
	}
	println("add err")
	return val, false
}

func (this *ActorSystem) Remove(aid int) (IActor, bool) {
	val, ok := this.remove(aid)
	if ok {
		val.Close()
		return val, true
	}
	return nil, false
}

func (this *ActorSystem) Actor(aid int) (IActor, bool) {
	return this.getActor(aid)
}

func (this *ActorSystem) Send(aid int, args ...interface{}) bool {
	if val, ok := this.getActor(aid); ok {
		return this.callActor(val, args)
	}
	return false
}

func (this *ActorSystem) Broadcast(args ...interface{}) {
	arr := this.getActors()
	//调度所有
	for i := range arr {
		this.callActor(arr[i], args)
	}
}

func (this *ActorSystem) callActor(ator IActor, datas []interface{}) bool {
	//这里可以交给路由处理下(然后路由转发)
	return ator.Commit(datas...)
}

func (this *ActorSystem) Shutdown() {
	this.removeAll()
	this.wgAtors.Wait()
}

//处理
func (this *ActorSystem) handleActor(ator IActor) {
	this.wgAtors.Add(1)
	go func() {
		//自身运载器运载自身
		ator.Runner().LoopActor(ator)
		//避免非正常删除(这里有点危险)
		this.removeAtor(ator)
		//最后被关闭
		ator.OnClosed()
		this.wgAtors.Done()
	}()
}
