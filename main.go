package main

import (
	"fmt"
	"time"
)

func sum(x, y int) {
	z := x + y
	fmt.Println(z)
}
func main() {
	for i := 0; i < 10; i++ {
		go sum(i,i)
	}
	time.Sleep(time.Millisecond * 10)
}