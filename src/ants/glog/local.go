package glog

//本地日志文件
import (
	"ants/base"
	"ants/core"
	"os"
)

//异步日志
type ILogStore interface {
	Open(string) bool
	Write(string, ...interface{})
	Close() bool
	SetEnd() //设置在最后追加的方式
	SetPos(int64)
	Pos() int64
}

//本地文件日志
type LogStore struct {
	core.Box
	//---
	path string
	pos  int64
	fd   *os.File
}

//独立文件系统
func NewStore(path string) ILogStore {
	this := new(LogStore)
	if this.Open(path) {
		core.RunAndThrowBox(this, nil)
	}
	return this
}

func (this *LogStore) OnReady() {
	this.SetName("本地日志进程")
	this.SetAgent(this)
}

func (this *LogStore) Handle(event interface{}) {
	str := event.(string)
	str = base.FromtALL() + "#" + str + "\n"
	n, err := this.fd.WriteAt([]byte(str), this.pos)
	if err == nil {
		this.SetPos(this.pos + int64(n))
	} else {
		Debug("[%s] file write err:%v", this.path, err)
	}
}

func (this *LogStore) OnDie() {
	this.fd.Close()
}

func (this *LogStore) SetPos(pos int64) {
	this.pos = pos
}

func (this *LogStore) Pos() int64 {
	return this.pos
}

func (this *LogStore) SetEnd() {
	pos, err := this.fd.Seek(0, os.SEEK_END)
	if err == nil {
		this.SetPos(pos)
	}
}

func (this *LogStore) Open(path string) bool {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err == nil {
		this.path = path
		this.fd = fd
		Debug("[%s] file open ok", path)
		return true
	}
	Debug("[%s] file open err:%v", path, err)
	return false
}

func (this *LogStore) Write(str string, args ...interface{}) {
	this.Push(base.Format(str, args...))
}

func (this *LogStore) Close() bool {
	return this.Die()
}
