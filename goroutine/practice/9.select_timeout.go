package main

import "time"

//9. 超时控制
//编写一个程序，goroutine 执行一个耗时操作（如 sleep 2秒），主程序使用 select 实现超时控制，如果超过 1 秒就打印 "Timeout"。

func main() {
	resultCh := make(chan string)
	go func() {
		// 模拟耗时操作
		time.Sleep(2 * time.Second)
		resultCh <- "操作完成"
	}()

	select {
	case result := <-resultCh:
		println(result)
	case <-time.After(1 * time.Second):
		println("超时")
		return
	}
}

// 总结：select 是"单次"操作，执行一个 case 后就结束
