package gnet

import (
	"fmt"
)

/*
这里为底层协议, 可以进行加密处理(一般不需要改动)
*/
type NetSocketHandler struct {
	//ISocketHandler
	ByteArray
	size     int16
	pack_pos int
}

func NewSocketHandler() ISocketHandler {
	this := new(NetSocketHandler)
	this.InitSocketHandler()
	return this
}

func (this *NetSocketHandler) InitSocketHandler() {
	this.InitByteArray()
	this.size = 0
	this.pack_pos = 0
}

//interfaces
//载入数据
func (this *NetSocketHandler) LoadBytes(bits []byte) {
	this.Append(bits)
}

func (this *NetSocketHandler) Pack() (interface{}, bool) {
	if ret := this.next(); ret > 0 {
		return this.body(), true
	}
	return nil, false
}

func (this *NetSocketHandler) BuffSize() int {
	return 1024 * 10
}

/*
0表示无消息,>0表示有消息, -1表示协议错误
*/
func (this *NetSocketHandler) next() int {
	this.flush()
	if this.size == 0 {
		if this.Available() >= HEAD_SIZE {
			this.ReadValue(&this.size)
			this.pack_pos = this.Pos()
			if this.size < 0 || this.size > PACKET_MAX {
				panic(fmt.Sprint("Message Size Err:", this.size))
				return 0
			}
			if this.Available() >= int(this.size) {
				return int(this.size)
			}
		}
	} else {
		if this.Available() >= int(this.size) {
			return int(this.size)
		}
	}
	return 0
}

func (this *NetSocketHandler) flush() {
	if this.size >= 0 {
		this.SetPos(int(this.size) + this.pack_pos)
		this.size = 0
		this.pack_pos = 0
	}
	if this.Available() == 0 {
		this.Reset()
	}
}

func (this *NetSocketHandler) body() ISocketPacket {
	this.SetPos(this.pack_pos - HEAD_SIZE)
	payload := this.ReadBytes(int(this.size) + HEAD_SIZE)
	return NewPacketWithBytes(payload).ReadBegin().(ISocketPacket)
}
