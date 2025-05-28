package main

import (
	"fmt"
	"time"
)

/*
*
4. 带缓冲的通道
创建一个带缓冲的通道，缓冲区大小为 3。启动一个 goroutine 发送 5 个值到通道，主程序接收并打印这些值。观察缓冲区的效果。
*/
func main() {
	// 创建一个缓冲区大小为3的通道
	ch := make(chan int, 3)

	// 启动 goroutine 发送数据
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i // 发送数字到通道
			fmt.Printf("发送数字 %d (缓冲区剩余容量: %d)\n", i, cap(ch)-len(ch))
			time.Sleep(500 * time.Millisecond) // 模拟处理延迟
		}
		close(ch) // 发送完成后关闭通道
	}()

	// 主程序接收并打印数据
	time.Sleep(2 * time.Second) // 等待发送方先发送一些数据
	fmt.Println("开始接收数据...")

	for num := range ch {
		fmt.Printf("接收到数字: %d (缓冲区剩余容量: %d)\n", num, cap(ch)-len(ch))
		time.Sleep(1 * time.Second) // 模拟处理延迟
	}

	fmt.Println("所有数字接收完成")
}
