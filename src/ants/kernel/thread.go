package kernel

import (
	"sync"
)

//控制所有的go
var wg sync.WaitGroup

func add() {
	wg.Add(1)
}

func done() {
	wg.Done()
}

func Wait() {
	wg.Wait()
}

//接口
type IGo interface {
	Run(fun interface{}, args ...interface{})
}

//可以被继承
type GoThread struct {
	Client interface{}
	OnDie  func(interface{}, interface{})
}

//新的线程
func NewGo(die func(interface{}, interface{}), client interface{}) IGo {
	return &GoThread{client, die}
}

func (this *GoThread) Run(fun interface{}, args ...interface{}) {
	add()
	go func() {
		//处理错误
		defer func() {
			if err := recover(); err == nil {
				this.OnDie(this.Client, nil)
			} else {
				this.OnDie(this.Client, err)
			}
			done()
		}()
		//执行
		threadHandle(fun, args)
	}()
}

//支持3个参数
func threadHandle(fun interface{}, args []interface{}) {
	switch f := fun.(type) {
	case func():
		f()
	case func(interface{}):
		f(args[0])
	case func(interface{}, interface{}):
		f(args[0], args[1])
	case func(interface{}, interface{}, interface{}):
		f(args[0], args[1], args[2])
	case func(...interface{}):
		f(args...)
	default:
		panic("not handle func")
	}
}
