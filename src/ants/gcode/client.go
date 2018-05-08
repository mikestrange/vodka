package gcode

//用于客户端的加密
import "fmt"
import "errors"
import "ants/base"

var private_size = 12
var private_key = []byte{1, 3, 5, 7, 9, 11, 13, 17, 23, 31, 37, 51}

//private 网络格式化处理
func NewClient() IByteCoder {
	this := new(BodyCode)
	return this
}

//简单的(解码/编码)操作
type BodyCode struct {
	base.ByteArray     //自身作为读
	head_size      int //包头
}

func (this *BodyCode) Decode(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = b[i] ^ private_key[i%private_size]
	}
	return b
}

func (this *BodyCode) Codeing(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = b[i] ^ private_key[i%private_size]
	}
	return b
}

func (this *BodyCode) Unmarshal(bits []byte) ([]interface{}, error) {
	this.Append(bits)
	var messages []interface{}
	for {
		if this.head_size == 0 {
			if this.Available() >= HEAD_SIZE {
				this.ReadValue(&this.head_size)
				//println("pack:", this.head_size, this.Available())
				if CheckSizeErr(this.head_size) {
					//panic("err size")
					return nil, errors.New(fmt.Sprintf("Unpack size err: head=%d", this.head_size))
				}
			} else {
				break
			}
		} else {
			if this.Available() >= this.head_size {
				b := this.ReadBytes(this.head_size)
				messages = append(messages, this.Decode(b))
				this.head_size = 0
				this.flush()
			} else {
				break
			}
		}
	}
	return messages, nil
}

func (this *BodyCode) flush() {
	if this.Available() == 0 {
		this.Reset()
	}
}

func (this *BodyCode) Marshal(args ...interface{}) ([]byte, error) {
	pack := base.NewByteArray()
	for i := range args {
		if bits := base.ToBytes(args[i]); bits != nil {
			if size := len(bits); CheckSizeOk(size) {
				pack.WriteValue(size, this.Codeing(bits))
			} else {
				return nil, errors.New(fmt.Sprintf("Pack size err: size=%d", len(bits)))
			}
		}
	}
	return pack.Bytes(), nil
}

func (this *BodyCode) BuffSize() int {
	return NET_BUFF_NEW_SIZE
}
