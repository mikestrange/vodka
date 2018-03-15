package gnet

//private 网络格式化处理
func NewSocketProcessor() INetProcessor {
	this := new(netProcessor)
	this.InitByteArray()
	return this
}

//简单的(解码/编码)操作
type netProcessor struct {
	ByteArray
	head_size int //包头
}

func (this *netProcessor) Unmarshal(bits []byte) []interface{} {
	this.Append(bits)
	var messages []interface{}
	for {
		if this.head_size == 0 {
			this.head_size = int(this.ReadShort())
			if !check_size(this.head_size) {
				panic("Unpack size error")
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
	for _, v := range args {
		if bits := ToBytes(v); bits != nil {
			if size := len(bits); check_size(size) {
				pack.WriteShort(int16(size))
				pack.WriteBytes(bits)
			}
		}
	}
	return pack.Bytes()
}
