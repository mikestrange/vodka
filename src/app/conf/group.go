package conf

import "ants/core"
import "ants/gnet"
import "ants/gcode"

type CellGroup struct {
	core.Box
}

func (this *CellGroup) OnReady() {
	this.SetAgent(this)
}

func (this *CellGroup) Handle(event interface{}) {
	//集群处理
	this.Broadcast(event)
}

//配置服务器组
type CellActor struct {
	core.Box
	data *SerConf
	conn gnet.Context
}

func (this *CellActor) OnReady() {
	this.SetAgent(this)
}

func (this *CellActor) Handle(event interface{}) {
	this.Lock()
	switch d := event.(type) {
	case *gnet.SocketEvent:
		this.connect(d.BeginPack())
	default:
		println("CellActor not type handle")
	}
	this.Unlock()
}

func (this *CellActor) connect(pack gcode.ISocketPacket) {
	if this.conn == nil {
		if conn, ok := gnet.Socket(this.data.Addr); ok {
			this.conn = conn
			conn.SetReceiver(func(b []byte) {

			})
			gnet.RunAndThrowAgent(conn)
		} else {
			return
		}
	}
	this.conn.Send(pack)
}

//路由分配>>
func init() {
	if !LOCAL_TEST {
		EachVo(func(vo *SerConf) {
			if box, ok := core.Main().Find(vo.Topic); ok {
				box.Join(vo.SerID, &CellActor{data: vo}, nil)
			} else {
				core.Main().Join(vo.Topic, &CellGroup{})
			}
		})
	}
}
