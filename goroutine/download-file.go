package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// DownloadResult 表示下载结果
type DownloadResult struct {
	URL      string
	Success  bool
	Duration time.Duration
	Size     int64 // 文件大小(字节)
	Error    error
}

// 模拟下载文件
func downloadFile(url string, progress chan<- int, result chan<- DownloadResult) {
	start := time.Now()
	defer func() {
		close(progress) // 关闭进度通道
	}()

	// 随机生成文件大小(1MB~10MB)
	size := rand.Int63n(10*1024*1024) + 1024*1024

	// 模拟下载过程
	for i := 0; i <= 100; i++ {
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond) // 模拟下载延迟
		progress <- i                                               // 发送进度
	}

	// 随机决定成功或失败(80%成功率)
	var err error
	if rand.Float32() < 0.2 {
		err = fmt.Errorf("download failed: connection timeout")
	}

	result <- DownloadResult{
		URL:      url,
		Success:  err == nil,
		Duration: time.Since(start),
		Size:     size,
		Error:    err,
	}
}

func main() {
	// 模拟要下载的文件URL
	urls := []string{
		"http://example.com/file1.zip",
		"http://example.com/file2.pdf",
		"http://example.com/file3.mp4",
		"http://example.com/file4.iso",
		"http://example.com/file5.exe",
	}

	var wg sync.WaitGroup
	resultChan := make(chan DownloadResult, len(urls))
	startTime := time.Now()

	// 启动下载goroutines
	for _, url := range urls {
		wg.Add(1)
		progressChan := make(chan int)

		// 启动进度显示器
		go func(url string, progress <-chan int) {
			for p := range progress {
				fmt.Printf("\rDownloading %s: %d%%", url, p)
			}
			fmt.Println() // 换行
		}(url, progressChan)

		// 启动下载器
		go func(url string) {
			defer wg.Done()
			downloadFile(url, progressChan, resultChan)
		}(url)
	}

	// 等待所有下载完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	var successCount, failCount int
	var totalBytes int64
	var totalDuration time.Duration

	for result := range resultChan {
		if result.Success {
			successCount++
			totalBytes += result.Size
			totalDuration += result.Duration
			fmt.Printf("Download succeeded: %s (Size: %.2f MB, Time: %v)\n",
				result.URL, float64(result.Size)/(1024*1024), result.Duration)
		} else {
			failCount++
			fmt.Printf("Download failed: %s (Error: %v)\n", result.URL, result.Error)
		}
	}

	// 计算统计信息
	totalTime := time.Since(startTime)
	averageSpeed := float64(totalBytes) / totalTime.Seconds() / (1024 * 1024) // MB/s

	// 显示统计信息
	fmt.Println("\n===== Download Statistics =====")
	fmt.Printf("Total time: %v\n", totalTime)
	fmt.Printf("Successfully downloaded: %d files\n", successCount)
	fmt.Printf("Failed downloads: %d files\n", failCount)
	fmt.Printf("Average download speed: %.2f MB/s\n", averageSpeed)
}
