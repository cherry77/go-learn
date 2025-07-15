package mode

/**
## **4. 定时任务（Ticker）题目**
### **题目 4：周期性爬取网页**
使用 `time.Ticker` 实现 **每 3 秒爬取一次网页**，并打印当前时间。
- 要求：
  - 使用 `time.NewTicker` 控制任务节奏。
  - 在 `10` 秒后停止任务（用 `context` 或 `time.After`）。
*/
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func crawl(url string) {
	fmt.Printf("Crawling at %v\n", time.Now())
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Fetched %s at %v, status: %s\n", url, time.Now().Format("15:04:05.000"), resp.Status)
}

func tickerCrawl(url string) {
	ticker := time.NewTicker(3 * time.Second) // 每3秒触发一次
	defer ticker.Stop()                       // 确保退出时停止Ticker

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 确保退出时取消上下文

	for {
		select {
		case <-ctx.Done(): // 10秒后停止 上下文结束，退出协程
			return
		case t := <-ticker.C: // 每次Ticker触发
			crawl(url)
			fmt.Println("Tick at", t)
		}
	}
}

func tickerCrawlGoRoutine(url string) {
	ticker := time.NewTicker(3 * time.Second) // 每3秒触发一次
	defer ticker.Stop()                       // 确保退出时停止Ticker

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 确保退出时取消上下文

	go func() {
		for {
			select {
			case <-ctx.Done(): // 10秒后停止 上下文结束，退出协程
				return
			case t := <-ticker.C: // 每次Ticker触发
				crawl(url)
				fmt.Println("Tick at", t)
			}
		}
	}()

	// 等待上下文结束
	<-ctx.Done()
	fmt.Println("Ticker stopped")
}

func timeAfterCrawl(url string) {
	ticker := time.NewTicker(3 * time.Second) // 每3秒触发一次
	defer ticker.Stop()                       // 确保退出时停止Ticker

	timeout := time.After(10 * time.Second) // 10秒后停止

	for {
		select {
		case <-timeout: // 10秒后停止
			fmt.Println("Time after stopped")
			return
		case t := <-ticker.C: // 每次Ticker触发
			crawl(url)
			fmt.Println("Tick at", t)
		}
	}
}

// 对于长期运行的服务，我会这样改进goroutine版本：
func StartPeriodicCrawl(url string, interval time.Duration, timeout time.Duration) context.CancelFunc {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				crawl(url)
				log.Printf("Crawled at %v", t)
			}
		}
	}()

	return cancel // 返回cancel函数让调用者可以提前停止
}
