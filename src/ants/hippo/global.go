package hippo

import "net"
import "fmt"
import "time"

func Dial(addr string) (IContext, bool) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		fmt.Println("Socket Connect Ok:", conn.RemoteAddr().String())
		return newContext(conn), true
	}
	fmt.Println("Socket Connect Err:", err)
	return nil, false
}

func Listen(port int) (net.Listener, bool) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err == nil {
		fmt.Println("Listen Service ok:", port)
		return ln, true
	}
	fmt.Println("Listen Service err:", err)
	return nil, false
}

func LoopService(ln net.Listener, block func(IContext)) error {
	var tempDelay time.Duration
	for {
		conn, err := ln.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				fmt.Printf("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0
		block(newContext(conn))
	}
	return nil
}
