package gcode

import "strings"
import "ants/base"
import "bufio"

const NOT_VALUE = -1
const UINT16_MAX int = 65535
const MAGIC_KEY = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

//读取行
func ReadLine(file *bufio.Reader) (string, bool) {
	line, _, err := file.ReadLine()
	if err != nil {
		//无数据可读
		//fmt.Println("Read File Err:", err)
		return "", false
	}
	return string(line), true
}

//返回WebSocket握手协议(报错证明不是websocket协议)
func WebShakeHands(bits []byte) ([]byte, bool) {
	header_map := make(map[string]string)
	//get map
	file := bufio.NewReader(strings.NewReader(string(bits)))
	http, ok := ReadLine(file)
	if ok {
		//println("首行:", head_info)
		if strings.Index(http, "HTTP") == NOT_VALUE || len(http) == 0 {
			return nil, false
		}
	} else {
		return nil, false
	}
	http = http[0 : len(http)-1]
	for {
		line, ok := ReadLine(file)
		if !ok || line == "\r" {
			break
		}
		//println("line=", line)
		end := strings.Index(line, ": ")
		if end == NOT_VALUE {
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
		message_digest := base.Sha1Encode([]byte(serverKey))
		if message_digest == nil {
			return nil, false
		}
		serverKey = base.Base64Encode(message_digest)
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
		bit := base.NewByteArrayWithSize(2)
		bit.SetEndian(base.BigEndian)
		bit.WriteBytes(bits[pos : pos+2])
		bit.SetBegin()
		payload_size = int(bit.ReadUShort())
		pos += 2
	} else if payload_size == 127 {
		bit := base.NewByteArrayWithSize(4)
		bit.SetEndian(base.BigEndian)
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
	bit := base.NewByteArray()
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
