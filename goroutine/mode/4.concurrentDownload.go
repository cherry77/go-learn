package mode

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
*
## **5. 综合应用题**
### **题目 5：并发下载 + 限速 + Worker Pool**
实现一个 **并发下载器**，要求：
1. 使用 **Worker Pool**（`3` 个 Worker）并发下载多个文件。
2. 使用 **速率限制**（每秒最多 `2` 个下载请求）。
3. 使用 `sync.WaitGroup` 等待所有下载完成。
*/
func download(url string) {
	fmt.Printf("Downloading %s...\n", url)
	time.Sleep(1 * time.Second) // 模拟下载耗时
}

func concurrentDownload(urls []string, qps int, burst int, worker int) {
	urlChan := make(chan string)
	var wg sync.WaitGroup

	// 速率限制器（令牌桶）
	burstLimiter := make(chan struct{}, burst)
	// 初始填充令牌
	for i := 0; i < qps; i++ {
		burstLimiter <- struct{}{}
	}

	// 令牌补充 goroutine（带停止机制）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(qps))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				select {
				case burstLimiter <- struct{}{}: // 尝试添加令牌
				default: // 桶满时丢弃令牌
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChan {
				<-burstLimiter // 限制速率
				download(url)
			}
		}()
	}

	go func() {
		for _, url := range urls {
			urlChan <- url
		}
		close(urlChan)
	}()

	wg.Wait()
	fmt.Println("All downloads completed!")
}
