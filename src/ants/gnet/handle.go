package gnet

import "fmt"

//private 网络格式化处理
func NewSocketProcessor() INetProcessor {
	this := new(netProcessor)
	this.init()
	return this
}

//简单的(解码/编码)操作
type netProcessor struct {
	ByteArray     //自身作为读
	head_size int //包头
}

func (this *netProcessor) init() {
	this.InitByteArray()
}

func (this *netProcessor) Unmarshal(bits []byte) []interface{} {
	this.Append(bits)
	return this.Message()
}

func (this *netProcessor) Message() []interface{} {
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
func (this *netProcessor) Marshal(args ...interface{}) []byte {
	pack := NewByteArray()
	for i := range args {
		if bits := ToBytes(args[i]); bits != nil {
			if size := len(bits); check_size_ok(size) {
				pack.WriteInt(int32(size))
				pack.WriteBytes(bits)
			}
		}
	}
	return pack.Bytes()
}

func (this *netProcessor) Commit() []byte {
	return nil
}
