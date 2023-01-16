package main

import (
	"fmt"
	"sync"
)

func hello(i int) {
	println("hello goroutine : " + fmt.Sprint(i))
}

func ManyGoWait() {
	var wg sync.WaitGroup
	wg.Add(5) // 计数器 +5
	for i := 0; i < 5; i++ {
		go func(j int) {
			defer wg.Done() // 计数器 -1
			hello(j)
		}(i)
	}
	wg.Wait() // 阻塞直到计数器为 0
}

func main() {
	ManyGoWait()
}
