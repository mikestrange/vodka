package gutil

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
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
func Add() {
	sys_group.Add(1)
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

func MemStr(val int64) string {
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

func OsArg(idx int) string {
	if idx >= len(os.Args) {
		return ""
	}
	return os.Args[idx]
}

func OsArgs() []string {
	return os.Args
}

func Match(idx int, str string) bool {
	return OsArg(idx) == str
}

func TraceData() {
	for i := range os.Args {
		println("os arg ", i, "=", os.Args[i])
	}
}

//获取参数
func Str(idx int, str []string, def string) string {
	if idx >= len(str) {
		return def
	}
	return str[idx]
}

func Int(idx int, str []string, def int) int {
	if idx >= len(str) {
		return def
	}
	if val, err := strconv.Atoi(str[idx]); err == nil {
		return val
	}
	return def
}
