package gwork

import (
	"ants/gsys"
	"sync"
)

//相关
type workChan chan []interface{}

//工作选择器(可以并行，也可以单线运行，默认单线)
type WorkSelector struct {
	L        sync.Mutex
	openFlag bool
	size     int
	step     int
	mqLen    int
	chans    []workChan
}

func NewWorks(num int, size int) *WorkSelector {
	this := new(WorkSelector)
	this.SetThreadNum(num)
	this.SetMqNum(size)
	return this
}

//线程数量，队列长度
func (this *WorkSelector) SetMqNum(val int) {
	this.L.Lock()
	if !this.openFlag {
		this.mqLen = val
	} else {
		println("do not set by open")
	}
	this.L.Unlock()
}

func (this *WorkSelector) SetThreadNum(val int) {
	this.L.Lock()
	if !this.openFlag {
		this.size = val
	} else {
		println("do not set by open")
	}
	this.L.Unlock()
}

//interfaces
func (this *WorkSelector) Push(args ...interface{}) bool {
	this.L.Lock()
	defer this.L.Unlock()
	if this.openFlag {
		if this.size == THREAD_DEF_SIZE {
			this.chans[0] <- args
		} else {
			this.chans[this.step%this.size] <- args
			this.step++
		}
		return true
	}
	return false
}

func (this *WorkSelector) Pull() ([]interface{}, bool) {
	return nil, false
}

func (this *WorkSelector) Exit() bool {
	this.L.Lock()
	if this.openFlag {
		this.openFlag = false
		for i := 0; i < this.size; i++ {
			close(this.chans[i])
		}
	}
	this.L.Unlock()
	return true
}

func (this *WorkSelector) Open() bool {
	this.L.Lock()
	defer this.L.Unlock()
	if this.openFlag {
		return false
	}
	this.step = 0
	this.openFlag = true
	if this.size < THREAD_DEF_SIZE {
		this.size = THREAD_DEF_SIZE
	}
	if this.mqLen < QUEUE_MIN_SIZE {
		this.mqLen = QUEUE_DEF_SIZE
	}
	this.chans = make([]workChan, this.size)
	for i := 0; i < this.size; i++ {
		this.chans[i] = make(workChan, this.mqLen)
	}
	return true
}

//等待所有结束
func (this *WorkSelector) Join(block func(...interface{})) {
	wg := gsys.NewWgGroup()
	for i := 0; i < this.size; i++ {
		wg.Add()
		go func(work workChan) {
			for {
				args, ok := <-work
				if ok {
					Handle(block, args)
				} else {
					break
				}
			}
			wg.Done()
		}(this.chans[i])
	}
	wg.Wait()
}
