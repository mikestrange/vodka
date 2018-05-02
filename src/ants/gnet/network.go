package gnet

import "ants/base"

//
type INetBase interface {
	//Listen(IHandle, int, int) bool
	Done(interface{}) bool
	Write(interface{}) bool
	Read(interface{}) bool
	Exit(int) int
	SetTimeout(int)
}

//网络基础
type NetBase struct {
	INetBase
	base.Locked
	openFlag bool
	used     int
	chans    *Channels
}

func newBase() INetBase {
	return new(NetBase)
}

func (this *NetBase) Listen(handle IHandle, size int, delay int) bool {
	ok := false
	this.Lock()
	if !this.openFlag {
		ok = true
		this.chans = newChannel(size)
		this.chans.SetTimeout(delay)
		this.chans.SetHandle(handle)
		this.openFlag = true
		go this.chans.Loop()
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

func (this *NetBase) Exit(code int) int {
	this.Lock()
	if this.openFlag {
		this.used = code
		this.openFlag = false
		this.chans.CloseOf(code)
	}
	this.Unlock()
	return this.used
}

func (this *NetBase) SetTimeout(delay int) {
	this.chans.SetTimeout(delay)
}

func init() {
	println("init")
	t := newBase()
	t.Exit(1)
	//t.Listen(t, 100, 1)
}
