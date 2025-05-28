package main

import (
	"fmt"
	"time"
)

//8. select 多路复用
//创建一个程序，有两个 goroutine 分别向两个不同的通道发送数据。使用 select 语句接收这两个通道的数据并打印。

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 启动第一个生产者goroutine
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- fmt.Sprintf("通道1-数据%d", i)
			time.Sleep(1 * time.Second) // 每秒发送一次
		}
		close(ch1) // 发送完成后关闭通道
	}()

	// 启动第二个生产者goroutine
	go func() {
		for i := 0; i < 3; i++ {
			ch2 <- fmt.Sprintf("通道2-数据%d", i)
			time.Sleep(2 * time.Second) // 每2秒发送一次
		}
		close(ch2) // 发送完成后关闭通道
	}()

	// 使用select多路复用接收数据
	ch1Closed, ch2Closed := false, false

	for {
		select {
		case data, ok := <-ch1:
			if !ok {
				fmt.Println("通道1已关闭")
				ch1Closed = true
			} else {
				fmt.Println("接收到:", data)
			}
		case data, ok := <-ch2:
			if !ok {
				fmt.Println("通道2已关闭")
				ch2Closed = true
			} else {
				fmt.Println("接收到:", data)
			}
		case <-time.After(1 * time.Second):
			fmt.Println("等待超时")
			return
		}

		// 当两个通道都关闭时退出循环
		if ch1Closed && ch2Closed {
			break
		}
	}

	fmt.Println("程序结束")
}

// 总结：
// 1. 两个通道关闭时一定要退出循环
// 2. 发送完成后一定要关闭通道
// 3. 注意退出循环的标签位置
