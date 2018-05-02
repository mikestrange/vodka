package base

import (
	"os"
	"strconv"
)

/*
内存计算
*/
const KB_BIT = 1024
const MB_BIT = KB_BIT * KB_BIT
const GB_BIT = KB_BIT * MB_BIT
const TB_BIT = KB_BIT * GB_BIT

func MemStr(val int64) string {
	if val > TB_BIT {
		return Format("%dTB", val/TB_BIT)
	} else if val > GB_BIT {
		return Format("%dGB", val/GB_BIT)
	} else if val > MB_BIT {
		return Format("%dMB", val/MB_BIT)
	} else if val > KB_BIT {
		return Format("%dKB", val/KB_BIT)
	}
	return Format("%dB", val)
}

//应用程序路径
func Pwd() string {
	return os.Args[0]
}

//系统参数列表
func OsArgs() []string {
	return os.Args
}

//系统参数
func OsArg(idx int) string {
	if len(os.Args) > idx {
		return os.Args[idx]
	}
	return ""
}

//系统参数比较
func OsCompare(idx int, str string) bool {
	return OsArg(idx) == str
}

//输出查看
func OsTrace() {
	for i := range os.Args {
		println("app Arg[", i, "]=", os.Args[i])
	}
}

//获取参数
func Str(idx int, str []string, set string) string {
	if len(str) > idx {
		return str[idx]
	}
	return set
}

//列表获取参数
func Int(idx int, str []string, set int) int {
	if len(str) > idx {
		if val, err := strconv.Atoi(str[idx]); err == nil {
			return val
		}
	}
	return set
}
