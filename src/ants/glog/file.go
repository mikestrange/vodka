package glog

import (
	"ants/actor"
	"fmt"
	"log"
	"os"
)

//文件日志
type LogStore struct {
	actor.BaseBox
	loger *log.Logger
	fd    *os.File
}

func NewLog() *LogStore {
	this := new(LogStore)
	actor.RunAndThrowBox(this, nil)
	return this
}

func (this *LogStore) OnReady() {
	this.SetActor(this)
}

func (this *LogStore) OpenFile(path string) bool {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err == nil {
		this.fd = fd
		this.loger = log.New(fd, "", log.Ldate|log.Ltime)
		return true
	}
	fmt.Println("Log Error:", err)
	return false
}

//实现
func (this *LogStore) OnMessage(args ...interface{}) {
	this.loger.Println(args...)
}

func (this *LogStore) OnDie() {
	this.fd.Close()
}
