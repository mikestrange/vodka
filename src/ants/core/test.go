package core

import "fmt"

func _init() {
	h := NewAgent(func(event interface{}) {
		fmt.Println("msg:", event)
	})
	s := NewBox(h, "main")
	RunAndThrowBox(s, nil, func() {
		fmt.Println("结束2")
	})
	t1 := NewBox(nil, "1 main")
	t2 := NewBox(nil, "2 main")
	s.Join(1, t1)
	s.Join(2, t2)

	for i := 0; i < 10; i++ {
		t1.Join(3+i, NewBox(nil, fmt.Sprintf("1.%d main", i)))
	}
	for i := 0; i < 10; i++ {
		t2.Join(23+i, NewBox(nil, fmt.Sprintf("2.%d main", i)))
	}
	//
	for i := 0; i < 1; i++ {
		s.Push(i)
		s.Send(1, "欧不1")
		s.Send(2, "欧不2")
	}
	s.Die()
}
