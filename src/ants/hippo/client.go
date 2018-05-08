package hippo

//客户端链接

type TcpClient struct {
	addr           string
	link_size      int
	auto_reconnect bool
	tx_list        []IContext
}

func (this *TcpClient) SetAutoReconnect(val bool) {
	this.auto_reconnect = val
}

func (this *TcpClient) Connect(addr string) {
	this.addr = addr
	for i := 0; i < this.link_size; i++ {
		this.loop_context(i)
	}
}

//循环连接
func (this *TcpClient) loop_context(idx int) {

}

func (this *TcpClient) SetLinkNum(size int) {
	this.link_size = size
}

func (this *TcpClient) Close() bool {
	this.auto_reconnect = false //自己关闭
	for i := 0; i < this.link_size; i++ {
		this.tx_list[i].Close()
	}
	return false
}

func (this *TcpClient) Send(b []byte) bool {

}

func (this *TcpClient) Handle(code int, target interface{}, data interface{}) {

}
