package base

/*
异步线程锁
*/
import "sync"

//为系统提供
var sys_group WaitGroup

func SysWrap(args ...func()) {
	sys_group.Wraps(args...)
}

func SysWait() {
	sys_group.Wait()
}

//interface
type IGroup interface {
	Add() IGroup
	Done() IGroup
	Wait() IGroup
	Wrap(func()) IGroup
	Wraps(...func()) IGroup
	WrapList(func(), int) IGroup
}

//信号组
func Group() IGroup {
	return new(WaitGroup)
}

//异步等待
func Wrap(f func()) IGroup {
	return Group().Wrap(f).Wait()
}

//多并行等待
func Wraps(args ...func()) IGroup {
	return Group().Wraps(args...).Wait()
}

//同并发等待
func WrapList(f func(), size int) IGroup {
	return Group().WrapList(f, size).Wait()
}

//class Locked
type WaitGroup struct {
	//IGroup
	wg sync.WaitGroup
}

//protected
func (this *WaitGroup) Add() IGroup {
	this.wg.Add(1)
	return this
}

func (this *WaitGroup) Done() IGroup {
	this.wg.Done()
	return this
}

//public
func (this *WaitGroup) Wait() IGroup {
	this.wg.Wait()
	return this
}

//异步等待方式
func (this *WaitGroup) Wrap(f func()) IGroup {
	this.Add()
	TryGo(f, func(ok bool) {
		this.Done()
	})
	return this
}

//异步等待方式
func (this *WaitGroup) Wraps(args ...func()) IGroup {
	for i := range args {
		this.Wrap(args[i])
	}
	return this
}

//异步等待方式
func (this *WaitGroup) WrapList(f func(), size int) IGroup {
	for i := 0; i < size; i++ {
		this.Wrap(f)
	}
	return this
}
