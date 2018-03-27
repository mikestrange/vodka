package glog

import (
	"ants/actor"
	"fmt"
	"log"
	"os"
)

type LogStore struct {
	actor.ActorRef  //继承
	actor.BaseActor //继承
	loger           *log.Logger
	fd              *os.File
}

func NewLog() *LogStore {
	this := new(LogStore)
	this.SetActor(this)
	actor.RunWithActor(this)
	return this
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

func (this *LogStore) OnClose() {
	this.fd.Close()
}
