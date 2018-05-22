package hippo

//事件
type IEvent interface {
	Check(int) bool
	Code() int
	Caller() interface{}
	Data() interface{}
}

type Event struct {
	caller interface{}
	code   int
	data   interface{}
}

func NewEvent(code int, caller interface{}, data interface{}) IEvent {
	return &Event{code: code, caller: caller, data: data}
}

func (this *Event) Check(code int) bool {
	return this.code == code
}

func (this *Event) Code() int {
	return this.code
}

func (this *Event) Caller() interface{} {
	return this.caller
}

func (this *Event) Data() interface{} {
	return this.data
}
