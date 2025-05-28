package main

import (
	"fmt"
	"time"
)

// 周期性定时器 time.Ticker 的使用，以及如何安全停止 Ticker 和协程
func main() {
	ticker := time.NewTicker(500 * time.Millisecond) // 每500ms触发一次
	done := make(chan bool)                          // 控制协程退出的信号通道

	go func() {
		for {
			select {
			case <-done: // 收到退出信号
				return // 结束协程
			case t := <-ticker.C: // 每次Ticker触发
				fmt.Println("Tick at", t)
			}
		}
	}()
	time.Sleep(1600 * time.Millisecond) // 让Ticker触发约3次（500ms×3=1500ms）
	ticker.Stop()                       // 停止Ticker（不再发送事件）
	done <- true                        // 通知协程退出
	fmt.Println("Ticker stopped")
}

/*### 关键点总结
1. **必须调用 `ticker.Stop()`**：
- 否则 Ticker 会持续占用资源（即使程序其他部分已结束）。
2. **通道关闭的替代方案**：
- 可直接 `close(done)`，协程中通过 `case <-done:` 检测到零值退出。
3. **定时精度**：
- Ticker 会尽力维持间隔，但受系统调度影响可能有微小偏差。

---

### 类比现实场景
- **Ticker**：像学校的上课铃，每隔45分钟响一次。
- **`done` 通道**：像校长广播“今天提前放学”，铃声停止，学生（协程）回家。
- **`ticker.Stop()`**：像物理关闭电铃电源，确保不会再响。

---

### 实际应用场景
- 定时数据采集（如每5秒读取传感器）。
- 周期性日志记录。
- 游戏中的自动保存功能。

---

### 扩展思考
- **若移除 `ticker.Stop()`**：
Ticker 持续触发，但协程因 `done` 信号会退出，导致 `ticker.C` 无接收者（资源泄漏）。
- **若移除 `done` 通道**：
协程无法退出，即使 Ticker 已停止，`for-select` 会永久阻塞（协程泄漏）。*/
