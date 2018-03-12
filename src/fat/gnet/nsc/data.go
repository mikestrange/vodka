package nsc

//路由的数据（可以添加更多）
type IDataRoute interface {
	RouteID() int
	Addr() string
	Name() string
	Topics() []int
	HasTopic(int) bool
	SetArgs(int, string, string, ...int)
	//新增
	Port() int
	Type() int
	State() int
}

//class data
type DataRoute struct {
	routeid int
	addr    string
	name    string
	topics  []int
}

func NewDataRouteWithArgs(id int, addr string, name string, topics ...int) IDataRoute {
	this := &DataRoute{routeid: id, addr: addr, name: name, topics: topics}
	return this
}

func (this *DataRoute) SetArgs(id int, addr string, name string, topics ...int) {
	this.routeid = id
	this.addr = addr
	this.name = name
	this.topics = topics
}

func (this *DataRoute) Addr() string {
	return this.addr
}

func (this *DataRoute) Name() string {
	return this.name
}

func (this *DataRoute) Topics() []int {
	return this.topics
}

func (this *DataRoute) Port() int {
	return 0
}

func (this *DataRoute) Type() int {
	return 0
}

func (this *DataRoute) State() int {
	return 0
}

func (this *DataRoute) HasTopic(topic int) bool {
	for i := range this.topics {
		if this.topics[i] == topic {
			return true
		}
	}
	return false
}

func (this *DataRoute) RouteID() int {
	return int(this.routeid)
}
