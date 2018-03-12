package chat

import "fmt"

type TableManager struct {
	chats map[int32]*ChatTable
}

func NewManager() *TableManager {
	this := &TableManager{chats: make(map[int32]*ChatTable)}
	//默认
	this.CreateTable(10086, 1)
	return this
}

func (this *TableManager) CreateTable(cid int32, ctype int8) (*ChatTable, bool) {
	if _, ok := this.chats[cid]; ok {
		fmt.Println("Build chat Err: have", cid)
		return nil, false
	}
	fmt.Println("Build chat Ok:", cid, ctype)
	val := NewTable(cid, ctype)
	this.chats[cid] = val
	return val, true
}

//通知所有用户频道被关闭
func (this *TableManager) RemoveTable(cid int32) (*ChatTable, bool) {
	table, ok := this.chats[cid]
	if ok {
		fmt.Println("Remove chat Ok:", cid)
		delete(this.chats, cid)
	} else {
		fmt.Println("Remove chat Err: no ", cid)
	}
	return table, ok
}

func (this *TableManager) GetTable(cid int32) (*ChatTable, bool) {
	table, ok := this.chats[cid]
	return table, ok
}
