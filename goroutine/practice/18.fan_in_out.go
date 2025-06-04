package main

import (
	"fmt"
	"sync"
)

//18. 扇出扇入模式
//实现一个扇出扇入模式：一个 goroutine 生成数据，多个 worker goroutine 处理数据，最后合并结果。

// 扇出扇入模式是一种并发模式，其中：
// 扇出：一个 goroutine（生产者）生成数据并分发给多个 worker goroutine
// 扇入：多个 worker goroutine 处理数据后，将结果合并到一个通道
func main() {
	data := make(chan int)
	results := make(chan int, 3) // 带缓冲

	// 数据生成器（Fan-out）
	go func() {
		for i := 0; i < 10; i++ {
			data <- i
		}
		close(data)
	}()

	// 启动多个 worker
	var wg sync.WaitGroup
	numWorkers := 3

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range data {
				results <- process(d)
			}
		}()
	}

	// 合并结果（Fan-in）
	go func() {
		wg.Wait()
		close(results)
	}()

	// 输出结果
	for r := range results {
		fmt.Println(r)
	}
}

// 数据处理函数
func process(data int) int {
	return data * 2
}
