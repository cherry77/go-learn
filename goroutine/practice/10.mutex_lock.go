package main

import "sync"

//10. 互斥锁使用
//创建一个共享计数器，多个 goroutine 并发地增加计数器值，使用 sync.Mutex 保证线程安全。

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()         // 获取锁
	defer c.mu.Unlock() // 确保锁会被释放
	c.count++           // 安全地修改共享变量
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	counter := SafeCounter{}

	var wg sync.WaitGroup

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()

	println(counter.Value())

}
