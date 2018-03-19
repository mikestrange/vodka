package gredis

import (
	"fmt"
	//第三方
	"github.com/garyburd/redigo/redis"
)

//用户参数格式
func ToUser(uid int, param string) string {
	return fmt.Sprintf("user.%d.%s", uid, param)
}

type IConn interface {
	Connect(string) bool
	//具体链接
	Conn() redis.Conn
	//地址
	Addr() string
	//关闭
	Close() error
	//清空(慎用)
	Clear() error
	//操作
	Do(string, ...interface{}) (interface{}, error)
	//获取
	Get(...interface{}) (interface{}, error)
	//删除(永远返回成功)
	Del(...interface{}) bool
	//设置(系统的，无任何处理)
	Set(...interface{}) bool
	//设置(自己的格式化)
	SetVal(string, interface{}) bool
	//秒(自己的格式化)
	SetEx(string, interface{}, int) bool
	//毫秒(自己的格式化)
	SetPx(string, interface{}, int) bool
	//设置一个用户数据(过期时间为秒,0为不过期)
	SetUser(int, string, interface{}, int) bool
	//获取用户数据
	GetUser(int, string) (interface{}, bool)
}

//本地
func NewConn() IConn {
	this := new(redis_client)
	this.Connect("127.0.0.1:6379")
	return this
}

func NewConnAddr(addr string) (IConn, bool) {
	this := new(redis_client)
	return this, this.Connect(addr)
}

type redis_client struct {
	conn redis.Conn
	addr string
}

func (this *redis_client) Connect(addr string) bool {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Redis Err:", err)
		return false
	}
	fmt.Println("Redis Ok:", addr)
	this.conn = conn
	this.addr = addr
	return true
}

func (this *redis_client) Conn() redis.Conn {
	return this.conn
}

func (this *redis_client) Addr() string {
	return this.addr
}

func (this *redis_client) Close() error {
	return this.conn.Close()
}

func (this *redis_client) Do(key string, args ...interface{}) (interface{}, error) {
	return this.conn.Do(key, args...)
}

func (this *redis_client) Clear() error {
	_, err := this.Do("flushall")
	if err != nil {
		fmt.Println("flushall err:", err)
	}
	return err
}

func (this *redis_client) Get(args ...interface{}) (interface{}, error) {
	ret, err := this.conn.Do("get", args...)
	if err == nil {
		return ret, nil
	}
	fmt.Println("get err:", err)
	return nil, err
}

func (this *redis_client) Set(args ...interface{}) bool {
	ret, err := this.conn.Do("set", args...)
	if err == nil {
		return true
	}
	fmt.Println(ret, "set err:", err)
	return false
}

func (this *redis_client) Del(args ...interface{}) bool {
	for i := range args {
		ret, err := this.conn.Do("del", args[i])
		if err == nil {
			fmt.Println("Del Ok:", args[i])
		} else {
			fmt.Println("Del Err:", err, ",ret=", ret, ",key=", args[i])
		}
	}
	return true
}

func (this *redis_client) SetVal(key string, data interface{}) bool {
	if val, ok := Marshal(data); ok {
		return this.Set(key, val)
	}
	return false
}

func (this *redis_client) SetEx(key string, data interface{}, delay int) bool {
	if val, ok := Marshal(data); ok {
		return this.Set(key, val, "EX", delay)
	}
	return false
}

func (this *redis_client) SetPx(key string, data interface{}, delay int) bool {
	if val, ok := Marshal(data); ok {
		return this.Set(key, val, "PX", delay)
	}
	return false
}

func (this *redis_client) SetUser(uid int, param string, val interface{}, delay int) bool {
	if delay <= 0 {
		return this.SetVal(ToUser(uid, param), val)
	}
	return this.SetEx(ToUser(uid, param), val, delay)
}

func (this *redis_client) GetUser(uid int, param string) (interface{}, bool) {
	if ret, err := this.Get(ToUser(uid, param)); err == nil {
		return ret, true
	}
	return nil, false
}
