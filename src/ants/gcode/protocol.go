package gcode

import "fmt"
import "errors"
import "ants/base"

////加密算法
//type IBodyCoder interface {
//	Decode([]byte) []byte
//	Codeing([]byte) []byte
//}

//private 网络格式化处理
func NewProtocoler() IByteCoder {
	this := new(netProtocoler)
	return this
}

//简单的(解码/编码)操作
type netProtocoler struct {
	base.ByteArray     //自身作为读
	head_size      int //包头
}

func (this *netProtocoler) Unmarshal(bits []byte) ([]interface{}, error) {
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
				messages = append(messages, this.ReadBytes(this.head_size))
				this.head_size = 0
				this.flush()
			} else {
				break
			}
		}
	}
	return messages, nil
}

func (this *netProtocoler) flush() {
	if this.Available() == 0 {
		this.Reset()
	}
}

func (this *netProtocoler) Marshal(args ...interface{}) ([]byte, error) {
	pack := base.NewByteArray()
	for i := range args {
		if bits := base.ToBytes(args[i]); bits != nil {
			if size := len(bits); CheckSizeOk(size) {
				pack.WriteValue(size, bits)
			} else {
				return nil, errors.New(fmt.Sprintf("Pack size err: size=%d", len(bits)))
			}
		}
	}
	return pack.Bytes(), nil
}

func (this *netProtocoler) BuffSize() int {
	return NET_BUFF_NEW_SIZE
}
