package main

import (
	"fmt"
	"sync"
	"time"
)

// 17. 漏桶限流器
// 实现一个漏桶算法限流器，控制 goroutine 的并发请求速率。
// 漏桶算法是一种常用的流量整形和速率限制算法，它以一个固定的速率处理请求。下面是一个Go语言实现的漏桶限流器，可以控制goroutine的并发请求速率。

// LeakyBucket 漏桶结构体
type LeakyBucket struct {
	capacity  int           // 桶的容量
	remaining int           // 桶中剩余的量
	rate      time.Duration // 漏出的速率
	last      time.Time     // 上次漏水时间
	mu        sync.Mutex    // 互斥锁
}

// NewLeakyBucket 创建一个新的漏桶
func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity:  capacity,
		remaining: capacity,
		rate:      rate,
		last:      time.Now(),
	}
}

// Allow 检查是否允许通过
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 计算从上一次到现在漏出了多少
	now := time.Now()
	elapsed := now.Sub(lb.last)
	leak := int(elapsed / lb.rate)

	if leak > 0 {
		lb.remaining += leak
		if lb.remaining > lb.capacity {
			lb.remaining = lb.capacity
		}
		lb.last = now
	}

	if lb.remaining > 0 {
		lb.remaining--
		return true
	}
	return false
}

// Wait 等待直到允许通过
func (lb *LeakyBucket) Wait() {
	for !lb.Allow() {
		time.Sleep(lb.rate)
	}
}

func main() {
	// 创建一个漏桶：容量为5，每秒漏出1个
	bucket := NewLeakyBucket(100, time.Second)

	var wg sync.WaitGroup

	// 模拟10个请求
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 等待直到允许通过
			bucket.Wait()

			fmt.Printf("Request %d allowed at %v\n", id, time.Now().Format("15:04:05.000"))
		}(i)
	}

	wg.Wait()
}
