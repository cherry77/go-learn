package main

import "fmt"

// select 语句与 default 分支的结合使用，用于实现非阻塞的通道操作
// 核心概念
// 非阻塞通道操作：
// 通过 select + default 实现：当通道操作无法立即完成时，直接执行 default 分支。
//
// 通道状态：
// 无缓冲通道（make(chan string)）的发送和接收会阻塞，直到另一端准备好。
//
// select 行为：
// 当所有 case 的通道操作均阻塞时，执行 default。

// ### 类比现实场景
// - **第一部分**：像查看信箱时发现没有信，直接离开（非阻塞检查）。
// - **第二部分**：像尝试快速投递包裹，但发现收件人不在家，直接放弃（非阻塞发送）。
// - **第三部分**：像同时监听电话和门铃，但两者均无动静，转而做其他事（多通道监听）。
func main() {
	messages := make(chan string)
	signals := make(chan bool)

	// 非阻塞接收（从 messages 接收）
	select {
	case msg := <-messages: // 尝试从 messages 接收
		fmt.Println("received message", msg)
	default: // 无数据时立即执行
		fmt.Println("no message received")
	}

	// 非阻塞发送（向 messages 发送）
	// messages 是无缓冲通道，且没有协程在接收，发送操作会阻塞，因此触发 default
	msg := "hi"
	select {
	case messages <- msg: // 尝试向 messages 发送 "hi"
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent") // 无接收者时立即执行
	}

	// messages 和 signals 都无数据，直接触发 default
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

//### 关键点总结
//1. **无缓冲通道的特性**：
//- 发送和接收必须同时准备好，否则操作会阻塞。
//- 通过 `select` + `default` 可以避免阻塞，实现“尝试性”操作。
//
//2. **`default` 的作用**：
//- 当所有 `case` 的通道操作无法立即完成时，执行 `default` 分支。
//- 若省略 `default`，`select` 会一直阻塞，直到某个 `case` 就绪。
//
//3. **实际应用场景**：
//- 检查通道是否有数据（非阻塞读）。
//- 尝试发送数据而不阻塞（非阻塞写）。
//- 同时监听多个通道，避免死锁。

//### 扩展思考
//- **若 `messages` 是缓冲通道**：
//```go
//  messages := make(chan string, 1)  // 缓冲大小为1
//  ```
//第二部分会成功发送（`sent message hi`），因为缓冲允许暂存数据。
//- **若移除 `default`**：
//`select` 会一直阻塞，直到某个 `case` 就绪（可能引发死锁）。
