package gutil

import (
	"fmt"
)

//事件回调
type EventBlock func(interface{}, int, ...interface{})

/*
消息注册器
*/
type IEventDispatcher interface {
	SetHandle(EventBlock)
	On(event int, block interface{})
	Off(event int, block interface{})
	OnTo(target interface{}, events []int)
	OffTo(block interface{})
	Has(event int, block interface{}) bool
	Done(event int, data ...interface{})
}

/*
基础事件对象(没有锁)
*/
type EventDispatcher struct {
	events map[int]IArrayObject
	handle EventBlock
}

func NewEvent() IEventDispatcher {
	this := new(EventDispatcher)
	this.InitEventDispatcher()
	return this
}

func NewEventWithHandle(val EventBlock) IEventDispatcher {
	this := NewEvent()
	this.SetHandle(val)
	return this
}

func (this *EventDispatcher) InitEventDispatcher() {
	this.events = make(map[int]IArrayObject)
}

func (this *EventDispatcher) SetHandle(val EventBlock) {
	this.handle = val
}

func (this *EventDispatcher) On(event int, block interface{}) {
	if _, ok := this.events[event]; !ok {
		this.events[event] = NewArray()
	}
	this.events[event].Push(block)
}

func (this *EventDispatcher) OnTo(target interface{}, events []int) {
	for i := range events {
		this.On(events[i], target)
	}
}

func (this *EventDispatcher) OffTo(target interface{}) {
	for k, list := range this.events {
		list.DelVals(target)
		if list.Empty() {
			delete(this.events, k)
		}
	}
}

func (this *EventDispatcher) getListeners(event int) []interface{} {
	if list, ok := this.events[event]; ok {
		return list.CopyVals()
	}
	return nil
}

func (this *EventDispatcher) Off(event int, block interface{}) {
	if list, ok := this.events[event]; ok {
		if list.DelVal(block); list.Empty() {
			delete(this.events, event)
		}
	}
}

func (this *EventDispatcher) Has(event int, block interface{}) bool {
	if list, ok := this.events[event]; ok {
		return NOT_VALUE != list.IndexOf(block)
	}
	return false
}

func (this *EventDispatcher) Done(event int, args ...interface{}) {
	if blocks := this.getListeners(event); blocks != nil {
		for _, block := range blocks {
			if this.handle != nil {
				this.handle(block, event, args...)
			} else {
				fmt.Println("Event no Block Error")
			}
		}
	} else {
		fmt.Println("Event no handle:", event)
	}
}
