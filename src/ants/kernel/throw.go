package kernel

//函数保持(无法继承)
type IKeeper interface {
	//Try(func()) IKeeper              //捕获错误
	Catch(func(interface{})) IKeeper //处理错误
	Die(func()) IKeeper              //无论错误都执行
	Done()                           //完成try
}

type Throw struct {
	tryBlock func()
	errBlock func(interface{})
	endBlock func()
}

//捕获错误
func Try(block func()) IKeeper {
	this := &Throw{}
	return this.Try(block)
}

//私有
func (this *Throw) Try(try func()) IKeeper {
	this.tryBlock = try
	return this
}

//public
func (this *Throw) Catch(err func(interface{})) IKeeper {
	this.errBlock = err
	return this
}

func (this *Throw) Die(end func()) IKeeper {
	this.endBlock = end
	return this
}

//行动
func (this *Throw) Done() {
	defer func() {
		if err := recover(); err != nil {
			this.OnErr(err)
		}
		this.OnFinal()
	}()
	this.tryBlock()
}

//错误处理
func (this *Throw) OnErr(err interface{}) {
	if this.errBlock != nil {
		this.errBlock(err)
	} else {
		println("ignore try err")
	}
}

//最终通知
func (this *Throw) OnFinal() {
	if this.endBlock != nil {
		this.endBlock()
	}
}
