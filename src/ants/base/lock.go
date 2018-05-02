package base

/*
线程锁
*/
import "sync"

//interface
type ILocked interface {
	Lock()
	Unlock()
	Auto(func())
}

func NewLocked() ILocked {
	return new(Locked)
}

//class Locked
type Locked struct {
	_m sync.Mutex
}

func (this *Locked) Lock() {
	this._m.Lock()
}

func (this *Locked) Unlock() {
	this._m.Unlock()
}

//自动执行加锁模式
func (this *Locked) Auto(f func()) {
	this.Lock()
	f()
	this.Unlock()
}
