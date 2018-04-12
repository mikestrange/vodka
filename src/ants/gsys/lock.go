package gsys

/*
线程锁
*/
import "sync"

//interface
type ILocked interface {
	Lock()
	Unlock()
	//自动加锁
	Auto(func())
	AutoBool(func() bool) bool
	AutoInt(func() int) int
	AutoArg(func() interface{}) interface{}
	AutoArgs(func() []interface{}) []interface{}
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

//自动执行加锁模式
func (this *Locked) Auto(f func()) {
	this.Lock()
	f()
	this.Unlock()
}

func (this *Locked) AutoBool(f func() bool) bool {
	this.Lock()
	ok := f()
	this.Unlock()
	return ok
}

func (this *Locked) AutoInt(f func() int) int {
	this.Lock()
	val := f()
	this.Unlock()
	return val
}

func (this *Locked) AutoArg(f func() interface{}) interface{} {
	this.Lock()
	val := f()
	this.Unlock()
	return val
}

func (this *Locked) AutoArgs(f func() []interface{}) []interface{} {
	this.Lock()
	vals := f()
	this.Unlock()
	return vals
}
