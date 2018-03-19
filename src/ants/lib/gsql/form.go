package gsql

//表单列表
type IRowsForm interface {
	//私有
	next() IRowItem
	//公开
	Item(int) IRowItem
	Length() int
	Empty() bool
	Elements() []IRowItem
}

//每一行数据
type IRowItem interface {
	//私有
	set(string, interface{})
	//公开
	Has(string) bool
	Value(string) (interface{}, bool)
	Int(string) int
	Str(string) string
	Int64(string) int64
	Bool(string) bool
	Elements() map[string]interface{}
	//CompareVal(string, interface{}) bool
}

//###########class item
func newItem() IRowItem {
	return &rowItem{items: make(map[string]interface{})}
}

type rowItem struct {
	items map[string]interface{}
}

func (this *rowItem) set(k string, v interface{}) {
	this.items[k] = v
}

func (this *rowItem) Elements() map[string]interface{} {
	return this.items
}

func (this *rowItem) Has(k string) bool {
	if _, ok := this.items[k]; ok {
		return true
	}
	return false
}

func (this *rowItem) Value(k string) (interface{}, bool) {
	if v, ok := this.items[k]; ok {
		return v, true
	}
	return nil, false
}

func (this *rowItem) Bool(k string) bool {
	return this.Int(k) > 0
}

func (this *rowItem) Int(k string) int {
	if v, ok := this.items[k]; ok {
		return Int(v)
	}
	return 0
}

func (this *rowItem) Int64(k string) int64 {
	if v, ok := this.items[k]; ok {
		return Int64(v)
	}
	return 0
}

func (this *rowItem) Str(k string) string {
	if v, ok := this.items[k]; ok {
		return Str(v)
	}
	return ""
}

//###########class form
func newForm() IRowsForm {
	return &rowsForm{}
}

type rowsForm struct {
	rows []IRowItem
}

func (this *rowsForm) next() IRowItem {
	item := newItem()
	this.rows = append(this.rows, item)
	return item
}

func (this *rowsForm) Elements() []IRowItem {
	return this.rows
}

func (this *rowsForm) Length() int {
	return len(this.rows)
}

func (this *rowsForm) Item(idx int) IRowItem {
	return this.rows[idx]
}

func (this *rowsForm) Empty() bool {
	return len(this.rows) == 0
}
