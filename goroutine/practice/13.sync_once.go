package main

import (
	"fmt"
	"sync"
)

//13. Once 使用
//使用 sync.Once 确保某个初始化操作（如配置文件加载）在并发环境下只执行一次。

func loadConfig() {
	fmt.Println("Loading config...")
}

var once sync.Once

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d is running\n", id)

			once.Do(loadConfig) // 确保 loadConfig 只被调用一次

			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}

	wg.Wait()
}
