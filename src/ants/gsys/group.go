package gsys

/*
线程锁
*/
import "sync"

//interface
type IWaitGroup interface {
	Add()
	Done()
	Wait()
}

func NewWgGroup() IWaitGroup {
	return new(WaitGroup)
}

//class Locked
type WaitGroup struct {
	wg sync.WaitGroup
}

func (this *WaitGroup) Add() {
	this.wg.Add(1)
}

func (this *WaitGroup) Done() {
	this.wg.Done()
}

func (this *WaitGroup) Wait() {
	this.wg.Wait()
}
