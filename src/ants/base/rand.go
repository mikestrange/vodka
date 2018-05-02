package base

//所有随机都写入里面
import "math/rand"

func init() {
	Seed(Nano())
}

func Seed(idx int64) {
	rand.Seed(idx)
}

//0-size
func Random(size int) int {
	return rand.Int() % size
}

func RandScope(b int, e int) int {
	return b + Random(e)
}
