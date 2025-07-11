package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

// 客户填写查询申请表（readOp 结构体）：
type readOp struct {
	key  int
	resp chan int
}

// 办理存取款（写入操作）
type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func main() {
	var readOps uint64
	var writeOps uint64

	reads := make(chan readOp)
	writes := make(chan writeOp)

	//保险柜管理员: 只通过两个窗口与外界沟通：一个接收查询请求，一个接收修改请求
	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.resp <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	// 100个只查询余额的客户（读取 goroutines）
	for range 100 { // 创建 100 个 goroutine
		go func() {
			for {
				read := readOp{key: rand.Intn(5), resp: make(chan int)}
				reads <- read
				<-read.resp // 阻塞等待：<-read.resp 会阻塞当前 goroutine，直到状态管理 goroutine 通过这个通道发送响应

				//银行记录这笔业务（atomic.AddUint64(&writeOps, 1)）
				atomic.AddUint64(&readOps, 1) // 第一个参数：要修改的变量的指针 (&readOps)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	//10个需要存款/取款的客户（写入 goroutines）
	for range 10 {
		go func() {
			for {
				write := writeOp{key: rand.Intn(5), val: rand.Intn(100), resp: make(chan bool)}
				writes <- write
				<-write.resp // 阻塞等待：<-write.resp 会阻塞当前 goroutine，直到状态管理 goroutine 通过这个通道发送响应

				//银行记录这笔业务（atomic.AddUint64(&writeOps, 1)）
				atomic.AddUint64(&writeOps, 1) // 第一个参数：要修改的变量的指针 (&writeOps)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)
	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}
