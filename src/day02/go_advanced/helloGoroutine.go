package main

import (
	"fmt"
	"time"
)

func hello(i int) {
	println("hello goroutine : " + fmt.Sprint(i))
}

func HelloGoRoutine() {
	for i := 0; i < 5; i++ {
		// go 创建协程
		go func(j int) {
			hello(j)
		}(i)
	}

	// 阻塞, 保证子协程执行完前主协程不退出, 可优化
	time.Sleep(time.Second)
}

func main() {
	HelloGoRoutine()
}
