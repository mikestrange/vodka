package gnet

import (
	"fat/gutil"
)

const UINT16_MAX int = 65535
const MAGIC_KEY = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

//返回WebSocket握手协议(报错证明不是websocket协议)
func WebShakeHands(bits []byte) ([]byte, bool) {
	header_map := make(map[string]string)
	//get map
	file := gutil.NewReadForString(string(bits))
	http, ok := gutil.ReadLine(file)
	if ok {
		//println("首行:", head_info)
		if gutil.Find(http, "HTTP") == gutil.NOT_VALUE || len(http) == 0 {
			return nil, false
		}
	} else {
		return nil, false
	}
	http = http[0 : len(http)-1]
	for {
		line, ok := gutil.ReadLine(file)
		if !ok || line == "\r" {
			break
		}
		//println("line=", line)
		end := gutil.Find(line, ": ")
		if end == gutil.NOT_VALUE {
			continue
		} else {
			key := line[0:end]
			value := line[end+2:]
			header_map[key] = value
			//println("set map:", key, value)
		}
	}
	//begin start
	request := "HTTP/1.1 101 Switching Protocols\r\n"
	request += "Connection: upgrade\r\n"
	request += "Sec-WebSocket-Accept: "
	//begin end
	//获取加密key
	serverKey, ok := header_map["Sec-WebSocket-Key"]
	if ok {
		serverKey += MAGIC_KEY
		//得到加密key
		message_digest := gutil.Sha1Encode([]byte(serverKey))
		if message_digest == nil {
			return nil, false
		}
		serverKey = gutil.Base64Encode(message_digest)
		serverKey += "\r\n"
	} else {
		return nil, false
	}
	//end start
	request += serverKey
	request += "Upgrade: websocket\r\n\r\n"
	//end end
	return []byte(request), true
}

//解密(收到解析)
func WebDecode(bits []byte) (uint8, []byte) {
	var masking_key [4]byte
	payload := make([]byte, len(bits))
	//var fin uint8 = bits[0] >> 7
	var opcode uint8 = bits[0] & 0x0f //8代表发送错误
	var pos int = 1
	mask := bits[pos] >> 7
	//get size
	payload_size := int(bits[pos] & 0x7f)
	//println("head:", opcode, payload_size)
	pos = pos + 1
	if payload_size == 126 {
		bit := NewByteArrayWithSize(2)
		bit.SetEndian(BigEndian)
		bit.WriteBytes(bits[pos : pos+2])
		bit.SetBegin()
		payload_size = int(bit.ReadUShort())
		pos += 2
	} else if payload_size == 127 {
		bit := NewByteArrayWithSize(4)
		bit.SetEndian(BigEndian)
		bit.WriteBytes(bits[pos : pos+4])
		bit.SetBegin()
		pos += 4
		payload_size = int(bit.ReadInt())
	}
	//mask keys
	if mask == 1 {
		for i := 0; i < 4; i++ {
			masking_key[i] = bits[pos+i]
		}
		pos = pos + 4
		for i := 0; i < payload_size; i++ {
			payload[i] = bits[pos+i] ^ masking_key[i%4]
		}
		//println("body size:", payload_size)
		return opcode, payload[0:payload_size]
	}
	return opcode, bits[pos:]
}

//加密(多线)
func WebEncrypt(payload []byte) []byte {
	payload_size := len(payload)
	var fin uint8 = 1
	var opcode uint8 = 2
	var mask uint8 = 0 //=1使用加密,目前会报错
	var masking_key [4]byte
	bit := NewByteArray()
	var r4 uint8 = fin << 7
	r4 |= opcode & 0x0f
	bit.WriteValue(r4)
	//bit.WriteUByte(r4)
	var m1 uint8 = mask << 7
	//超出的长度
	if payload_size < 126 {
		m1 |= uint8(payload_size & 0x7f)
		bit.WriteValue(m1)
	} else if payload_size >= 126 && payload_size < UINT16_MAX {
		m1 |= 126
		bit.WriteValue(m1)
		bit.WriteValue(int16(payload_size & 0xffff))
	} else {
		m1 |= 127
		bit.WriteValue(m1)
		bit.WriteValue(int32(payload_size))
	}
	//encrypt
	if mask == 1 {
		for i := 0; i < 4; i++ {
			bit.WriteValue(int8(masking_key[i]))
		}
		for i := 0; i < payload_size; i++ {
			bit.WriteValue(int8(payload[i] ^ masking_key[i%4]))
		}
	} else {
		bit.WriteBytes(payload[0:])
	}
	return bit.Bytes()
}

//IByteParser
//type WebSocketParser struct {
//	IByteParser
//}

//func NewWebSocketParser() IByteParser {
//	return new(WebSocketParser)
//}

//func (this *WebSocketParser) Decode(bits []byte) ([]byte, bool) {
//	ret, data := WebDecode(bits)
//	if ret == 8 {
//		return nil, false
//	}
//	return data, true
//}

//func (this *WebSocketParser) Encode(bits []byte) ([]byte, bool) {
//	return WebEncrypt(bits), true
//}

//package gnet

//import (
//	"fmt"
//)

///*
//这里为底层协议, 可以进行加密处理(一般不需要改动)
//*/
//type NetDepack struct {
//	ISocketHandler
//	ByteArray
//	parser   IByteParser //解析
//	size     int16
//	pack_pos int32
//}

//func NewDepack() ISocketHandler {
//	this := new(NetDepack)
//	this.InitNetDepack()
//	return this
//}

//func (this *NetDepack) InitNetDepack() {
//	this.InitByteArray()
//	this.size = 0
//	this.pack_pos = 0
//	this.parser = nil
//}

////interfaces
////载入数据
//func (this *NetDepack) LoadBytes(tx INetContext, bits []byte) {
//	if tx.AsSocket() {
//		this.parser = this
//	} else {
//		if this.parser == nil {
//			if data, ok := WebShakeHands(bits); ok {
//				this.parser = NewWebSocketParser()
//				fmt.Println("[WebSocke Connect]")
//				tx.Send(data)
//				return
//			} else {
//				this.parser = this
//			}
//		}
//	}
//	//解析了
//	if datas, ok := this.parser.Decode(bits); ok {
//		this.Append(datas)
//	}
//}

//func (this *NetDepack) Pack() (interface{}, bool) {
//	if ret := this.next(); ret > 0 {
//		return this.body(), true
//	}
//	return nil, false
//}

////发送载入包(目前直接发送) (如果包太大不能一次发送,分几次)
//func (this *NetDepack) PushPack(tx INetContext, data interface{}) {
//	bits := ToBytes(data)
//	if tx.AsSocket() {
//		if datas, ok := this.Encode(bits); ok {
//			SendConn(tx.Conn(), datas)
//		}
//	} else {
//		if datas, ok := this.parser.Encode(bits); ok {
//			SendConn(tx.Conn(), datas)
//		}
//	}
//}

//func (this *NetDepack) UnPack() ([]byte, bool) {
//	return nil, false
//}

//func (this *NetDepack) BuffSize() int {
//	return 1024 * 10
//}

///*
//0表示无消息,>0表示有消息, -1表示协议错误
//*/
//func (this *NetDepack) next() int {
//	this.flush()
//	if this.size == 0 {
//		if this.Available() >= HEAD_SIZE {
//			this.ReadValue(&this.size)
//			this.pack_pos = this.GetPosition()
//			if this.size < 0 || this.size > MESSAGE_PACKET_MAX {
//				panic(fmt.Sprint("消息长度错误", this.size))
//				return NET_ERROR
//			}
//			if this.Available() >= int32(this.size) {
//				return int(this.size)
//			}
//		}
//	} else {
//		if this.Available() >= int32(this.size) {
//			return int(this.size)
//		}
//	}
//	return 0
//}

//func (this *NetDepack) flush() {
//	if this.size >= 0 {
//		this.SetPosition(int32(this.size) + this.pack_pos)
//		this.size = 0
//		this.pack_pos = 0
//	}
//	if this.Available() == 0 {
//		this.Reset()
//	}
//}

//func (this *NetDepack) body() ISocketPacket {
//	this.SetPosition(this.pack_pos - HEAD_SIZE)
//	payload := this.ReadBytes(int32(this.size) + HEAD_SIZE)
//	return NewPacketWithBytes(payload).ReadBegin().(ISocketPacket)
//}

////
//func (this *NetDepack) Decode(bits []byte) ([]byte, bool) {
//	return bits, true
//}

//func (this *NetDepack) Encode(bits []byte) ([]byte, bool) {
//	return bits, true
//}
