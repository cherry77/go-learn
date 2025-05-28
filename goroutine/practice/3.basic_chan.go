package main

import "fmt"

/*
*
3. 通道基础
创建一个 goroutine 生成 1 到 10 的数字并通过通道发送，主程序接收并打印这些数字。
*/
func main() {
	// 创建一个无缓冲通道
	ch := make(chan int)

	// 启动 goroutine 发送数据
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i // 发送数字到通道
		}
		close(ch) // 发送完成后关闭通道
	}()

	for num := range ch {
		fmt.Println("接收到的数字:", num)
	}
}
