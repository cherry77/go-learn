package main

import (
	"fmt"
	"sync"
)

// Counter 计数器服务
type Counter struct {
	value     int            // 实际计数值
	opsChan   chan func()    // 操作通道
	closeChan chan struct{}  // 关闭信号
	wg        sync.WaitGroup // 用于优雅关闭
}

// NewCounter 创建新计数器
func NewCounter() *Counter {
	c := &Counter{
		opsChan:   make(chan func()),
		closeChan: make(chan struct{}),
	}

	// 启动状态管理goroutine
	c.wg.Add(1)
	go c.loop()

	return c
}

// loop 状态管理主循环
func (c *Counter) loop() {
	defer c.wg.Done()

	for {
		select {
		case op := <-c.opsChan:
			op() // 执行操作
		case <-c.closeChan:
			return // 退出goroutine
		}
	}
}

// Add 增加计数值
func (c *Counter) Add(n int) {
	// 使用通道发送操作请求
	done := make(chan struct{})
	c.opsChan <- func() {
		c.value += n
		close(done)
	}
	<-done // 等待操作完成
}

// Get 获取当前计数值
func (c *Counter) Get() int {
	// 使用通道发送查询请求
	result := make(chan int)
	c.opsChan <- func() {
		result <- c.value
	}
	return <-result
}

// Close 关闭计数器，释放资源
func (c *Counter) Close() {
	close(c.closeChan)
	c.wg.Wait()
}

func main() {
	counter := NewCounter()
	defer counter.Close()

	var wg sync.WaitGroup

	// 启动多个goroutine并发增加计数器
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Add(1)
			}
		}()
	}

	wg.Wait()
	fmt.Println("Final counter value:", counter.Get()) // 应该输出10000
}
