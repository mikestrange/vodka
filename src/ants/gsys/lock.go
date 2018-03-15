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

func NewLocked() ILocked {
	return new(Locked)
}

//class Locked
type Locked struct {
	mutex sync.Mutex
}

func (this *Locked) Lock() {
	this.mutex.Lock()
}

func (this *Locked) Unlock() {
	this.mutex.Unlock()
}
