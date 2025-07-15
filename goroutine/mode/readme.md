### **并发模式练习题（Worker Pool、信号量、速率限制、定时任务）**
以下题目帮助你掌握 **Go 并发模式**，涵盖 **Worker Pool、信号量、速率限制（令牌桶）、定时任务** 等场景。

---

## **1. Worker Pool 题目**
### **题目 1：批量处理 HTTP 请求**
实现一个 **Worker Pool**，并发发送 HTTP 请求，但最多同时运行 `5` 个 Goroutine。
- 输入：`urls []string`（待请求的 URL 列表）。
- 要求：
    - 使用 `chan` + `WaitGroup` 实现 Worker Pool。
    - 每个 Worker 从 Channel 读取 URL 并发送 HTTP GET 请求。
    - 打印响应状态码（如 `200`）或错误信息。

**示例代码框架：**
```go
func fetchURL(url string) (int, error) {
    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()
    return resp.StatusCode, nil
}

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

func main() {
    urls := []string{"https://google.com", "https://github.com", ...}
    // TODO: 实现 Worker Pool
}
```

---

## **2. 信号量（Semaphore）题目**
### **题目 2：限制文件并发写入**
使用 **信号量模式** 控制最多 `3` 个 Goroutine 同时写入文件。
- 输入：`data []string`（待写入的字符串列表）。
- 要求：
    - 每个 Goroutine 写入一个文件（如 `file_1.txt`, `file_2.txt`）。
    - 使用 `chan struct{}` 作为信号量，限制并发数。

**示例代码框架：**
```go
func writeToFile(filename, content string) {
    // 模拟写入延迟
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Written to %s\n", filename)
}

func main() {
    data := []string{"data1", "data2", "data3", "data4", "data5"}
    sem := make(chan struct{}, 3) // 允许 3 个并发
    var wg sync.WaitGroup

    for i, d := range data {
        wg.Add(1)
        go func(id int, content string) {
            defer wg.Done()
            sem <- struct{}{} // 获取信号量
            defer func() { <-sem }()
            writeToFile(fmt.Sprintf("file_%d.txt", id), content)
        }(i, d)
    }
    wg.Wait()
}
```

---

## **3. 速率限制（Rate Limiting）题目**
### **题目 3：API 限流（令牌桶）**
实现一个 **速率限制器**，限制每秒最多调用 `2` 次 `mockAPICall()`。
- 要求：
    - 使用 `chan time.Time` + `time.Tick` 实现令牌桶。
    - 初始允许 `2` 次突发请求，之后每 `500ms` 补充 `1` 个令牌。

**示例代码框架：**
```go
func mockAPICall(id int) {
    fmt.Printf("API call %d at %v\n", id, time.Now())
}

func main() {
    burstyLimiter := make(chan time.Time, 2) // 初始 2 个令牌
    for i := 0; i < 2; i++ {
        burstyLimiter <- time.Now()
    }

    go func() {
        for t := range time.Tick(500 * time.Millisecond) { // 每 500ms 补充 1 个
            burstyLimiter <- t
        }
    }()

    for i := 1; i <= 5; i++ {
        <-burstyLimiter // 等待令牌
        go mockAPICall(i)
    }
    time.Sleep(2 * time.Second) // 等待所有请求完成
}
```

---

## **4. 定时任务（Ticker）题目**
### **题目 4：周期性爬取网页**
使用 `time.Ticker` 实现 **每 3 秒爬取一次网页**，并打印当前时间。
- 要求：
    - 使用 `time.NewTicker` 控制任务节奏。
    - 在 `10` 秒后停止任务（用 `context` 或 `time.After`）。

**示例代码框架：**
```go
func crawl() {
    fmt.Printf("Crawling at %v\n", time.Now())
}

func main() {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()
    done := time.After(10 * time.Second) // 10 秒后停止

    for {
        select {
        case <-ticker.C:
            crawl()
        case <-done:
            fmt.Println("Stopped!")
            return
        }
    }
}
```

---

## **5. 综合应用题**
### **题目 5：并发下载 + 限速 + Worker Pool**
实现一个 **并发下载器**，要求：
1. 使用 **Worker Pool**（`3` 个 Worker）并发下载多个文件。
2. 使用 **速率限制**（每秒最多 `2` 个下载请求）。
3. 使用 `sync.WaitGroup` 等待所有下载完成。

**提示：**
- 结合 `chan`（任务队列） + `time.Tick`（限速） + `WaitGroup`（同步）。
- 模拟下载函数：
  ```go
  func download(url string) {
      fmt.Printf("Downloading %s...\n", url)
      time.Sleep(1 * time.Second) // 模拟下载耗时
  }
  ```

---

## **答案 & 思路**
### **关键点总结**
1. **Worker Pool**
    - 固定数量的 Goroutine + 任务 Channel。
    - 适合长期运行的任务（如 HTTP 服务）。

2. **信号量模式**
    - `chan struct{}` 控制并发数。
    - 适合短期任务 + 资源保护（如文件写入）。

3. **速率限制**
    - `chan time.Time` + `time.Tick` 实现令牌桶。
    - 适合 API 调用限流。

4. **定时任务**
    - `time.Ticker` 控制任务节奏。
    - 适合心跳检测、周期性任务。

通过练习这些题目，你可以掌握 Go 并发编程的核心模式！ 🚀
