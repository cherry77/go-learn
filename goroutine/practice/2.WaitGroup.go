package main

import (
	"fmt"
	"sync"
)

/*
*
2. 使用 WaitGroup
修改上面的程序，使用 sync.WaitGroup 来确保主程序在所有 goroutine 完成后才退出。
*/
func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1) // 增加 WaitGroup 计数器
		go func(id int) {
			defer wg.Done() // 完成后减少计数器
			fmt.Printf("Goroutine %d 正在执行\n", id)
		}(i)
	}

	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("All goroutines finished.")
}
