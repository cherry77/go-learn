package main

import (
	"fmt"
	"time"
)

func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second) // 模拟耗时任务（1秒）
	fmt.Println("done")

	done <- true // 发送完成信号
}

//类比现实场景
//主协程：像项目经理，分配任务后等待员工汇报。
//
//工作协程：像员工，完成任务后通过通道（类似邮件）通知经理。
//
//<-done：经理必须收到邮件才能继续下一步。

//扩展思考
//若移除 <-done，主协程可能直接退出（不等待 worker 完成）。
//
//若通道无缓冲（make(chan bool)），工作协程的 done <- true 会阻塞，直到主协程准备好接收。

func main() {
	done := make(chan bool, 1) // 创建缓冲为1的布尔通道

	go worker(done) // 启动工作协程,go worker(done) 启动一个并发执行的 worker 协程，不会阻塞主协程。

	<-done // 阻塞主协程，直到收到完成信号,<-done 会使主协程阻塞，直到从通道接收到数据（即 worker 协程发送的 true）。
}
