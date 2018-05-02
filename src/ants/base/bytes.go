package base

import (
	"encoding/binary"
	"errors"
	"fmt"
)

//默认字节
const BYTE_DEF_SIZE = 100

//字节编译格式
var LittleEndian = binary.LittleEndian
var BigEndian = binary.BigEndian

//目前默认大端
var DefEndian = BigEndian //默认

//获得字节接口
type IBytes interface {
	Bytes() []byte
}

//转化字节
func ToBytes(data interface{}) []byte {
	switch data.(type) {
	case IBytes:
		return data.(IBytes).Bytes()
	case []byte:
		return data.([]byte)
	case string:
		return []byte(data.(string))
	}
	return nil
}

//interfaces
type IByteArray interface {
	SetEndian(binary.ByteOrder)
	Endian() binary.ByteOrder
	//直接追加(指针不动)
	Append([]byte)
	//writes
	//Write([]byte) (int, error)
	WriteValue(...interface{})
	WriteBool(bool)
	WriteByte(int)
	WriteUByte(int)
	WriteShort(int)
	WriteUShort(int)
	WriteInt(int)
	WriteString(str string)
	WriteBytes([]byte)
	WriteUInt(uint32)
	WriteInt64(int64)
	WriteUInt64(uint64)
	//reads
	//Read([]byte) (int, error)
	ReadValue(...interface{})
	ReadByte() int
	ReadBool() bool
	ReadUByte() int
	ReadShort() int
	ReadUShort() int
	ReadInt() int
	ReadString() string
	ReadBytes(int) []byte
	ReadRemaining() []byte
	ReadUInt() uint32
	ReadInt64() int64
	ReadUInt64() uint64
	//datas
	Bytes() []byte
	GetByte(int) byte
	SetByte(int, byte)
	SetBegin()
	SetEnd()
	Length() int
	CapLen() int
	Pos() int
	SetPos(int)
	Available() int
	Reset()
	Clean()
}

//class ByteArray
type ByteArray struct {
	IByteArray
	pos      int
	size     int
	cap_size int
	bytes    []byte
	endian   binary.ByteOrder
}

func NewByteArray() IByteArray {
	return &ByteArray{}
}

func NewByteArrayWithVals(args ...interface{}) IByteArray {
	this := &ByteArray{}
	this.WriteValue(args...)
	return this
}

func NewByteArrayWithSize(size int) IByteArray {
	this := &ByteArray{}
	this.SetSize(size)
	return this
}

//这里直接引用，不会复制
func NewByteArrayWithBytes(bits []byte) IByteArray {
	this := &ByteArray{}
	this.SetBits(bits)
	return this
}

//inits

func (this *ByteArray) SetSize(size int) {
	this.pos = 0
	this.size = 0
	this.cap_size = size
	this.bytes = make([]byte, size)
}

//直接开始引用
func (this *ByteArray) SetBits(bits []byte) {
	this.bytes = bits
	this.size = len(bits)
	this.pos = this.size
	this.cap_size = this.size
}

func (this *ByteArray) SetEndian(val binary.ByteOrder) {
	this.endian = val
}

func (this *ByteArray) Endian() binary.ByteOrder {
	if this.endian == nil {
		return DefEndian
	}
	return this.endian
}

func (this *ByteArray) Pos() int {
	return this.pos
}

func (this *ByteArray) SetPos(val int) {
	this.pos = val
}

func (this *ByteArray) Length() int {
	return this.size
}

func (this *ByteArray) CapLen() int {
	return this.cap_size
}

func (this *ByteArray) Available() int {
	return this.size - this.pos
}

func (this *ByteArray) SetBegin() {
	this.pos = 0
}

func (this *ByteArray) SetEnd() {
	this.pos = this.size
}

func (this *ByteArray) Reset() {
	this.pos = 0
	this.size = 0
}

func (this *ByteArray) Clean() {
	this.pos = 0
	this.size = 0
	this.cap_size = 0
	this.bytes = make([]byte, 0)
}

func (this *ByteArray) Bytes() []byte {
	return this.bytes[0:this.size]
}

func (this *ByteArray) GetByte(pos int) byte {
	return this.bytes[pos]
}

func (this *ByteArray) SetByte(pos int, val byte) {
	this.bytes[pos] = val
}

//read
func (this *ByteArray) _read(bit interface{}) {
	if err := binary.Read(this, this.Endian().(binary.ByteOrder), bit); err != nil {
		//panic(fmt.Sprintf("Read bits Err:%v", err))
	}
}

func (this *ByteArray) Read(p []byte) (n int, err error) {
	if len(p) == 0 { //读取的参数存在问题
		panic(errors.New("This P size = 0 or nil"))
	}
	if len(p)+this.pos > this.size { //读取超过
		panic(fmt.Sprintf("This P Size is over bytes pos=%d size=%d readsize=%d", this.pos, this.size, len(p)))
	}
	size := copy(p, this.bytes[this.pos:])
	this.SetPos(this.pos + size)
	return size, nil
}

func (this *ByteArray) ReadValue(vals ...interface{}) {
	for _, val := range vals {
		this._read_val(val)
	}
}

func (this *ByteArray) _read_val(val interface{}) {
	switch val.(type) {
	case *string:
		*val.(*string) = this.ReadString()
	case *int:
		*val.(*int) = int(this.ReadInt())
	default:
		this._read(val)
	}
}

func (this *ByteArray) ReadBool() bool {
	return this.ReadByte() != 0
}

func (this *ByteArray) ReadByte() int {
	var bit int8
	this._read(&bit)
	return int(bit)
}

func (this *ByteArray) ReadUByte() int {
	var bit uint8
	this._read(&bit)
	return int(bit)
}

func (this *ByteArray) ReadShort() int {
	var bit int16
	this._read(&bit)
	return int(bit)
}

func (this *ByteArray) ReadUShort() int {
	var bit uint16
	this._read(&bit)
	return int(bit)
}

func (this *ByteArray) ReadInt() int {
	var bit int32
	this._read(&bit)
	return int(bit)
}

func (this *ByteArray) ReadUInt() uint32 {
	var bit uint32
	this._read(&bit)
	return bit
}

func (this *ByteArray) ReadInt64() int64 {
	var bit int64
	this._read(&bit)
	return bit
}

func (this *ByteArray) ReadUInt64() uint64 {
	var bit uint64
	this._read(&bit)
	return bit
}

func (this *ByteArray) ReadString() string {
	size := this.ReadInt()
	if size > 0 {
		bits := make([]byte, size)
		this._read(&bits)
		return string(bits)
	}
	return ""
}

func (this *ByteArray) ReadBytes(size int) []byte {
	if size <= 0 || size > this.Available() {
		size = this.Available()
	}
	bits := make([]byte, size)
	p := this.pos
	for i := 0; i < size; i++ {
		bits[i] = this.bytes[p+i]
	}
	this.SetPos(p + size)
	return bits
}

func (this *ByteArray) ReadRemaining() []byte {
	return this.ReadBytes(0)
}

//write
func (this *ByteArray) Write(p []byte) (n int, err error) {
	if size := len(p); size > 0 {
		this.grow(size + this.pos)
		ret := copy(this.bytes[this.pos:], p)
		this.pos = this.pos + ret
		if this.pos > this.size {
			this.size = this.pos
		}
		return size, nil
	}
	return 0, nil
}

//追加 不移动指针
func (this *ByteArray) Append(p []byte) {
	if len(p) == 0 {
		println("无效字段")
		return
	}
	this.grow(len(p) + this.size)
	ret := copy(this.bytes[this.size:], p)
	this.size = this.size + ret
}

func (this *ByteArray) grow(new_size int) {
	if new_size > this.cap_size {
		size := BYTE_DEF_SIZE * (new_size/BYTE_DEF_SIZE + 1)
		bits := make([]byte, size)
		copy(bits, this.bytes)
		this.bytes = bits
		this.cap_size = size
	}
}

func (this *ByteArray) _write(bit interface{}) {
	//系统的比较耗时(相当于泛写入)
	if err := binary.Write(this, this.Endian().(binary.ByteOrder), bit); err != nil {
		panic(fmt.Sprintln("Write bits Err:", err))
	}
}

func (this *ByteArray) _write_val(val interface{}) {
	switch v := val.(type) {
	case string:
		this.WriteString(v)
	case *string:
		this.WriteString(*v)
	case int:
		this.WriteInt(v)
	case *int:
		this.WriteInt(*v)
	case IByteArray:
		this.WriteBytes(v.Bytes())
	case []byte:
		this.WriteBytes(v)
	default:
		this._write(val)
	}
}

func (this *ByteArray) WriteValue(vals ...interface{}) {
	for _, val := range vals {
		this._write_val(val)
	}
}

func (this *ByteArray) WriteBool(val bool) {
	if val {
		this.WriteByte(1)
	} else {
		this.WriteByte(0)
	}
}

func (this *ByteArray) WriteByte(val int) {
	num := int8(val)
	this._write(&num)
}

func (this *ByteArray) WriteUByte(val int) {
	num := uint8(val)
	this._write(&num)
}

func (this *ByteArray) WriteShort(val int) {
	num := int16(val)
	this._write(&num)
}

func (this *ByteArray) WriteUShort(val int) {
	num := uint16(val)
	this._write(&num)
}

func (this *ByteArray) WriteInt(val int) {
	num := int32(val)
	this._write(&num)
}

func (this *ByteArray) WriteUInt(val uint32) {
	this._write(&val)
}

func (this *ByteArray) WriteInt64(val int64) {
	this._write(&val)
}

func (this *ByteArray) WriteUInt64(val uint64) {
	this._write(&val)
}

func (this *ByteArray) WriteString(str string) {
	bits := []byte(str)
	this._write(int32(len(bits)))
	this._write(bits)
}

func (this *ByteArray) WriteBytes(bits []byte) {
	this.Write(bits)
}
