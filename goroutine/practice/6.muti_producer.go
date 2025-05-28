package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 多生产者单消费者
// 创建 3 个生产者 goroutine 和一个消费者 goroutine。生产者生成随机数发送到通道，消费者接收并打印这些数字。
func main() {
	ch := make(chan int, 10) // 关键：带缓冲!!!!

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go producer(i+1, ch, &wg)
	}
	wg.Wait()

	go consumer(ch)

	close(ch)

	time.Sleep(500 * time.Millisecond)
	fmt.Println("程序结束")
}

func producer(id int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	num := rand.Int()
	ch <- num
	println("Producer", id, "produced:", num)

}

func consumer(ch <-chan int) {
	for num := range ch {
		println(num)
		time.Sleep(300 * time.Millisecond) // 模拟处理时间
	}
}
