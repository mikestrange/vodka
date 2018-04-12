package gutil

import "math/rand"

//0-size
func Random(size int) int {
	return rand.Int() % size
}

func RandScope(b int, e int) int {
	return b + Random(e)
}
