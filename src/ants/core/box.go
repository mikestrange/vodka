package core

import (
	"ants/base"
	"errors"
	"sync"
)

type boxSet map[int]IBox

var CloseSign = errors.New("boxCloseSign")

//上帝禁区
type IBox interface {
	Fork
	setParent(IBox)                 //设置父亲(加入后即使退出也不会，除非重新写入其他box)
	Parent() IBox                   //父亲
	PushSuper(interface{}) bool     //推送到父亲
	SetName(string)                 //设置盒子名称
	Name() string                   //盒子名称
	SetAgent(IAgent)                //设置处理者
	Agent() IAgent                  //代理
	Join(int, IBox, ...func()) bool //加入子进程
	Find(int) (IBox, bool)          //寻找
	OnReady()                       //准备进程
	CloseAll()                      //关闭所有子进程
	Send(int, interface{}) bool     //通知子节点
	Broadcast(interface{})          //通知所有子节点
	//IWork
	Push(interface{}) bool //推送自己运行
	NewTimer() ITimer      //基于时间
	Make(interface{}) bool //建立通道
}

//运行盒子
type Box struct {
	//IBox
	work    Work   //8
	boxs    boxSet //8
	wgBox   base.WaitGroup
	wgChild base.WaitGroup
	_name   string
	_super  IBox
	_agent  IAgent
	_m      sync.Mutex
}

func NewBox(agent IAgent, name string) IBox {
	this := new(Box)
	this.SetName(name)
	this.SetAgent(agent)
	return this
}

//公开给子类
func (this *Box) Lock() {
	this._m.Lock()
}

func (this *Box) Unlock() {
	this._m.Unlock()
}

//interface (只允许被设置一次)
func (this *Box) SetAgent(val IAgent) {
	if val == nil {
		this._agent = nil
	} else {
		if this._agent == nil {
			this._agent = val
		}
	}
}

func (this *Box) Agent() IAgent {
	return this._agent
}

func (this *Box) Name() string {
	return this._name
}

func (this *Box) SetName(val string) {
	this._name = val
}

func (this *Box) Run() {
	base.Assert(this.Agent() == nil, "box is not handle %s", this.Name())
	this.Loop(this.Agent())
}

func (this *Box) Wrap(block func()) {
	this.wgBox.Wrap(block)
}

func (this *Box) WrapChild(block func()) {
	this.wgChild.Wrap(block)
}

func (this *Box) Join(idx int, box IBox, args ...func()) bool {
	if this.setBox(idx, box) {
		this.handleBox(idx, box, args)
		return true
	}
	return false
}

func (this *Box) handleBox(idx int, box IBox, args []func()) {
	//设置后，不会注销
	box.setParent(this)
	//自己需要设置代理，在里面执行
	box.OnReady()
	//设置过不会再设置
	box.SetAgent(this.Agent())
	//建立过不会再建立
	box.Make(nil)
	//守护子进程
	this.WrapChild(func() {
		box.Wrap(box.Run)
		box.Wait()
		this.unBox(idx, box)
		//box.setParent(nil)
		box.OnDie()
		each_func(args)
	})
}

func (this *Box) setBox(idx int, box IBox) bool {
	this.Lock()
	if this.boxs == nil {
		this.boxs = make(boxSet)
	} else {
		if _, ok := this.boxs[idx]; ok {
			this.Unlock()
			return false
		}
	}
	this.boxs[idx] = box
	this.Unlock()
	return true
}

func (this *Box) unBox(idx int, box IBox) {
	this.Lock()
	if this.boxs != nil {
		if val, ok := this.boxs[idx]; ok {
			//判断是否一致
			if val == box {
				delete(this.boxs, idx)
			}
		}
	}
	this.Unlock()
}

func (this *Box) Parent() IBox {
	return this._super
}

func (this *Box) setParent(val IBox) {
	this._super = val
}

func (this *Box) PushSuper(event interface{}) bool {
	if this._super != nil {
		return this._super.Push(event)
	}
	return false
}

func (this *Box) Find(idx int) (IBox, bool) {
	this.Lock()
	if this.boxs != nil {
		if box, ok := this.boxs[idx]; ok {
			this.Unlock()
			return box, true
		}
	}
	this.Unlock()
	return nil, false
}

func (this *Box) CloseAll() {
	list := this.boxList(true)
	for i := range list {
		list[i].Die()
	}
}

func (this *Box) OnReady() {
	//this.SetAgent(this)
	//this.Make(val)
}

func (this *Box) OnDie() {
	println(this.Name(), "## die")
}

//=======work functions begin
func (this *Box) Die() bool {
	return this.work.Close()
}

func (this *Box) Push(event interface{}) bool {
	return this.work.Put(event)
}

func (this *Box) Make(val interface{}) bool {
	return this.work.Make(val)
}

func (this *Box) NewTimer() ITimer {
	return this.work.NewTimer()
}

//必须继承才能使用
func (this *Box) Loop(handle IAgent) {
	this.work.Loop(handle.Handle)
}

//======work functions end

func (this *Box) Wait() {
	this.wgBox.Wait()
	this.CloseAll() //关闭所有子进程
	this.wgChild.Wait()
}

func (this *Box) Send(idx int, event interface{}) bool {
	if box, ok := this.Find(idx); ok {
		return box.Push(event)
	}
	return false
}

func (this *Box) Broadcast(event interface{}) {
	list := this.boxList(false)
	//集体推送
	for i := range list {
		list[i].Push(event)
	}
}

//protected
func (this *Box) boxList(del bool) []IBox {
	var list []IBox
	this.Lock()
	if this.boxs != nil {
		for _, box := range this.boxs {
			list = append(list, box)
		}
		if del {
			this.boxs = nil
		}
	}
	this.Unlock()
	return list
}
