package gnet

import "fmt"

//解码-编码
type IProtocoler interface {
	//解码直接获得消息
	Unmarshal([]byte) []interface{}
	//编码
	Marshal(interface{}) []byte
	//编码每次尺寸
	BuffSize() int
}

//private 网络格式化处理
func NewProtocoler() IProtocoler {
	this := new(netProtocoler)
	this.init()
	return this
}

//简单的(解码/编码)操作
type netProtocoler struct {
	ByteArray     //自身作为读
	head_size int //包头
}

func (this *netProtocoler) init() {
	this.InitByteArray()
}

func (this *netProtocoler) Unmarshal(bits []byte) []interface{} {
	this.Append(bits)
	var messages []interface{}
	for {
		if this.head_size == 0 {
			this.head_size = int(this.ReadInt())
			if !check_size_ok(this.head_size) {
				panic(fmt.Sprintf("Unpack size err: head=%d", this.head_size))
			}
		} else {
			if this.Available() >= this.head_size {
				messages = append(messages, this.ReadBytes(this.head_size))
				this.head_size = 0
			} else {
				break
			}
		}
		if this.Available() == 0 {
			this.Reset()
			break
		}
	}
	return messages
}

//这里选择了复制
func (this *netProtocoler) Marshal(data interface{}) []byte {
	pack := NewByteArray()
	if bits := ToBytes(data); bits != nil {
		if size := len(bits); check_size_ok(size) {
			pack.WriteInt(int32(size))
			pack.WriteBytes(bits)
		} else {
			panic(fmt.Sprintf("Pack size err: size=%d", len(bits)))
		}
	}
	return pack.Bytes()
}

func (this *netProtocoler) BuffSize() int {
	return NET_BUFF_NEW_SIZE
}
