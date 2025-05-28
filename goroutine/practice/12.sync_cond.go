package main

import (
	"fmt"
	"sync"
	"time"
)

//12. 条件变量
//使用 sync.Cond 实现一个生产者-消费者模型，当队列满时生产者等待，队列空时消费者等待。

type Queue struct {
	items    []int
	capacity int
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

func NewQueue(capacity int) *Queue {
	q := &Queue{
		items:    make([]int, 0, capacity),
		capacity: capacity,
	}
	q.notFull = sync.NewCond(&q.mu)
	q.notEmpty = sync.NewCond(&q.mu)
	return q
}

// 生产者方法
func (q *Queue) Produce(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 队列满时等待
	for len(q.items) == q.capacity {
		fmt.Println("队列已满，生产者等待")
		q.notFull.Wait()
	}

	// 添加元素并通知消费者
	q.items = append(q.items, item)
	fmt.Printf("生产: %d (队列长度: %d)\n", item, len(q.items))
	q.notEmpty.Signal() // 通知一个等待的消费者
}

// 消费者方法
func (q *Queue) Consume() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 队列空时等待
	for len(q.items) == 0 {
		fmt.Println("队列为空，消费者等待")
		q.notEmpty.Wait()
	}

	// 取出元素并通知生产者
	item := q.items[0]
	q.items = q.items[1:]
	fmt.Printf("消费: %d (队列长度: %d)\n", item, len(q.items))
	q.notFull.Signal() // 通知一个等待的生产者
	return item
}

func main() {
	queue := NewQueue(3) // 容量为3的队列

	var wg sync.WaitGroup

	// 启动2个消费者
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				item := queue.Consume()
				time.Sleep(time.Duration(200+id*100) * time.Millisecond) // 模拟处理时间
				fmt.Printf("消费者%d处理完成: %d\n", id, item)
			}
		}(i)
	}

	// 启动3个生产者
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 4; j++ {
				item := id*100 + j
				queue.Produce(item)
				time.Sleep(time.Duration(100+id*50) * time.Millisecond) // 模拟生产间隔
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("所有生产消费完成")
}
