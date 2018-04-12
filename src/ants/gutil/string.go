package gutil

/*
基础类型的一些转换
*/

import "strings"
import "fmt"
import "strconv"

//strings
func Find(src string, f string) int {
	return strings.Index(src, f)
}

func Format(str string, args ...interface{}) string {
	return fmt.Sprintf(str, args...)
}

//序列化
func Itoa(val int) string {
	return strconv.Itoa(val)
}

func Atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err == nil {
		return val
	}
	return 0
}

func Atol(str string) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return val
	}
	return 0
}

func Ltoa(val int64) string {
	return strconv.FormatInt(val, 10)
}
