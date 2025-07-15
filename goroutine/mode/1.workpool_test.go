package mode

import "testing"

func Test1_WorkerPool(t *testing.T) {
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.medium.com",
		"https://www.youtube.com",
		"https://www.amazon.com",
		"https://invalid.url", // 这个会出错
	}

	workerCount := 5 // 最大并发 Goroutine 数
	processURLs(urls, workerCount)
}
