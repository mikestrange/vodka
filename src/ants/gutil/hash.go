package gutil

/*
一般情况不用，泛表可以用
*/
type IHashMap interface {
	Val(interface{}) interface{}
	Get(interface{}) interface{}
	Set(interface{}, interface{}) interface{}
	Has(interface{}) bool
	Del(interface{}) interface{}
	Seach(func(interface{}) bool) interface{}
	Vals() []interface{}
	Keys() []interface{}
	Each(func(interface{}, interface{}))
	Clean() []interface{}
}

type HashMap struct {
	m_map map[interface{}]interface{}
}

func NewHashMap() IHashMap {
	this := new(HashMap)
	this.InitHashMap()
	return this
}

func (this *HashMap) InitHashMap() {
	this.m_map = make(map[interface{}]interface{})
}

func (this *HashMap) Get(k interface{}) interface{} {
	if val, ok := this.m_map[k]; ok {
		return val
	}
	return nil
}

func (this *HashMap) Val(k interface{}) interface{} {
	if val, ok := this.m_map[k]; ok {
		return val
	}
	return nil
}

func (this *HashMap) Set(k interface{}, v interface{}) interface{} {
	old := this.Val(k)
	this.m_map[k] = v
	return old
}

func (this *HashMap) Has(k interface{}) bool {
	if _, ok := this.m_map[k]; ok {
		return true
	}
	return false
}

func (this *HashMap) Del(k interface{}) interface{} {
	if val, ok := this.m_map[k]; ok {
		delete(this.m_map, k)
		return val
	}
	return nil
}

func (this *HashMap) Vals() []interface{} {
	var vals []interface{}
	for _, v := range this.m_map {
		vals = append(vals, v)
	}
	return vals
}

func (this *HashMap) Keys() []interface{} {
	var keys []interface{}
	for k, _ := range this.m_map {
		keys = append(keys, k)
	}
	return keys
}

func (this *HashMap) Each(block func(interface{}, interface{})) {
	for k, v := range this.m_map {
		block(k, v)
	}
}

func (this *HashMap) Seach(block func(interface{}) bool) interface{} {
	for _, v := range this.m_map {
		if block(v) {
			return v
		}
	}
	return nil
}

func (this *HashMap) Clean() []interface{} {
	vals := this.Vals()
	this.m_map = make(map[interface{}]interface{})
	return vals
}
