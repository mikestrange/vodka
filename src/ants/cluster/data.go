package cluster

import "ants/conf"

//路由的数据（可以添加更多）
type IDataRoute interface {
	Addr() string
	Name() string
	Topic() int
	Port() int
	Type() int
	HasTopic(int) bool
	SetArgs(string, string, int)
}

//class data
type DataRoute struct {
	addr   string
	name   string
	topics int
	port   int
	mtype  int
	state  int
}

func NewDataWithVo(vo *conf.RouteVo) IDataRoute {
	return &DataRoute{addr: vo.Addr, name: vo.Name, topics: vo.Topic, port: vo.Port, mtype: vo.Type}
}

func NewData(port int) IDataRoute {
	return NewDataWithVo(conf.GetRouter(port))
}

func (this *DataRoute) SetArgs(addr string, name string, topics int) {
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

func (this *DataRoute) Topic() int {
	return this.topics
}

func (this *DataRoute) Port() int {
	return this.port
}

func (this *DataRoute) Type() int {
	return this.mtype
}

func (this *DataRoute) HasTopic(topic int) bool {
	return this.topics == topic
}
