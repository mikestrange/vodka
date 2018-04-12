package gutil

import "sort"
import "reflect"

//快速的排序
func Sort(list interface{}, less func(int, int) bool) {
	sort.Sort(new(SortInterface).init(list, less))
}

//interfaces
type SortInterface struct {
	list   reflect.Value
	handle func(int, int) bool
}

func (this *SortInterface) init(data interface{}, less func(int, int) bool) *SortInterface {
	this.handle = less
	this.list = reflect.ValueOf(data)
	return this
}

func (this *SortInterface) Len() int {
	return this.list.Len()
}

func (this *SortInterface) Less(i, j int) bool {
	return this.handle(i, j)
}

func (this *SortInterface) Swap(i, j int) {
	temp := this.list.Index(i).Interface()
	this.list.Index(i).Set(this.list.Index(j))
	this.list.Index(j).Set(reflect.ValueOf(temp))
}
