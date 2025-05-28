package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 17. 漏桶限流器
// 实现一个漏桶算法限流器，控制 goroutine 的并发请求速率。
//高级实现（带上下文取消）
//下面是一个更高级的实现，支持上下文取消和更精确的速率控制：

// LeakyBucket 高级漏桶实现
type LeakyBucket struct {
	capacity  int           // 桶的容量
	remaining int           // 桶中剩余的量
	rate      time.Duration // 漏出的速率
	last      time.Time     // 上次漏水时间
	mu        sync.Mutex    // 互斥锁
	ticker    *time.Ticker  // 定时器
	closeCh   chan struct{} // 关闭通道
}

// NewLeakyBucket 创建漏桶
func NewLeakyBucket(capacity int, per time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		capacity:  capacity,
		remaining: capacity,
		rate:      per,
		last:      time.Now(),
		closeCh:   make(chan struct{}),
	}

	// 启动后台漏水协程
	lb.ticker = time.NewTicker(lb.rate)
	go lb.leak()

	return lb
}

// leak 后台漏水
func (lb *LeakyBucket) leak() {
	for {
		select {
		case <-lb.ticker.C:
			lb.mu.Lock()
			if lb.remaining < lb.capacity {
				lb.remaining++
			}
			lb.mu.Unlock()
		case <-lb.closeCh:
			return
		}
	}
}

// Close 关闭漏桶
func (lb *LeakyBucket) Close() {
	close(lb.closeCh)
	lb.ticker.Stop()
}

// Allow 检查是否允许通过
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if lb.remaining > 0 {
		lb.remaining--
		return true
	}
	return false
}

// Wait 等待直到允许通过
func (lb *LeakyBucket) Wait(ctx context.Context) error {
	for {
		if lb.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(lb.rate):
			// 继续尝试
		}
	}
}

func main() {
	// 创建一个漏桶：容量为3，每秒漏出1个
	bucket := NewLeakyBucket(3, time.Second)
	defer bucket.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// 模拟10个请求
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 等待直到允许通过
			if err := bucket.Wait(ctx); err != nil {
				fmt.Printf("Request %d failed: %v\n", id, err)
				return
			}

			fmt.Printf("Request %d allowed at %v\n", id, time.Now().Format("15:04:05.000"))
		}(i)
	}

	wg.Wait()
}
