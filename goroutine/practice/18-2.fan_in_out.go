package main

import (
	"fmt"
	"sync"
)

// 模拟数据生产者
func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// worker 处理函数，这里简单计算平方
func worker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n // 计算平方
		}
	}()
	return out
}

// 合并多个通道的结果
func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// 为每个输入通道启动一个 goroutine
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// 等待所有 goroutine 完成后关闭通道
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// 1. 创建输入数据
	in := producer(1, 2, 3, 4, 5, 6, 7, 8)

	// 2. 扇出 - 创建多个 worker 处理数据
	const numWorkers = 3
	workers := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = worker(in)
	}

	// 3. 扇入 - 合并所有 worker 的结果
	for result := range merge(workers...) {
		fmt.Println(result) // 打印处理后的结果
	}
}
