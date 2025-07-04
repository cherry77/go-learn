package main

import (
	"fmt"
	"time"
)

func main() {
	//requests := make(chan int, 5)
	//for i := 1; i <= 5; i++ {
	//	requests <- i
	//}
	//close(requests)
	//
	//limiter := time.Tick(200 * time.Millisecond)
	//
	//for req := range requests {
	//	<-limiter
	//	fmt.Println("request", req, time.Now())
	//}

	// 1. 创建突发限制器
	burstyLimiter := make(chan time.Time, 3) // 创建一个缓冲大小为3的时间通道burstyLimiter
	for range 3 {
		burstyLimiter <- time.Now() // 预先填充3个当前时间值，允许立即处理3个请求（突发能力）
	}

	// 2. 定时补充令牌
	go func() { // 启动一个goroutine，每200毫秒向限制器通道发送一个时间戳， 这实现了持续以200ms/个的速率补充处理能力
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	// 3. 创建请求队列
	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		<-burstyLimiter // 等待获取处理令牌
		fmt.Println("request", req, time.Now())
	}

}

//运行效果
//前3个请求会立即处理（利用预先填充的令牌）
//后2个请求会每200ms处理一个（等待定时器补充令牌）
//输出会显示前3个请求的时间几乎相同，后2个间隔约200ms
