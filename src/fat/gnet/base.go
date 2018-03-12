package gnet

import (
	"fmt"
	"net"
)

//conn.size
const (
	HEAD_SIZE   = 2             //包头部长度
	PACKET_MAX  = 1024 * 10     //包最大10MB
	BUFFER_SIZE = 1024 * 2      //缓冲长度(默认)
	PING_TIME   = 1000 * 60 * 5 //心跳时间，秒(1000 * 60 * 5)
)

//conn.type
const (
	NET_TYPE_CONN   = 1 //连接conn类型
	NET_TYPE_SOCKET = 2 //连接socket类型
)

//conn.close.used
const (
	SIGN_CLOSE_OK     = 0 //自己主动关闭
	SIGN_CLOSE_ERROR  = 1 //报错引起关闭
	SIGN_CLOSE_SELF   = 2 //自己关闭
	SIGN_CLOSE_DISTAL = 3 //远端关闭
	SIGN_CLOSE_PING   = 4 //心跳不及时关闭
)

//events
const (
	EVENT_CONTEXT_OPEN  = 1 //运行
	EVENT_CONTEXT_CLOSE = 2 //最终关闭
	EVENT_CONTEXT_READ  = 3 //消息
	EVENT_CONTEXT_SEND  = 4 //发送
	//错误的处理
	EVENT_CONTEXT_CLOSED      = 5 //被关闭
	EVENT_CONTEXT_CLOSE_SELF  = 6 //自己关闭
	EVENT_CONTEXT_CLOSE_ERROR = 7 //错误的关闭
	//其他
	EVENT_CONTEXT_HEARTBEAT = 8 //心跳
)

//event desc
var EVENT_STR map[int]string = map[int]string{
	EVENT_CONTEXT_OPEN:        "连接打开",
	EVENT_CONTEXT_CLOSE:       "连接最终关闭",
	EVENT_CONTEXT_READ:        "连接接受数据",
	EVENT_CONTEXT_SEND:        "连接发送数据",
	EVENT_CONTEXT_CLOSED:      "连接被关闭",
	EVENT_CONTEXT_CLOSE_SELF:  "连接主动关闭",
	EVENT_CONTEXT_CLOSE_ERROR: "连接出错",
	EVENT_CONTEXT_HEARTBEAT:   "心跳连接",
}

//监听服务器端口(阻塞状态)
func ListenAndRunServer(port int, block func(interface{})) int {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err == nil {
		defer ln.Close()
		fmt.Println("Run Service Ok:", port)
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Accept Err:", err)
				continue
			} else {
				go block(conn)
			}
		}
	} else {
		fmt.Println("Run Server Err:", err)
	}
	return 0
}

//客户端链接
func DialConn(addr string) (interface{}, bool) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", addr)
		return conn, true
	}
	fmt.Println("Socket Connect Err:", err)
	return nil, false
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
