package mode

import (
	"context"
	"fmt"
	"time"
)

/*
*
## **3. 速率限制（Rate Limiting）题目**
### **题目 3：API 限流（令牌桶）**
实现一个 **速率限制器**，限制每秒最多调用 `2` 次 `mockAPICall()`。
- 要求：
  - 使用 `chan time.Time` + `time.Tick` 实现令牌桶。
  - 初始允许 `2` 次突发请求，之后每 `500ms` 补充 `1` 个令牌。
*/
func mockAPICall(id int) {
	fmt.Printf("API call %d at %v\n", id, time.Now())
}

func burstLimit(qps int, burst int) {
	burstyLimiter := make(chan time.Time, burst)
	for i := 0; i < burst; i++ {
		burstyLimiter <- time.Now() // 初始突发令牌
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保退出时停止goroutine

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(qps)) // 每秒补充令牌
		defer ticker.Stop()                                        // 确保退出时停止Ticker

		for {
			select {
			case t := <-ticker.C:
				select {
				case burstyLimiter <- t:
					fmt.Printf("Token added at %v\n", t.Format("15:04:05.000"))
				default:
					fmt.Println("Token discarded (bucket full)")
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// 模拟调用时可以加个延迟，让输出更易观察
	for i := 0; i < 10; i++ { // 模拟10次API调用
		<-burstyLimiter // 等待令牌
		mockAPICall(i + 1)
		time.Sleep(200 * time.Millisecond) // 添加延迟
	}
}
