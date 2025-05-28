package main

import "fmt"

// 编写一个程序，goroutine 发送数据后关闭通道，主程序使用 for range 接收数据并在通道关闭后退出。
func main() {
	ch := make(chan int)
	go producer(ch)

	// 主程序作为消费者接收数据
	consumer(ch)

	fmt.Println("程序结束")
}

func producer(ch chan<- int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for num := range ch {
		println(num)
	}
	fmt.Println("检测到通道已关闭，停止接收")
}

// 总结
// 1. producer ch参数类型，是需要写数据到 ch 里，ch的类型是 chan<- int, <-chan int 是只读的
