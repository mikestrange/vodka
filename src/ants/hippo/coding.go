package hippo

import "encoding/binary"
import "ants/base"

const HEAD_SIZE = 4
const MAX_SIZE = 1024
const MIN_SIZE = 1

//协议
type ICoding interface {
	Decode([]byte) ([]interface{}, bool) //解码(Body)
	Encode(interface{}) ([]byte, bool)   //编码(Body)
	Buffer() []byte                      //缓冲
}

//简单的(解码/编码)操作
type NetCoding struct {
	size      int
	pos       int
	pack      []byte
	body_size int
	buff      []byte
	endian    binary.ByteOrder //大小端
}

func NewBigCoding() ICoding {
	return &NetCoding{endian: binary.BigEndian}
}

func NewLittleCoding() ICoding {
	return &NetCoding{endian: binary.LittleEndian}
}

func (this *NetCoding) Decode(b []byte) ([]interface{}, bool) {
	if len(this.pack)-this.size >= len(b) {
		copy(this.pack[this.size:], b)
	} else {
		this.pack = append(this.pack, b...)
	}
	this.size += len(b)
	var messages []interface{}
	for {
		if this.body_size == 0 {
			if this.Available() >= HEAD_SIZE {
				this.body_size = int(this.endian.Uint32(this.pack[this.pos:]))
				this.flush(HEAD_SIZE)
			} else {
				break
			}
		} else {
			if this.Available() >= this.body_size {
				messages = append(messages, this.body())
				this.flush(this.body_size)
				this.body_size = 0
			} else {
				break
			}
		}
	}
	return messages, true
}

func (this *NetCoding) Encode(data interface{}) ([]byte, bool) {
	bits := base.ToBytes(data)
	b := make([]byte, HEAD_SIZE+len(bits))
	this.endian.PutUint32(b, uint32(len(bits)))
	copy(b[HEAD_SIZE:], bits)
	return b, true
}

func (this *NetCoding) Buffer() []byte {
	if this.buff == nil {
		this.buff = make([]byte, 1024)
	}
	return this.buff
}

func (this *NetCoding) SetEndian(val binary.ByteOrder) {
	this.endian = val
}

func (this *NetCoding) Endian() binary.ByteOrder {
	return this.endian
}

func (this *NetCoding) Available() int {
	return this.size - this.pos
}

func (this *NetCoding) body() []byte {
	b := make([]byte, this.body_size)
	copy(b, this.pack[this.pos:])
	return b
}

func (this *NetCoding) flush(size int) {
	this.pos += size
	if this.Available() == 0 {
		this.pos = 0
		this.size = 0
	}
}

func _init() {
	b := NewLittleCoding()
	m, _ := b.Encode([]byte{2, 3, 1})
	b.Decode(m)
}
