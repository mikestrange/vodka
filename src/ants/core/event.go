package core

/*
这是一个扩展的方式，如果需要自定义可以借鉴
*/

//默认方式
type EventHandle func(interface{})

//直接回调
type IBlock interface {
	Result()
}

//事件的形式
type IEvent interface {
	Type() int
	Data() interface{}
}

func NewEvent(event int, data interface{}) IEvent {
	return &Event{event: event, data: data}
}

type Event struct {
	event int
	data  interface{}
}

func (this *Event) Type() int {
	return this.event
}

func (this *Event) Data() interface{} {
	return this.data
}

//基础的盒子，保护了一些处理机制
type BaseBox struct {
	Box
	evts  map[int]EventHandle
	block EventHandle
}

func NewBase(val EventHandle, name string) *BaseBox {
	this := new(BaseBox)
	this.SetName(name)
	this.SetAgent(this) //默认自己处理
	this.SetBlock(val)  //设置回调
	return this
}

func (this *BaseBox) SetBlock(val EventHandle) {
	this.block = val
}

func (this *BaseBox) On(event int, block EventHandle) {
	if this.evts == nil {
		this.evts = make(map[int]EventHandle)
	}
	this.evts[event] = block
}

func (this *BaseBox) Off(event int) {
	if this.evts != nil {
		delete(this.evts, event)
	}
}

func (this *BaseBox) Do(event IEvent) {
	if this.evts != nil {
		if f, ok := this.evts[event.Type()]; ok {
			f(event.Data())
		}
	}
}

func (this *BaseBox) doBlcok(data interface{}) {
	if this.block != nil {
		this.block(data)
	} else {
		println("not box block:", this.Name())
	}
}

func (this *BaseBox) Handle(event interface{}) {
	switch f := event.(type) {
	case func():
		f()
	case IBlock:
		f.Result()
	case IEvent:
		this.Do(f)
	case error:
		this.Die()
	default:
		this.doBlcok(event)
	}
}
