package cluster

import "ants/conf"

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
	SetStatus(int)
}

//class data
type DataRoute struct {
	routeid int
	addr    string
	name    string
	topics  []int
	port    int
	mtype   int
	state   int
}

func NewDataRouteWithVo(vo *conf.RouteVo) IDataRoute {
	return &DataRoute{routeid: vo.Id, addr: vo.Addr, name: vo.Name, topics: vo.Topic, port: vo.Port, mtype: vo.Type}
}

func NewDataRoute(port int) IDataRoute {
	return NewDataRouteWithVo(conf.GetRouter(port))
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
	return this.port
}

func (this *DataRoute) Type() int {
	return this.mtype
}

func (this *DataRoute) State() int {
	return this.state
}

func (this *DataRoute) SetStatus(val int) {
	this.state = val
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
