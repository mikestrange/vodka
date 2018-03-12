package gutil

/*
文本操作
*/
import (
	"bufio"
	//	"fmt"
	"os"
	"strings"
)

func NewReadForString(str string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(str))
}

//读取行
func ReadLine(red *bufio.Reader) (string, bool) {
	line, _, err := red.ReadLine()
	if err != nil {
		//无数据可读
		//fmt.Println("Read File Err:", err)
		return "", false
	}
	return string(line), true
}

//读取路径
func ReadPath(path string) []string {
	if file, err := os.Open(path); CheckSucceed(err) {
		var strs []string
		buff := bufio.NewReader(file)
		for {
			if str, ok := ReadLine(buff); ok {
				strs = append(strs, str)
			} else {
				break
			}
		}
		return strs
	}
	return nil
}
