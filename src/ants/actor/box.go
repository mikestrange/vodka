package actor

import "ants/gsys"

type RefSet map[interface{}]IBoxRef

//主要
var Main IBoxSystem = new(BoxSystem)

//静态
type BoxSystem struct {
	BaseBox //设置自身吧
	wgRef   gsys.WaitGroup
	locked  gsys.Locked
	refs    RefSet
}

//private
func (this *BoxSystem) lock() {
	this.locked.Lock()
}

func (this *BoxSystem) unlock() {
	this.locked.Unlock()
}

//public
func (this *BoxSystem) ActorOf(mark interface{}, val IBoxRef) bool {
	this.lock()
	ref, ok := this.setRef(mark, val)
	this.unlock()
	if ok {
		this.handleRefer(ref)
		return true
	}
	return false
}

//遵循释放顺序
func (this *BoxSystem) handleRefer(ref IBoxRef) {
	ref.setFather(this) //没有其他太多的权限
	ref.OnReady()
	this.wgRef.Wrap(func() {
		ref.PerformRunning()
		this.UnRef(ref)
		ref.OnRelease()
		ref.setFather(nil)
	})
}

//查找移除
func (this *BoxSystem) UnRef(ref IBoxRef) bool {
	ok := false
	this.lock()
	if this.refs != nil {
		for key, val := range this.refs {
			if val == ref {
				delete(this.refs, key)
				ok = true
				break
			}
		}
	}
	this.unlock()
	return ok
}

func (this *BoxSystem) CloseRef(mark interface{}) (IBoxRef, bool) {
	this.lock()
	ref, ok := this.delRef(mark)
	this.unlock()
	if ok {
		ref.Die()
		return ref, true
	}
	return nil, false
}

func (this *BoxSystem) FindRef(mark interface{}) (IBoxRef, bool) {
	this.lock()
	if this.refs != nil {
		if ref, ok := this.refs[mark]; ok {
			this.unlock()
			return ref, true
		}
	}
	this.unlock()
	return nil, false
}

func (this *BoxSystem) Send(mark interface{}, args ...interface{}) bool {
	if ref, ok := this.FindRef(mark); ok {
		return ref.Router(args...)
	}
	return false
}

func (this *BoxSystem) Broadcast(args ...interface{}) {
	ref_list := this.RefList(false)
	//集体推送
	for i := range ref_list {
		ref_list[i].Router(args...)
	}
}

func (this *BoxSystem) CloseAll() {
	ref_list := this.RefList(true)
	//集体推送
	for i := range ref_list {
		ref_list[i].Die()
	}
	this.wgRef.Wait()
}

//protected
func (this *BoxSystem) RefList(del bool) []IBoxRef {
	var ref_list []IBoxRef
	this.lock()
	if this.refs != nil {
		for _, ref := range this.refs {
			ref_list = append(ref_list, ref)
		}
		if del {
			this.refs = nil
		}
	}
	this.unlock()
	return ref_list
}

func (this *BoxSystem) setRef(mark interface{}, ref IBoxRef) (IBoxRef, bool) {
	if this.refs == nil {
		this.refs = make(RefSet)
		this.refs[mark] = ref
		return ref, true
	}
	if old, ok := this.refs[mark]; ok {
		return old, false
	}
	this.refs[mark] = ref
	return ref, true
}

func (this *BoxSystem) delRef(mark interface{}) (IBoxRef, bool) {
	if this.refs == nil {
		return nil, false
	}
	if ref, ok := this.refs[mark]; ok {
		delete(this.refs, mark)
		return ref, true
	}
	return nil, false
}
