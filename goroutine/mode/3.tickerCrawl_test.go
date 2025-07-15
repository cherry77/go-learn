package mode

import (
	"log"
	"testing"
	"time"
)

func Test3_TickerCrawl(t *testing.T) {
	tickerCrawl("https://example.com") // 每3秒爬取一次网页，10秒后停止任务
}

func Test3_TickerCrawlGoRoutine(t *testing.T) {
	tickerCrawlGoRoutine("https://example.com") // 每3秒爬取一次网页，10秒后停止任务
}

func Test3_TimeAfterCrawl(t *testing.T) {
	timeAfterCrawl("https://example.com") // 每3秒爬取一次网页，10秒后停止任务
}

func Test3_StartPeriodicCrawl(t *testing.T) {
	log.Println("Starting program...")

	// 启动两个爬虫任务
	cancel1 := StartPeriodicCrawl("https://www.example.com", 2*time.Second, 8*time.Second)
	StartPeriodicCrawl("https://www.google.com", 3*time.Second, 12*time.Second)

	// 模拟程序运行...
	time.Sleep(5 * time.Second)

	// 提前取消第一个爬虫（原本应该运行8秒）
	log.Println("Manually stopping first crawler")
	cancel1()

	// 继续等待
	time.Sleep(10 * time.Second)
	log.Println("Main program exiting")
}
