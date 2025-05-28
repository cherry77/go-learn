package main

import (
	"fmt"
	"time"
)

/*
*
select 语句与 time.After 超时机制的结合使用，用于处理通道操作的超时控制
代码核心功能
select 语句：监听多个通道操作，执行第一个就绪的 case。

time.After：返回一个通道，在指定时间后发送当前时间（用于超时控制）。

缓冲通道：make(chan string, 1) 允许发送操作不阻塞（缓冲大小为1）。
*/
func main() {
	c1 := make(chan string, 1)

	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "hello"
	}()
	select {
	case msg := <-c1:
		fmt.Println("Received:", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout after 1 second")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}

//关键机制说明
//select 的工作逻辑：
//
//同时监听所有 case 的通道操作（接收或发送）。
//
//执行 第一个就绪的 case，其他 case 被忽略。
//
//如果多个 case 同时就绪，随机选择一个执行。
//
//time.After 的超时控制：
//
//返回一个单向通道（<-chan time.Time）。
//
//在指定时间后，该通道会自动收到一个时间值（类似定时器）。
