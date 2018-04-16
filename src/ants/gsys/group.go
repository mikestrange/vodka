package gsys

/*
线程锁
*/
import "sync"

//interface
type IWaitGroup interface {
	Add() IWaitGroup
	Done() IWaitGroup
	Wait() IWaitGroup
	Wrap(func()) IWaitGroup
}

//信号组
func Group() IWaitGroup {
	return new(WaitGroup)
}

//异步等待
func Wrap(block func()) IWaitGroup {
	return Group().Wrap(block).Wait()
}

func WrapList(block func(), size int) IWaitGroup {
	wg := Group()
	for i := 0; i < size; i++ {
		wg.Wrap(block)
	}
	return wg.Wait()
}

func Wraps(args ...func()) IWaitGroup {
	wg := Group()
	for i := range args {
		wg.Wrap(args[i])
	}
	return wg.Wait()
}

//class Locked
type WaitGroup struct {
	wg sync.WaitGroup
}

//protected
func (this *WaitGroup) Add() IWaitGroup {
	this.wg.Add(1)
	return this
}

func (this *WaitGroup) Done() IWaitGroup {
	this.wg.Done()
	return this
}

//public
func (this *WaitGroup) Wait() IWaitGroup {
	this.wg.Wait()
	return this
}

//异步等待方式
func (this *WaitGroup) Wrap(block func()) IWaitGroup {
	this.Add()
	go func() {
		block()
		this.Done()
	}()
	return this
}
