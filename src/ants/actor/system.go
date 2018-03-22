package actor

import "sync"

type ActorSet map[int]IActor
type ActorList []IActor

type IActorSystem interface {
	Shutdown()
	Added(IActor) (IActor, bool)
	Remove(int) (IActor, bool)
	Send(int, interface{}) bool
	Broadcast(interface{})
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
	val, ok := this.actors[ator.Idx()]
	if ok {
		return val, false
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

//interfaces
func (this *ActorSystem) Added(ator IActor) (IActor, bool) {
	val, ok := this.put(ator)
	if ok {
		this.handleActor(ator)
		return ator, true
	}
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

func (this *ActorSystem) Send(aid int, data interface{}) bool {
	if val, ok := this.getActor(aid); ok {
		return this.callActor(val, data)
	}
	return false
}

func (this *ActorSystem) Broadcast(data interface{}) {
	arr := this.getActors()
	//调度所有
	for i := range arr {
		this.callActor(arr[i], data)
	}
}

func (this *ActorSystem) callActor(ator IActor, data interface{}) bool {
	if ator.WillReceive(data) {
		return ator.Commit(data)
	}
	return false
}

func (this *ActorSystem) Shutdown() {
	this.removeAll()
	this.wgAtors.Wait()
}

//处理
func (this *ActorSystem) handleActor(ator IActor) {
	this.wgAtors.Add(1)
	go func() {
		ator.Runner().LoopActor(ator)
		this.Remove(ator.Idx())
		ator.OnClosed() //被关闭
		this.wgAtors.Done()
	}()
}
