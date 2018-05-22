package hippo

import "ants/base"

type IBoxRef interface {
	Run()
	OnDie()
	Die() bool
	Push(interface{}) bool
}

type Thread interface {
	Done(interface{}) bool             //自己线程做事
	Die() bool                         //自杀
	Ref() IBoxRef                      //工作运行
	Super() Thread                     //父项
	Main() Thread                      //顶级
	Child(int, IBoxRef) (Thread, bool) //添加子项
	Close(int) (Thread, bool)          //关闭子项
	Find(int) (Thread, bool)           //寻找子项
	Send(int, interface{}) bool        //通知子项
	Broadcast(interface{})             //子项集群通知
}

type NodeThread struct {
	ref      IBoxRef
	super    *NodeThread
	wg       base.WaitGroup
	wgChild  base.WaitGroup
	childSet map[int]Thread
}

func (this *NodeThread) Child(idx int, ref IBoxRef) (Thread, bool) {
	if child, ok := this.newChild(idx); ok {
		child.setRef(ref)
		child.setSuper(this)
		this.wgChild.Wrap(func() {
			child.wg.Wrap(func() {
				ref.Run()
			})
			child.Wait()
			this.unChild(idx, child)
			ref.OnDie()
		})
		return child, true
	}
	return nil, false
}

func (this *NodeThread) Close(idx int) (Thread, bool) {
	if tx, ok := this.Find(idx); ok {
		return tx, tx.Die()
	}
	return nil, false
}

func (this *NodeThread) Send(idx int, event interface{}) bool {
	if tx, ok := this.Find(idx); ok {
		return tx.Done(event)
	}
	return false
}

func (this *NodeThread) Broadcast(event interface{}) {
	list := this.boxList(false)
	for i := range list {
		list[i].Done(event)
	}
}

func (this *NodeThread) Find(idx int) (Thread, bool) {
	if this.childSet != nil {
		if tx, ok := this.childSet[idx]; ok {
			return tx, true
		}
	}
	return nil, false
}

func (this *NodeThread) Main() Thread {
	return &stage
}

func (this *NodeThread) Super() Thread {
	return this.super
}

func (this *NodeThread) setSuper(val *NodeThread) {
	this.super = val
}

func (this *NodeThread) Ref() IBoxRef {
	return this.ref
}

func (this *NodeThread) setRef(ref IBoxRef) {
	this.ref = ref
}

func (this *NodeThread) Done(event interface{}) bool {
	return this.ref.Push(event)
}

func (this *NodeThread) Die() bool {
	return this.ref.Die()
}

//private
func (this *NodeThread) newChild(idx int) (*NodeThread, bool) {
	if this.childSet == nil {
		this.childSet = make(map[int]Thread)
	}
	if _, ok := this.childSet[idx]; ok {
		return nil, false
	}
	child := new(NodeThread)
	this.childSet[idx] = child
	return child, true
}

func (this *NodeThread) unChild(idx int, box Thread) {
	old, ok := this.childSet[idx]
	if ok && old == box {
		delete(this.childSet, idx)
	}
}

func (this *NodeThread) boxList(del bool) []Thread {
	var list []Thread
	//this.Lock()
	if this.childSet != nil {
		for _, tx := range this.childSet {
			list = append(list, tx)
		}
		if del {
			this.childSet = nil
		}
	}
	//this.Unlock()
	return list
}

func (this *NodeThread) Wait() {
	this.wg.Wait()
	this.CloseAll()
	this.wgChild.Wait()
}

func (this *NodeThread) CloseAll() {
	list := this.boxList(true)
	for i := range list {
		list[i].Die()
	}
}

//顶级
var stage TopThread

type TopThread struct {
	NodeThread
	running bool
}

func (this *TopThread) RunSelf(ref IBoxRef) bool {
	if !this.running {
		this.running = true
		this.setRef(ref)
		//异步启动
		this.wg.Wrap(func() {
			ref.Run()
		})
		this.Wait()
		ref.OnDie()
		return true
	}
	return false
}

func init() {

}
