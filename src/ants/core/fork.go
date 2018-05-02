package core

//进程
type Fork interface {
	WrapChild(func()) //子进程守护
	Wrap(func())      //守护本进程
	Die() bool        //关闭自身
	OnDie()           //监听本进程的关闭
	Wait()            //等待进程和所有子进程结束
	Run()             //保存运行
}
