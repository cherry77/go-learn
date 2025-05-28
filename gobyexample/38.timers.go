package main

import (
	"fmt"
	"time"
)

func main() {
	timer1 := time.NewTimer(2 * time.Second) // // 2秒后触发 time.NewTimer(d Duration) 创建一个定时器，在 d 时间后向 timer.C 通道发送当前时间。
	<-timer1.C                               // 阻塞直到2秒后 阻塞当前协程，直到定时器触发（通道收到时间值）。
	fmt.Println("Timer 1 fired")

	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()
	stop2 := timer2.Stop() // // 尝试停止定时器 停止定时器，防止其触发。如果定时器已停止，返回 true；如果已触发，返回 false。
	if stop2 {
		fmt.Println("Timer 2 stopped") // 成功停止时执行
	}

	time.Sleep(2 * time.Second)
}

/*### 关键点总结
1. **定时器的本质**：
- `Timer` 通过通道 `C` 通知触发，适合单次延迟任务。
- 对比 `time.Sleep`：`Sleep` 是纯阻塞，而 `Timer` 可停止或重置。

2. **停止定时器的时机**：
- 必须在定时器触发前调用 `Stop()` 才有效。
- 如果定时器已触发或已停止，`Stop()` 返回 `false`。

3. **资源释放**：
- 未被停止的定时器会正常触发，但需要确保通道被读取（否则可能泄漏资源）。

---

### 类比现实场景
- **`Timer1`**：像设置一个2秒的厨房定时器，时间到后响铃（`Timer 1 fired`）。
- **`Timer2`**：像设置1秒的备用定时器，但你在它响铃前按下了取消按钮（`Timer 2 stopped`）。
- **协程阻塞**：像备用定时器的铃铛永远等不到响的那一刻（协程中的 `<-timer2.C` 一直阻塞）。

---

### 实际应用场景
- **任务超时控制**：结合 `select` 实现操作限时。
- **延迟操作**：如游戏中的技能冷却。
- **资源清理**：延迟关闭文件或连接。

---

### 扩展思考
- **若移除 `timer2.Stop()`**：
输出会多一行 `"Timer 2 fired"`（协程中的打印生效）。
- **使用 `time.After` 替代**：
```go
  select {
  case <-time.After(2 * time.Second):
      fmt.Println("Timeout")
  }
  ```
`time.After` 更适合简单的超时场景，但无法主动停止。*/
