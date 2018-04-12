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
	Wrap(func())
	WrapArg(func(interface{}), interface{})
	WrapArgs(func(...interface{}), ...interface{})
}

func NewGroup() IWaitGroup {
	return new(WaitGroup)
}

func NewGroupWrap(f func()) IWaitGroup {
	wg := NewGroup()
	wg.Wrap(f)
	wg.Wait()
	return wg
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

//异步等待方式
func (this *WaitGroup) Wrap(block func()) {
	this.Add()
	go func() {
		block()
		this.Done()
	}()
}

func (this *WaitGroup) WrapArg(block func(interface{}), data interface{}) {
	this.Add()
	go func() {
		block(data)
		this.Done()
	}()
}

func (this *WaitGroup) WrapArgs(block func(...interface{}), args ...interface{}) {
	this.Add()
	go func() {
		block(args...)
		this.Done()
	}()
}
