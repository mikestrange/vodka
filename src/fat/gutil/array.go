package gutil

import (
	"sort"
)

//比较合理的数组
const NOT_VALUE = -1

//成等比曲线增长
const MIN_GROW_SIZE = 100

type IArrayObject interface {
	//add
	SetVal(int, interface{}) interface{}
	Push(...interface{})
	UnShift(...interface{})
	Insert(int, ...interface{})
	//del
	Pop() interface{}
	Shift() interface{}
	//seach
	Each(func(int, interface{}))
	IndexOf(interface{}) int
	SeachValue(func(interface{}) bool) (int, interface{})
	IndexOfLast(interface{}) int
	//del
	DelIndex(int) interface{}
	DelVal(interface{}) bool
	DelVals(interface{})
	//copy
	CopyVals() []interface{}
	Elements() []interface{}
	//get
	GetVal(int) interface{}
	BeginVal() interface{}
	EndVal() interface{}
	//others
	Reset()
	Empty() bool
	Size() int
	CapSize() int
	Sort(func(interface{}, interface{}) bool)
}

type Array struct {
	IArrayObject
	m_list []interface{}
}

func NewArray() IArrayObject {
	return NewArrayWithSize(0)
}

func NewArrayWithSize(size int) IArrayObject {
	this := new(Array)
	this.InitArray(size)
	return this
}

func (this *Array) InitArray(size int) {
	this.m_list = make([]interface{}, size)
}

//关键函数
func (this *Array) grow(size int) {
	cap_size := this.CapSize()
	if size > cap_size {
		if cap_size == 0 {
			cap_size = MIN_GROW_SIZE
		}
		cap_size := cap_size * (size/cap_size + 1)
		this.m_list = append(make([]interface{}, 0, cap_size), this.m_list...)
	}
}

func (this *Array) Push(args ...interface{}) {
	this.grow(this.Size() + len(args))
	this.m_list = append(this.m_list, args...)
}

func (this *Array) UnShift(args ...interface{}) {
	sub_size := len(args)
	size := this.Size()
	this.grow(size + sub_size)
	//先追加为了扩充内容
	this.m_list = append(this.m_list, args...)
	//原始部分后移
	copy(this.m_list[sub_size:], this.m_list[0:size])
	//首部添加新元素
	copy(this.m_list[0:], args)
}

func (this *Array) SetVal(pos int, val interface{}) interface{} {
	old := this.m_list[pos]
	this.m_list[pos] = val
	return old
}

func (this *Array) Pop() interface{} {
	if size := this.Size(); size > 0 {
		pos := size - 1
		val := this.m_list[pos]
		this.m_list = this.m_list[0:pos]
		return val
	}
	return nil
}

func (this *Array) Shift() interface{} {
	if size := this.Size(); size > 0 {
		val := this.m_list[0]
		this.m_list = this.m_list[1:size]
		return val
	}
	return nil
}

func (this *Array) GetVal(idx int) interface{} {
	return this.m_list[idx]
}

func (this *Array) BeginVal() interface{} {
	return this.m_list[0]
}

func (this *Array) EndVal() interface{} {
	return this.m_list[this.Size()-1]
}

func (this *Array) Insert(pos int, args ...interface{}) {
	if pos >= this.Size() {
		this.Push(args...)
	} else if pos <= 0 {
		this.UnShift(args...)
	} else {
		this.grow(this.Size() + len(args))
		//先追加为了扩充内容
		this.m_list = append(this.m_list, args...)
		//再把中间的拷贝到后面
		copy(this.m_list[(pos+len(args)):], this.m_list[pos:])
		//最后新的部分拷贝进来
		copy(this.m_list[pos:], args)
	}
}

func (this *Array) IndexOf(val interface{}) int {
	for i := range this.m_list {
		if val == this.m_list[i] {
			return i
		}
	}
	return NOT_VALUE
}

func (this *Array) Each(block func(int, interface{})) {
	for i := range this.m_list {
		block(i, this.m_list[i])
	}
}

func (this *Array) SeachValue(block func(interface{}) bool) (int, interface{}) {
	for i := range this.m_list {
		if block(this.m_list[i]) {
			return i, this.m_list[i]
		}
	}
	return NOT_VALUE, nil
}

func (this *Array) IndexOfLast(val interface{}) int {
	for i := range this.m_list {
		if this.m_list[i] == val {
			return i
		}
	}
	return NOT_VALUE
}

func (this *Array) DelIndex(i int) interface{} {
	if i >= 0 && i < this.Size() {
		val := this.m_list[i]
		this.m_list = append(this.m_list[:i], this.m_list[i+1:]...)
		return val
	}
	return nil
}

func (this *Array) DelVal(val interface{}) bool {
	return this.DelIndex(this.IndexOf(val)) != nil
}

func (this *Array) DelVals(val interface{}) {
	for this.DelVal(val) {
		//清理相同的值
	}
}

func (this *Array) CopyVals() []interface{} {
	return append([]interface{}{}, this.m_list...)
}

func (this *Array) Elements() []interface{} {
	return this.m_list
}

func (this *Array) Reset() {
	this.m_list = this.m_list[:0]
}

func (this *Array) Empty() bool {
	return this.Size() == 0
}

func (this *Array) Size() int {
	return len(this.m_list)
}

func (this *Array) CapSize() int {
	return cap(this.m_list)
}

func (this *Array) Sort(block func(interface{}, interface{}) bool) {
	SortHandle(this.m_list, block)
}

/*
基于系统一种简单的排序
*/
func SortHandle(data []interface{}, block func(interface{}, interface{}) bool) {
	sort.Sort(&SortInterface{data, block})
}

type SortInterface struct {
	//sort.Interface
	SortList   []interface{}
	SortHandle func(interface{}, interface{}) bool
}

func (this *SortInterface) Len() int {
	return len(this.SortList)
}

func (this *SortInterface) Less(i, j int) bool {
	return this.SortHandle(this.SortList[i], this.SortList[j])
}

func (this *SortInterface) Swap(i, j int) {
	this.SortList[i], this.SortList[j] = this.SortList[j], this.SortList[i]
	//temp := this.SortList[i]
	//this.SortList[i] = this.SortList[j]
	//this.SortList[j] = temp
}
