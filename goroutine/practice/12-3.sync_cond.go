package main

import (
	"fmt"
	"sync"
	"time"
)

//在 12-2 的基础上增加
// 任务1：增加紧急菜品通道

// Queue 代表餐厅的餐盘架
type Item struct {
	value    int
	priority bool // true表示紧急菜品
}

type Queue struct {
	items    []Item // 统一存储，通过字段区分优先级
	capacity int
	mu       sync.Mutex
	notFull  *sync.Cond
	notEmpty *sync.Cond
}

func NewQueue(capacity int) *Queue {
	q := &Queue{capacity: capacity}
	q.notFull = sync.NewCond(&q.mu)
	q.notEmpty = sync.NewCond(&q.mu)
	return q
}

// 生产普通菜品
func (q *Queue) ProduceNormal(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == q.capacity {
		q.notFull.Wait()
	}

	q.items = append(q.items, Item{item, false})
	fmt.Printf("生产普通菜品: %d (总数: %d)\n", item, len(q.items))
	q.notEmpty.Signal()
}

// 生产紧急菜品
func (q *Queue) ProducePriority(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == q.capacity {
		q.notFull.Wait()
	}

	// 紧急菜品插入队列头部
	q.items = append([]Item{{item, true}}, q.items...)
	fmt.Printf("‼️ 生产紧急菜品: %d (总数: %d)\n", item, len(q.items))
	q.notEmpty.Broadcast() // 紧急情况广播通知
}

// 消费菜品（优先取紧急）
func (q *Queue) Consume() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 {
		q.notEmpty.Wait()
	}

	// 查找第一个紧急菜品
	for i, item := range q.items {
		if item.priority {
			ret := item.value
			q.items = append(q.items[:i], q.items[i+1:]...)
			fmt.Printf("取出紧急菜品: %d (剩余: %d)\n", ret, len(q.items))
			q.notFull.Signal()
			return ret
		}
	}

	// 没有紧急则取第一个普通
	ret := q.items[0].value
	q.items = q.items[1:]
	fmt.Printf("取出普通菜品: %d (剩余: %d)\n", ret, len(q.items))
	q.notFull.Signal()
	return ret
}

func main() {
	queue := NewQueue(3)

	go func() {
		for i := 0; i < 5; i++ {
			queue.ProduceNormal(100 + i)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		time.Sleep(300 * time.Millisecond) // 稍后插入紧急
		queue.ProducePriority(999)
		queue.ProducePriority(888)
	}()

	for i := 0; i < 7; i++ {
		item := queue.Consume()
		fmt.Printf("处理菜品: %d\n", item)
		time.Sleep(300 * time.Millisecond)
	}
}
