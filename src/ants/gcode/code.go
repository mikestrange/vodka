package gcode

const (
	//默认: pack min > i < pack max
	HEAD_SIZE         = 4
	NET_BUFF_MINLEN   = 1
	NET_BUFF_MAXLEN   = 1024 * 1024 * 50 //(50MB)
	NET_BUFF_NEW_SIZE = 1024 * 10        //new read bytes
)

//解码-编码
type IByteCoder interface {
	//解码直接获得消息
	Unmarshal([]byte) ([][]byte, error)
	//编码
	Marshal(...interface{}) ([]byte, error)
	//编码每次尺寸
	BuffSize() int
}

//查看缓冲是否正确
func check_size_ok(size int) bool {
	return size >= NET_BUFF_MINLEN && size <= NET_BUFF_MAXLEN
}

//查看缓冲是否错误
func check_size_err(size int) bool {
	return size < NET_BUFF_MINLEN || size > NET_BUFF_MAXLEN
}
