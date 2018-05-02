package base

import "sort"
import "reflect"

//快速的排序
func Sort(list interface{}, less func(int, int) bool) {
	sort.Sort(new(sortObj).init(list, less))
}

//sort.Interface
type sortObj struct {
	list   reflect.Value
	handle func(int, int) bool
}

func (this *sortObj) init(data interface{}, less func(int, int) bool) *sortObj {
	this.handle = less
	this.list = reflect.ValueOf(data)
	return this
}

func (this *sortObj) Len() int {
	return this.list.Len()
}

func (this *sortObj) Less(i, j int) bool {
	return this.handle(i, j)
}

func (this *sortObj) Swap(i, j int) {
	temp := this.list.Index(i).Interface()
	this.list.Index(i).Set(this.list.Index(j))
	this.list.Index(j).Set(reflect.ValueOf(temp))
}
