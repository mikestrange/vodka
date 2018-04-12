package kernel

import "ants/gutil"

func _init() {
	r := NewWork()
	r.Make(1002)

	go func() {
		r.ReadMsg(NewReceiver(func(args ...interface{}) {
			println(len(args))
		}))
	}()

	str := gutil.TryFun(func() {
		for i := 0; i < 10; i++ {
			r.Push(12, 23)
		}
	})

	r.Die()
	println(str)

	println("kernel test")
}
