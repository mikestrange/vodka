package gnet

import "fmt"
import "errors"

//解码-编码
type IProtocoler interface {
	//解码直接获得消息
	Unmarshal([]byte) ([][]byte, error)
	//编码
	Marshal(...interface{}) ([]byte, error)
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

func (this *netProtocoler) Unmarshal(bits []byte) ([][]byte, error) {
	this.Append(bits)
	var messages [][]byte
	for {
		if this.head_size == 0 {
			if this.Available() >= HEAD_SIZE {
				this.ReadValue(&this.head_size)
				if check_size_err(this.head_size) {
					println("err size")
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

//这里选择了复制
func (this *netProtocoler) Marshal(args ...interface{}) ([]byte, error) {
	pack := NewByteArray()
	for i := range args {
		if bits := ToBytes(args[i]); bits != nil {
			if size := len(bits); check_size_ok(size) {
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
