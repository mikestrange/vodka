package gutil

import (
	"fmt"
)

//回调
type ModeBlock func(interface{}, ...interface{})

//模块控制
type IModeAccessor interface {
	SetHandle(ModeBlock)
	On(int, interface{}) interface{}
	Off(int) interface{}
	Get(cmd int) interface{}
	Done(int, ...interface{}) bool
}

type ModeAccessor struct {
	handle   ModeBlock
	commands map[int]interface{}
}

func NewMode() IModeAccessor {
	this := new(ModeAccessor)
	this.InitModeAccessor()
	return this
}

func NewModeWithHandle(val ModeBlock) IModeAccessor {
	this := NewMode()
	this.SetHandle(val)
	return this
}

func (this *ModeAccessor) InitModeAccessor() {
	this.commands = make(map[int]interface{})
}

func (this *ModeAccessor) SetHandle(val ModeBlock) {
	this.handle = val
}

func (this *ModeAccessor) On(cmd int, block interface{}) interface{} {
	val, ok := this.commands[cmd]
	this.commands[cmd] = block
	if ok {
		return val
	}
	return nil
}

func (this *ModeAccessor) Off(cmd int) interface{} {
	val, ok := this.commands[cmd]
	if ok {
		delete(this.commands, cmd)
		return val
	}
	return nil
}

func (this *ModeAccessor) Get(cmd int) interface{} {
	fun, ok := this.commands[cmd]
	if ok {
		return fun
	}
	return nil
}

func (this *ModeAccessor) Done(cmd int, args ...interface{}) bool {
	if block := this.Get(cmd); block != nil {
		if this.handle != nil {
			this.handle(block, args...)
			return true
		} else {
			fmt.Println("No Mode Block Err:", cmd)
		}
	} else {
		fmt.Println("No cmd handle:", cmd)
	}
	return false
}
