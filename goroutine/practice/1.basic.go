package main

import (
	"fmt"
	"time"
)

/*
*
1. 基本 Goroutine
编写一个程序，启动 5 个 goroutine，每个 goroutine 打印自己的编号（1到5），主程序等待所有 goroutine 完成。
*/
func main() {
	for i := 0; i < 5; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d 正在执行\n", id)
		}(i)
	}

	time.Sleep(1 * time.Second) // Wait for goroutines to finish
}
