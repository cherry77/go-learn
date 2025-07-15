package mode

import (
	"fmt"
	"net/http"
	"sync"
)

/*
## **1. Worker Pool 题目**
### **题目 1：批量处理 HTTP 请求**
实现一个 **Worker Pool**，并发发送 HTTP 请求，但最多同时运行 `5` 个 Goroutine。
- 输入：`urls []string`（待请求的 URL 列表）。
- 要求：
- 使用 `chan` + `WaitGroup` 实现 Worker Pool。
- 每个 Worker 从 Channel 读取 URL 并发送 HTTP GET 请求。
- 打印响应状态码（如 `200`）或错误信息。
*/

// fetchURL 发送HTTP GET请求并返回状态码
func fetchURL(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// worker 从channel读取URL并处理
func worker(urlChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urlChan {
		status, err := fetchURL(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", url, err)
		} else {
			fmt.Printf("%s -> %d\n", url, status)
		}
	}
}
func processURLs(urls []string, workerCount int) {
	urlChan := make(chan string)
	var wg sync.WaitGroup

	// 先启动 worker!!!!
	// 启动worker (这里设置为5个并发worker)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(urlChan, &wg)
	}

	// 然后发送URL到channel
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan) // 关闭channel通知worker没有更多任务

	wg.Wait() // 等待所有worker完成
}

/**
  如果将 for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)
挪到 worker 启动之前，会导致Channel 阻塞问题：
1. 在主 Goroutine 中直接向无缓冲 channel urlChan 发送数据，但没有接收者准备好
2. 这会导致死锁，因为无缓冲 channel 需要发送和接收同时准备好

关键修正点
1. 启动顺序调整：
先启动所有 worker Goroutine
然后再向 channel 发送数据

2. Channel 使用改进：
保持无缓冲 channel，但确保接收者先就绪
发送完成后正确关闭 channel
*/
