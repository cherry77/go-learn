package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// SharedData 包含读写锁保护的数据
type SharedData struct {
	mu      sync.RWMutex
	content string
	version int
}

// Reader 并发读取数据
func (d *SharedData) Reader(id int) {
	d.mu.RLock() // 获取读锁
	defer d.mu.RUnlock()

	fmt.Printf("Reader %d: 版本=%d, 内容=%s\n", id, d.version, d.content)
	time.Sleep(time.Duration(rand.Intn(100))) // 模拟读取耗时
}

// Writer 更新数据
func (d *SharedData) Writer(id int) {
	d.mu.Lock() // 获取写锁
	defer d.mu.Unlock()

	d.version++
	d.content = fmt.Sprintf("Writer %d 更新于 %v", id, time.Now())
	fmt.Printf("Writer %d: 新版本=%d\n", id, d.version)
	time.Sleep(time.Duration(100 + rand.Intn(100))) // 模拟写入耗时
}

func main() {
	rand.Seed(time.Now().UnixNano())
	data := SharedData{
		content: "初始内容",
		version: 1,
	}

	var wg sync.WaitGroup

	// 启动5个reader goroutine
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ { // 内部循环让每个 goroutine 执行 多次工作单元，比"一任务一goroutine"更符合生产环境实际用法
				data.Reader(id)
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			}
		}(i)
	}

	// 启动2个writer goroutine
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 2; j++ {
				data.Writer(id)
				time.Sleep(time.Duration(300+rand.Intn(200)) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("最终数据版本:", data.version)
}
