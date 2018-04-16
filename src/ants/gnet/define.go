package gnet

import "time"
import "fmt"
import "net"

//全局方法和常量

const (
	//handle
	EVENT_CONN_READ      = 1
	EVENT_CONN_SEND      = 2
	EVENT_HEARTBEAT_PINT = 3 //心跳
	//close sign
	SIGN_CLOSE_OK     = 0 //自己关闭
	SIGN_CLOSE_DISTAL = 1 //对方关闭
	//	SIGN_CLOSE_ERROR      = 2
	//	SIGN_CLOSE_ERROR_READ = 3
	//	SIGN_CLOSE_ERROR_SEND = 4
	//	SIGN_CLOSE_HEARTBEAT  = 5
	//	SIGN_READ_ERROR       = 6
	//	SIGN_SEND_ERROR       = 7
	//默认: chan size
	NET_CHAN_SIZE = 1000
	//默认: max server conn (2万表示无压力)
	NET_SERVER_CONN_SIZE = 20000
	//默认: pack min > i < pack max
	HEAD_SIZE         = 4
	NET_BUFF_MINLEN   = 1
	NET_BUFF_MAXLEN   = 1024 * 1024 * 50 //(50MB)
	NET_BUFF_NEW_SIZE = 1024 * 10        //new read bytes
	//默认: Heartbeat time
	PING_TIME = 1000 * 60 * 5 //5分钟
)

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

//查看缓冲是否正确
func check_size_ok(size int) bool {
	return size >= NET_BUFF_MINLEN && size <= NET_BUFF_MAXLEN
}

//查看缓冲是否错误
func check_size_err(size int) bool {
	return size < NET_BUFF_MINLEN || size > NET_BUFF_MAXLEN
}

//查看端口是否错误，错误就关闭服务
func check_server_error(err error) bool {
	var tempDelay time.Duration = 0
	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}
		fmt.Println("Accept error: ", err, "; retrying in ", tempDelay)
		time.Sleep(tempDelay)
	} else {
		fmt.Println("Accept Err:", err)
		return true
	}
	return false
}
