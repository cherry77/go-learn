package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//16. 上下文取消
//使用 context 包实现一个可以取消的长时间运行任务。主程序可以在任意时刻取消所有 goroutine。

// 模拟一个长时间运行的任务
func longRunningTask(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 模拟随机处理时间
	processTime := time.Duration(rand.Intn(5)+1) * time.Second
	fmt.Printf("Task %d started, will take %v\n", id, processTime)

	select {
	case <-time.After(processTime):
		// 任务正常完成
		fmt.Printf("Task %d completed successfully\n", id)
	case <-ctx.Done():
		// 任务被取消
		fmt.Printf("Task %d canceled: %v\n", id, ctx.Err())
	}
}

func main() {
	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 确保所有资源都被释放

	var wg sync.WaitGroup

	// 启动多个goroutine执行任务
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go longRunningTask(ctx, i, &wg)
	}

	// 模拟主程序运行一段时间后决定取消所有任务
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\nMain program initiating cancellation...")
		cancel() // 取消所有goroutine
	}()

	// 等待所有goroutine完成
	wg.Wait()
	fmt.Println("All tasks have finished")
}
