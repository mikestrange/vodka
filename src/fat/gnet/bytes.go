package gnet

import (
	"encoding/binary"
	"errors"
	"fmt"
)

//默认字节
const BYTE_DEF_SIZE = 1024

//字节编译格式
var LittleEndian = binary.LittleEndian
var BigEndian = binary.BigEndian

//目前默认大端
var DefEndian = BigEndian //默认

//interfaces
type IByteArray interface {
	SetEndian(binary.ByteOrder)
	Endian() binary.ByteOrder
	//直接追加(指针不动)
	Append([]byte)
	//writes
	//Write([]byte) (int, error)
	WriteValue(...interface{})
	WriteByte(int8)
	WriteUByte(uint8)
	WriteShort(int16)
	WriteUShort(uint16)
	WriteInt(int32)
	WriteUInt(uint32)
	WriteInt64(int64)
	WriteUInt64(uint64)
	WriteString(str string)
	WriteBytes([]byte)
	//reads
	//Read([]byte) (int, error)
	ReadValue(...interface{})
	ReadByte() int8
	ReadUByte() uint8
	ReadShort() int16
	ReadUShort() uint16
	ReadInt() int32
	ReadUInt() uint32
	ReadInt64() int64
	ReadUInt64() uint64
	ReadString() string
	ReadBytes(int) []byte
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
	m_endian binary.ByteOrder
}

func NewByteArray() IByteArray {
	return NewByteArrayWithSize(BYTE_DEF_SIZE)
}

func NewByteArrayWithVals(args ...interface{}) IByteArray {
	this := NewByteArray()
	this.WriteValue(args...)
	return this
}

func NewByteArrayWithSize(size int) IByteArray {
	this := new(ByteArray)
	this.InitByteArrayWithSize(size)
	return this
}

func NewByteArrayWithBytes(bits []byte) IByteArray {
	this := new(ByteArray)
	this.InitByteArrayWithBits(bits)
	return this
}

//inits
func (this *ByteArray) InitByteArray() {
	this.InitByteArrayWithSize(BYTE_DEF_SIZE)
}

func (this *ByteArray) InitByteArrayWithSize(size int) {
	this.pos = 0
	this.size = 0
	this.cap_size = size
	this.bytes = make([]byte, size)
}

//默认在开始字段
func (this *ByteArray) InitByteArrayWithBits(bits []byte) {
	this.InitByteArrayWithSize(len(bits))
	this.WriteBytes(bits)
	this.SetBegin()
}

func (this *ByteArray) SetEndian(val binary.ByteOrder) {
	this.m_endian = val
}

func (this *ByteArray) Endian() binary.ByteOrder {
	if this.m_endian == nil {
		return DefEndian
	}
	return this.m_endian
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
		panic(errors.New("This P Size is over bytes"))
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
	default:
		this._read(val)
	}
}

func (this *ByteArray) ReadByte() int8 {
	var bit int8
	this._read(&bit)
	return bit
}

func (this *ByteArray) ReadUByte() uint8 {
	var bit uint8
	this._read(&bit)
	return bit
}

func (this *ByteArray) ReadShort() int16 {
	var bit int16
	this._read(&bit)
	return bit
}

func (this *ByteArray) ReadUShort() uint16 {
	var bit uint16
	this._read(&bit)
	return bit
}

func (this *ByteArray) ReadInt() int32 {
	var bit int32
	this._read(&bit)
	return bit
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
	pos := this.pos
	if size == 0 { //读完所有
		this.SetEnd()
		return this.bytes[pos:this.size]
	}
	if this.pos == this.size { //无字节可读
		return []byte{}
	}
	if this.pos+size > this.size { //读出超过，返回所有
		this.SetEnd()
		return this.bytes[pos:this.size]
	}
	this.SetPos(pos + size)
	return this.bytes[pos : pos+size]
}

//write
func (this *ByteArray) Write(p []byte) (n int, err error) {
	if size := len(p); size > 0 {
		this.grow(size + this.pos)
		code := copy(this.bytes[this.pos:], p)
		this.pos = this.pos + code
		if this.pos > this.size {
			this.size = this.pos
		}
		return size, nil
	}
	return 0, nil
}

//追加 不移动指针
func (this *ByteArray) Append(p []byte) {
	this.grow(len(p) + int(this.size))
	code := copy(this.bytes[this.size:], p)
	this.size = this.size + code
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
	switch val.(type) {
	case string:
		this.WriteString(val.(string))
	default:
		this._write(val)
	}
}

func (this *ByteArray) WriteValue(vals ...interface{}) {
	for _, val := range vals {
		this._write_val(val)
	}
}

func (this *ByteArray) WriteByte(val int8) {
	this._write(&val)
}

func (this *ByteArray) WriteUByte(val uint8) {
	this._write(&val)
}

func (this *ByteArray) WriteShort(val int16) {
	this._write(&val)
}

func (this *ByteArray) WriteUShort(val uint16) {
	this._write(&val)
}

func (this *ByteArray) WriteInt(val int32) {
	this._write(&val)
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
