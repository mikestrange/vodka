package gutil

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

//
var sys_group *sync.WaitGroup

func init() {
	//系统的
	sys_group = new(sync.WaitGroup)
	//系统运行时间
	BeginTime = GetTimer()
	//随机因子
	UpdateSeed()
	//当前路径
	Pwd()
}

func UpdateSeed() {
	rand.Seed(time.Now().UnixNano())
}

//等待系统
func Add(val int) {
	sys_group.Add(val)
}

func Done() {
	sys_group.Done()
}

func Wait() {
	sys_group.Wait()
}

/*
暂停毫秒(毫秒计算)
*/
func Sleep(val int) {
	time.Sleep(time.Millisecond * time.Duration(val))
}

/*
退出系统
*/
func ExitSystem(code int) {
	os.Exit(code)
}

/*
有错误
*/
func CheckError(err interface{}) bool {
	if err == nil {
		return false
	}
	fmt.Println("CheckError :", err)
	return true
}

/*
无错误
*/
func CheckSucceed(err interface{}) bool {
	if err == nil {
		return true
	}
	fmt.Println("CheckSucceed :", err)
	return false
}

//输出错误
func TryError(args ...interface{}) {
	if err := recover(); err != nil {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			fuc := runtime.FuncForPC(pc)
			fmt.Println("[", file, line, fuc.Name(), "] Err=", err, ",Arg=", args)
		} else {
			fmt.Println("TryError :", err)
		}
	}
}

func TryCatch(okfunc func(), errfunc func(interface{})) {
	if err := recover(); err != nil {
		errfunc(err)
	} else {
		okfunc()
	}
}

/*
内存计算
*/
const KB_BIT = 1024
const MB_BIT = KB_BIT * KB_BIT
const GB_BIT = KB_BIT * MB_BIT
const TB_BIT = KB_BIT * GB_BIT

func MemString(val uint64) string {
	if val > TB_BIT {
		return fmt.Sprintf("%dTB", val/TB_BIT)
	} else if val > GB_BIT {
		return fmt.Sprintf("%dGB", val/GB_BIT)
	} else if val > MB_BIT {
		return fmt.Sprintf("%dMB", val/MB_BIT)
	} else if val > KB_BIT {
		return fmt.Sprintf("%dKB", val/KB_BIT)
	}
	return fmt.Sprintf("%dB", val)
}

/*
获取当前路径
*/
var root_path interface{}

func Pwd() string {
	if root_path == nil {
		index := strings.LastIndex(os.Args[0], "/")
		if index != -1 {
			root_path = os.Args[0][0:index]
		} else {
			root_path = "."
		}
	}
	return root_path.(string)
}

//匹配系统参数
func MatchSys(idx int, str string) bool {
	if idx >= len(os.Args) {
		return false
	}
	return os.Args[idx] == str
}

func GetArgs(idx int) string {
	return os.Args[idx]
}

func TraceData() {
	for i := range os.Args {
		fmt.Println(i, "系统参数:", os.Args[i])
	}
}
