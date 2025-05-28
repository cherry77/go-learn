package main

import (
	"fmt"
	"time"
)

/*
  - id：Worker 的唯一标识。
    jobs：只读通道，接收任务（类型 int）。
    results：只写通道，发送结果（类型 int）。
*/
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs { // 自动从jobs通道接收任务
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2 // 将结果发送到results通道
	}
}

// 实现了一个工作池（Worker Pool）模式，展示了如何用多个协程（Worker）并发处理任务，并通过通道进行任务分发和结果收集
// 3 个 Worker 并发处理 5 个任务
func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)    // 缓冲通道（容量=任务数）
	results := make(chan int, numJobs) // 缓冲通道（容量=任务数）

	// 启动3个Worker协程
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	// 发送5个任务到jobs通道
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs) // 关闭通道（通知Worker无新任务）

	for a := 1; a <= numJobs; a++ {
		<-results // 从results通道接收结果
	}
}

/*### 关键点总结
1. **Worker 数量与任务分配**：
- Worker 数（3） < 任务数（5）时，Worker 会复用处理多个任务。
- 如果 Worker 数 ≥ 任务数，部分 Worker 可能闲置。

2. **通道关闭的最佳实践**：
- 由发送方（主协程）关闭通道，避免向已关闭通道发送数据引发 panic。

3. **资源释放**：
- Worker 协程在 `jobs` 通道关闭后自动退出，无协程泄漏。

---

### 类比现实场景
- **`jobs` 通道**：像工厂的任务公告板，工人（Worker）主动领取任务。
- **`results` 通道**：像工人完成任务后把产品放到传送带上。
- **主协程**：像经理，分配完任务后等待所有产品质检完成。

---

### 实际应用场景
- 批量处理文件（如并发解析日志）。
- 高并发网络请求（如爬虫任务分发）。
- 计算密集型任务的分片处理。

---

### 扩展思考
- **动态调整 Worker 数量**：
可通过 `sync.WaitGroup` 实现 Worker 的弹性扩缩容。
- **错误处理**：
增加错误通道（`errChan`）收集 Worker 的处理异常。
- **任务优先级**：
使用多个 `jobs` 通道或优先级队列实现差异化调度。*/
