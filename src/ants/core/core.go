package core

import "ants/base"

//系统盒子
var sys_box IBox

func init() {
	sys_box = RunAndThrowBox(NewBase(nil, "主程序"), nil)
}

func Main() IBox {
	return sys_box
}

//独立运行的盒子: 没有上级运行
func RunAndThrowBox(box IBox, val interface{}, args ...func()) IBox {
	if box.Make(val) {
		box.OnReady()
		box.Wrap(box.Run)
		WaitDaemon(box)
	}
	return box
}

//进程守护
func WaitDaemon(fork Fork, args ...func()) {
	base.TryGo(func() {
		fork.Wait()
		fork.OnDie()
		each_func(args)
	}, func(ok bool) {
		println("exit fork: ", ok)
	})
}

//执行关闭函数
func each_func(args []func()) {
	for i := range args {
		args[i]()
	}
}
