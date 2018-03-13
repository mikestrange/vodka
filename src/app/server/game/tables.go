package game

import "fat/gsys"
import "fmt"

type ITableLogic interface {
	TableID() int
	Type() int
	OnFree()
	OnLaunch(int, interface{})
	//推送消息
	PushNotice(...interface{})
}

var tables map[int]ITableLogic = make(map[int]ITableLogic)
var lock gsys.ILocked = gsys.NewLocked()

func init() {
	RegTable(1024, NewTexasLogic(), nil)
}

//注册启动
func RegTable(tid int, table ITableLogic, data interface{}) bool {
	lock.Lock()
	if _, ok := tables[tid]; ok {
		lock.Unlock()
		return false
	}
	tables[tid] = table
	lock.Unlock()
	fmt.Println("注册房间:", tid, table.Type())
	table.OnLaunch(tid, data)
	return true
}

//自身移除
func unRegTable(table ITableLogic) bool {
	lock.Lock()
	if table, ok := tables[table.TableID()]; ok {
		delete(tables, table.TableID())
		lock.Unlock()
		return true
	}
	lock.Unlock()
	return false
}

//外部移除他
func UnRegTableByID(tid int) (ITableLogic, bool) {
	lock.Lock()
	if table, ok := tables[tid]; ok {
		delete(tables, tid)
		lock.Unlock()
		table.OnFree()
		return table, true
	}
	lock.Unlock()
	return nil, false
}

//派送
func SendTable(tid int, args ...interface{}) {
	lock.Lock()
	if table, ok := tables[tid]; ok {
		lock.Unlock()
		table.PushNotice(args...)
	} else {
		lock.Unlock()
	}
}
