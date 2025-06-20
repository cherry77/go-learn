package main

import "fmt"

/*
### 代码功能概述
1. **主协程**：
- 创建一个缓冲通道 `jobs`（容量5）和一个信号通道 `done`。
- 向 `jobs` 发送3个任务（1, 2, 3）。
- 关闭 `jobs` 通道，表示任务发送完毕。
- 通过 `<-done` 等待工作协程完成。

2. **工作协程**：
- 持续从 `jobs` 接收任务，直到通道关闭。
- 通过 `more` 判断通道是否关闭。
- 所有任务处理完成后，向 `done` 发送信号。
*/
func main() {
	jobs := make(chan int, 5)
	done := make(chan bool) // 无缓冲通道, 用于工作协程通知主协程任务处理完成。

	go func() {
		for {
			job, more := <-jobs // 循环接收任务
			if more {
				fmt.Println("received job", job)
			} else {
				fmt.Println("received all jobs")
				done <- true // 通知主协程
				return       // 退出协程
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs) // 通知接收方不再有新任务发送。已关闭的通道仍可读取剩余数据，读完后再读取会返回零值和 false。
	fmt.Println("sent all jobs")

	<-done // 阻塞直到收到信号
}

/*### 关键点总结
1. **通道关闭的最佳实践**：
- 由发送方关闭通道（避免接收方关闭引发 panic）。
- 关闭通道是一种广播机制，通知接收方数据流结束。

2. **缓冲通道的作用**：
- 容量为5的缓冲允许主协程快速发送任务，无需等待工作协程立即接收。

3. **协程同步模式**：
- 通过 `done` 通道实现“等待子协程完成”的同步，这是 Go 中常见的模式。

---

### 类比现实场景
- **`jobs` 通道**：像工厂的任务队列，主协程是生产部（投递任务），工作协程是工人（处理任务）。
- **`close(jobs)`**：生产部宣布“任务全部下发”，工人处理完队列后下班。
- **`done` 信号**：工人下班前打卡，通知工厂可以关门。

---

### 扩展思考
- **若移除 `close(jobs)`**：
工作协程会一直阻塞在 `<-jobs`，导致 `done` 信号无法发送（死锁）。
- **若通道有剩余数据**：
关闭后仍可读取剩余数据，直到通道为空时 `more` 返回 `false`。*/
