package hippo

import "net"
import "ants/base"

type NetService struct {
	ln     net.Listener
	wgSer  base.WaitGroup
	wgConn base.WaitGroup
}

func (this *NetService) Start(port int) {
	if ln, ok := Listen(port); ok {
		this.ln = ln
	}
}

func (this *NetService) OnReady() {

}

func (this *NetService) Wrap(f func()) {
	this.wgSer.Wrap(f)
}

func (this *NetService) Wait() {
	this.wgSer.Wait()

}

func (this *NetService) Run() {

}

func (this *NetService) OnDie() {

}

func (this *NetService) Close() error {
	return this.ln.Close()
}
