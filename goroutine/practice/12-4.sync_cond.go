package main

import (
	"fmt"
	"sync"
	"time"
)

//在 12-2 的基础上增加 方式二
// 任务1：增加紧急菜品通道

type Queue struct {
	normalItems   []int
	priorityItems []int
	capacity      int
	mu            sync.Mutex
	notFull       *sync.Cond // 队列未满条件
	notEmpty      *sync.Cond // 队列非空条件（普通+紧急）
	hasPriority   *sync.Cond // 紧急专用条件
}

func NewQueue(capacity int) *Queue {
	q := &Queue{capacity: capacity}
	q.notFull = sync.NewCond(&q.mu)
	q.notEmpty = sync.NewCond(&q.mu)
	q.hasPriority = sync.NewCond(&q.mu)
	return q
}

// 生产普通菜品
func (q *Queue) ProduceNormal(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.normalItems)+len(q.priorityItems) == q.capacity {
		fmt.Println("队列满，普通生产者等待")
		q.notFull.Wait()
	}

	q.normalItems = append(q.normalItems, item)
	fmt.Printf("生产普通: %d (普:%d 急:%d)\n", item, len(q.normalItems), len(q.priorityItems))

	// 同时通知两种消费者
	q.notEmpty.Signal()    // 唤醒任意一个消费者
	q.hasPriority.Signal() // 确保紧急消费者能被唤醒
}

// 生产紧急菜品
func (q *Queue) ProducePriority(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.normalItems)+len(q.priorityItems) == q.capacity {
		fmt.Println("队列满，紧急生产者等待")
		q.notFull.Wait()
	}

	q.priorityItems = append(q.priorityItems, item)
	fmt.Printf("‼️生产紧急: %d (普:%d 急:%d)\n", item, len(q.normalItems), len(q.priorityItems))

	// 优先唤醒紧急消费者
	q.hasPriority.Broadcast() // 广播唤醒所有紧急消费者
	q.notEmpty.Signal()       // 同时唤醒普通消费者
}

// 消费菜品（优先紧急）
func (q *Queue) Consume() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 优先检查紧急菜品
	for len(q.priorityItems) == 0 && len(q.normalItems) == 0 {
		fmt.Println("队列全空，消费者等待")
		q.notEmpty.Wait() // 等待任何菜品
	}

	var item int
	if len(q.priorityItems) > 0 {
		item = q.priorityItems[0]
		q.priorityItems = q.priorityItems[1:]
		fmt.Printf("取出紧急: %d (普:%d 急:%d)\n", item, len(q.normalItems), len(q.priorityItems))
	} else {
		item = q.normalItems[0]
		q.normalItems = q.normalItems[1:]
		fmt.Printf("取出普通: %d (普:%d)\n", item, len(q.normalItems))
	}

	q.notFull.Signal() // 通知生产者有空位
	return item
}

// 专用紧急消费者
func (q *Queue) ConsumePriority() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.priorityItems) == 0 {
		fmt.Println("无紧急菜品，VIP等待")
		q.hasPriority.Wait() // 只等待紧急菜品
	}

	item := q.priorityItems[0]
	q.priorityItems = q.priorityItems[1:]
	fmt.Printf("VIP取紧急: %d (急:%d)\n", item, len(q.priorityItems))

	q.notFull.Signal()
	return item
}

func main() {
	queue := NewQueue(2) // 小容量更容易触发边界条件

	var wg sync.WaitGroup

	// 1个普通消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			item := queue.Consume()
			fmt.Printf("普通消费者处理: %d\n", item)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// 1个VIP消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			item := queue.ConsumePriority()
			fmt.Printf("VIP消费者处理: %d\n", item)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// 2个生产者（1普通1紧急）
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			queue.ProduceNormal(100 + i)
			time.Sleep(400 * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond) // 让普通生产者先运行
		for i := 0; i < 2; i++ {
			queue.ProducePriority(999 + i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	wg.Wait()
	fmt.Println("=== 运行结束 ===")
}
