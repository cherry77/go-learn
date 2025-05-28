package main

import (
	"fmt"
	"sync"
	"time"
)

//7. 工作池模式
//实现一个工作池：主程序创建 3 个 worker goroutine，通过通道分发任务（比如计算数字的平方），主程序收集结果。

func main() {
	// 1. 创建任务通道和结果通道（带缓冲）
	const workerNums = 3
	const jobNums = 5
	jobs := make(chan int, jobNums)
	results := make(chan int, jobNums)

	// 2. 创建WaitGroup用于等待所有worker完成
	var wg sync.WaitGroup

	// 3. 启动worker协程
	for i := 0; i < workerNums; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// 4. 发送任务到任务通道
	for i := 0; i < jobNums; i++ {
		jobs <- i
	}
	close(jobs) // 关闭任务通道（通知worker没有新任务了）

	// 5. 等待所有worker完成
	wg.Wait()
	close(results) // 关闭结果通道（所有结果已收集）

	// 6. 收集并打印结果
	for result := range results {
		fmt.Println("结果:", result)
	}
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // 通知WaitGroup当前worker已完成

	for num := range jobs {
		time.Sleep(time.Second) // 模拟处理时间
		println("Worker", id, "processing job:", num)
		results <- num * num
	}
}
