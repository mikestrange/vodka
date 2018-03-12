package gsys

/*
线程锁
*/
import "sync"

//interface
type ILocked interface {
	Lock()
	Unlock()
}

//class Locked
type Locked struct {
	ILocked
	mutex *sync.Mutex
}

func NewLocked() ILocked {
	this := new(Locked)
	this.InitLocked()
	return this
}

func (this *Locked) InitLocked() {
	this.mutex = new(sync.Mutex)
}

func (this *Locked) Lock() {
	this.mutex.Lock()
}

func (this *Locked) Unlock() {
	this.mutex.Unlock()
}
