package gredis

import (
	"ants/base"
	"fmt"
	//第三方
	"github.com/garyburd/redigo/redis"
)

const (
	//基础类型
	TYPE_INT   = 1
	TYPE_INT64 = 2
	TYPE_STR   = 3
	//数组
	TYPE_BYTES  = 1
	TYPE_INTS   = 2
	TYPE_INT64S = 3
	TYPE_STRS   = 4
	//int表
	TYPE_INT_MAP_INT   = 5
	TYPE_INT_MAP_INT64 = 6
	TYPE_INT_MAP_STR   = 7
	//string表
	TYPE_STR_MAP_INT   = 8
	TYPE_STR_MAP_INT64 = 9
	TYPE_STR_MAP_STR   = 10
)

//gets
func Int(conn IConn, args ...interface{}) (int, bool) {
	ret, err := redis.Int(conn.Get(args...))
	if err == nil {
		return ret, true
	}
	return 0, false
}

func Str(conn IConn, args ...interface{}) (string, bool) {
	ret, err := redis.String(conn.Get(args...))
	if err == nil {
		return ret, true
	}
	return "", false
}

func Int64(conn IConn, args ...interface{}) (int64, bool) {
	ret, err := redis.Int64(conn.Get(args...))
	if err == nil {
		return ret, true
	}
	return 0, false
}

func Bytes(conn IConn, args ...interface{}) ([]byte, bool) {
	ret, err := redis.Bytes(conn.Get(args...))
	if err == nil {
		return ret, true
	}
	return nil, false
}

//list
func Ints(conn IConn, args ...interface{}) ([]int, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_INTS)
		if ok {
			return val.([]int), true
		}
	}
	return []int{}, false
}

func Int64s(conn IConn, args ...interface{}) ([]int64, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_INT64S)
		if ok {
			return val.([]int64), true
		}
	}
	return []int64{}, false
}

func Strs(conn IConn, args ...interface{}) ([]string, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_STRS)
		if ok {
			return val.([]string), true
		}
	}
	return []string{}, false
}

//int maps
func IntMapInt(conn IConn, args ...interface{}) (map[int]int, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_INT_MAP_INT)
		if ok {
			return val.(map[int]int), true
		}
	}
	return map[int]int{}, false
}

func IntMapInt64(conn IConn, args ...interface{}) (map[int]int64, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_INT_MAP_INT64)
		if ok {
			return val.(map[int]int64), true
		}
	}
	return map[int]int64{}, false
}

func IntMapStr(conn IConn, args ...interface{}) (map[int]string, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_INT_MAP_STR)
		if ok {
			return val.(map[int]string), true
		}
	}
	return map[int]string{}, false
}

//str maps
func StrMapInt(conn IConn, args ...interface{}) (map[string]int, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_STR_MAP_INT)
		if ok {
			return val.(map[string]int), true
		}
	}
	return map[string]int{}, false
}

func StrMapInt64(conn IConn, args ...interface{}) (map[string]int64, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_STR_MAP_INT64)
		if ok {
			return val.(map[string]int64), true
		}
	}
	return map[string]int64{}, false
}

func StrMapStr(conn IConn, args ...interface{}) (map[string]string, bool) {
	bits, ok := Bytes(conn, args...)
	if ok {
		val, ok := Unmarshal(bits, TYPE_STR_MAP_STR)
		if ok {
			return val.(map[string]string), true
		}
	}
	return map[string]string{}, false
}

//sets(转化值)
//解码(复制类型)
func Unmarshal(bits []byte, code int) (interface{}, bool) {
	pack := base.NewByteArrayWithBytes(bits)
	pack.SetBegin()
	//类型不匹配
	if ctype := int(pack.ReadByte()); code != ctype {
		panic(fmt.Sprintf("no code err code=%d, result type=%d", code, ctype))
		return nil, false
	}
	size := int(pack.ReadInt())
	switch code {
	case TYPE_INTS:
		{
			b := make([]int, size)
			for i := 0; i < size; i++ {
				b[i] = int(pack.ReadInt())
			}
			return b, true
		}
	case TYPE_INT64S:
		{
			b := make([]int64, size)
			for i := 0; i < size; i++ {
				b[i] = pack.ReadInt64()
			}
			return b, true
		}
	case TYPE_STRS:
		{
			b := make([]string, size)
			for i := 0; i < size; i++ {
				b[i] = pack.ReadString()
			}
			return b, true
		}
	case TYPE_INT_MAP_INT:
		{
			b := make(map[int]int)
			for i := 0; i < size; i++ {
				k := int(pack.ReadInt())
				v := int(pack.ReadInt())
				b[k] = v
			}
			return b, true
		}
	case TYPE_INT_MAP_INT64:
		{
			b := make(map[int]int64)
			for i := 0; i < size; i++ {
				k := int(pack.ReadInt())
				v := pack.ReadInt64()
				b[k] = v
			}
			return b, true
		}
	case TYPE_INT_MAP_STR:
		{
			b := make(map[int]string)
			for i := 0; i < size; i++ {
				k := int(pack.ReadInt())
				v := pack.ReadString()
				b[k] = v
			}
			return b, true
		}
	case TYPE_STR_MAP_INT:
		{
			b := make(map[string]int)
			for i := 0; i < size; i++ {
				k := pack.ReadString()
				v := int(pack.ReadInt())
				b[k] = v
			}
			return b, true
		}
	case TYPE_STR_MAP_INT64:
		{
			b := make(map[string]int64)
			for i := 0; i < size; i++ {
				k := pack.ReadString()
				v := pack.ReadInt64()
				b[k] = v
			}
			return b, true
		}
	case TYPE_STR_MAP_STR:
		{
			b := make(map[string]string)
			for i := 0; i < size; i++ {
				k := pack.ReadString()
				v := pack.ReadString()
				b[k] = v
			}
			return b, true
		}
	}
	return nil, false
}

//编码
func Marshal(d interface{}) (interface{}, bool) {
	switch d.(type) {
	case int, int64, string, byte, []byte:
		return d, true
	case []int:
		return packInts(d.([]int)...), true
	case []int64:
		return packInt64s(d.([]int64)...), true
	case []string:
		return packStrs(d.([]string)...), true
	case map[string]int:
		return packStrMapInt(d.(map[string]int)), true
	case map[string]int64:
		return packStrMapInt64(d.(map[string]int64)), true
	case map[string]string:
		return packStrMapStr(d.(map[string]string)), true
	case map[int]int:
		return packIntMapInt(d.(map[int]int)), true
	case map[int]int64:
		return packIntMapInt64(d.(map[int]int64)), true
	case map[int]string:
		return packIntMapStr(d.(map[int]string)), true
	default:
		panic("no redis val type")
	}
	return nil, false
}

func packInts(args ...int) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_INTS)
	pack.WriteValue(len(args))
	for i := range args {
		pack.WriteValue(args[i])
	}
	return pack.Bytes()
}

func packInt64s(args ...int64) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_INT64S)
	pack.WriteValue(len(args))
	for i := range args {
		pack.WriteValue(args[i])
	}
	return pack.Bytes()
}

func packStrs(args ...string) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_STRS)
	pack.WriteValue(len(args))
	for i := range args {
		pack.WriteValue(args[i])
	}
	return pack.Bytes()
}

//int maps
func packIntMapInt(m map[int]int) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_INT_MAP_INT)
	pack.WriteValue(len(m))
	for k, v := range m {
		pack.WriteValue(k, v)
	}
	return pack.Bytes()
}

func packIntMapInt64(m map[int]int64) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_INT_MAP_INT64)
	pack.WriteValue(len(m))
	for k, v := range m {
		pack.WriteValue(k, v)
	}
	return pack.Bytes()
}

func packIntMapStr(m map[int]string) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_INT_MAP_STR)
	pack.WriteValue(len(m))
	for k, v := range m {
		pack.WriteValue(k, v)
	}
	return pack.Bytes()
}

//str maps
func packStrMapInt(m map[string]int) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_STR_MAP_INT)
	pack.WriteValue(len(m))
	for k, v := range m {
		pack.WriteValue(k, v)
	}
	return pack.Bytes()
}

func packStrMapInt64(m map[string]int64) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_STR_MAP_INT64)
	pack.WriteValue(len(m))
	for k, v := range m {
		pack.WriteValue(k, v)
	}
	return pack.Bytes()
}

func packStrMapStr(m map[string]string) interface{} {
	pack := base.NewByteArray()
	pack.WriteByte(TYPE_STR_MAP_STR)
	pack.WriteValue(len(m))
	for k, v := range m {
		pack.WriteValue(k, v)
	}
	return pack.Bytes()
}
