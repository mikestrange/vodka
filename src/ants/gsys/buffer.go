package gsys

//缓冲通道(可以用于发送)
import (
	"sync"
)

//private
type buffItem struct {
	args []interface{}
	next *buffItem
}

//消息队列
type WorkBuff struct {
	size     int
	openFlag bool
	//其他
	current int
	bitem   *buffItem
	eitem   *buffItem
	cond    *sync.Cond
}

func NewBuffer(size int) *WorkBuff {
	return new(WorkBuff).init(size)
}

func (this *WorkBuff) init(sz int) *WorkBuff {
	this.SetSize(sz)
	return this
}

func (this *WorkBuff) SetSize(sz int) {
	this.size = sz
	if this.cond == nil {
		this.cond = sync.NewCond(new(sync.Mutex))
	}
}

func (this *WorkBuff) Open() bool {
	this.cond.L.Lock()
	if !this.openFlag {
		this.openFlag = true
		this.current = 0
		this.cond.L.Unlock()
		return true
	}
	this.cond.L.Unlock()
	return false
}

func (this *WorkBuff) Exit() bool {
	this.cond.L.Lock()
	if this.openFlag {
		this.openFlag = false
		this.cond.Broadcast()
	}
	this.cond.L.Unlock()
	return true
}

func (this *WorkBuff) Push(args ...interface{}) bool {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	//overfull
	if this.current >= this.size {
		return false
	}
	if this.openFlag {
		val := &buffItem{args, nil}
		if this.current == 0 {
			this.bitem, this.eitem = val, val
		} else {
			this.eitem.next, this.eitem = val, val
		}
		this.current++
		this.cond.Signal()
		return true
	}
	return false
}

func (this *WorkBuff) Pull() ([]interface{}, bool) {
	for {
		this.cond.L.Lock()
		if this.current > 0 {
			val := this.bitem.args
			this.bitem = this.bitem.next
			this.current--
			this.cond.L.Unlock()
			return val, true
		}
		if !this.openFlag {
			this.cond.L.Unlock()
			return nil, false
		}
		this.cond.Wait()
		this.cond.L.Unlock()
	}
}

func (this *WorkBuff) Join(block func(...interface{})) {
	for {
		if args, ok := this.Pull(); ok {
			block(args...)
		} else {
			break
		}
	}
}
