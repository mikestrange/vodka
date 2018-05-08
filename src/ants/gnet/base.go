package gnet

import "ants/base"

type INetBase interface {
	//生成监听
	Listen(IHandle, int) bool
	//超时
	SetTimeout(int)
	//处理
	Done(interface{}) bool
	Write(interface{}) bool
	Read(interface{}) bool
	Loop()
	LoopTimeout()
	//退出
	Exit(int)
}

//网络基础
type NetBase struct {
	INetBase
	base.Locked
	openFlag bool
	chans    *Channels
}

func newBase() INetBase {
	return new(NetBase)
}

func (this *NetBase) Listen(handle IHandle, size int) bool {
	ok := false
	this.Lock()
	if !this.openFlag {
		ok = true
		this.chans = newChannel(size)
		this.chans.SetHandle(handle)
		this.openFlag = true
	}
	this.Unlock()
	return ok
}

func (this *NetBase) Write(b interface{}) bool {
	ok := false
	this.Lock()
	if this.openFlag {
		ok = true
		this.chans.Send(b)
	}
	this.Unlock()
	return ok
}

func (this *NetBase) Read(b interface{}) bool {
	ok := false
	this.Lock()
	if this.openFlag {
		ok = true
		this.chans.Read(b)
	}
	this.Unlock()
	return ok
}

func (this *NetBase) Done(b interface{}) bool {
	ok := false
	this.Lock()
	if this.openFlag {
		ok = true
		this.chans.Done(b)
	}
	this.Unlock()
	return ok
}

func (this *NetBase) Exit(code int) {
	this.Lock()
	if this.openFlag {
		this.openFlag = false
		this.chans.CloseOf(code)
	}
	this.Unlock()
}

func (this *NetBase) SetTimeout(delay int) {
	this.Lock()
	if this.openFlag {
		this.chans.SetTimeout(delay)
	}
	this.Unlock()
}

func (this *NetBase) Loop() {
	this.chans.Loop()
}

func (this *NetBase) LoopTimeout() {
	this.chans.LoopWithTimeout()
}

func init2() {
	println("init")
	t := newBase()
	t.Exit(1)
	//t.Listen(t, 100, 1)
}
