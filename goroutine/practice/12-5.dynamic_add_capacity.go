package main

import (
	"fmt"
	"sync"
	"time"
)

//任务2：实现动态扩容
//当普通区满时，临时扩大容量（思考如何保证线程安全）

type Item struct {
	value    int
	priority bool
}

type Queue struct {
	items           []Item
	baseCapacity    int
	currentCapacity int
	mu              sync.Mutex
	notFull         *sync.Cond
	notEmpty        *sync.Cond
	expanded        bool
}

func NewQueue(baseCapacity int) *Queue {
	q := &Queue{
		baseCapacity:    baseCapacity,
		currentCapacity: baseCapacity,
	}
	q.notFull = sync.NewCond(&q.mu)
	q.notEmpty = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) expand() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.expanded {
		return
	}

	oldCap := q.currentCapacity
	q.currentCapacity = q.baseCapacity * 2
	q.expanded = true
	fmt.Printf("【扩容】%d → %d\n", oldCap, q.currentCapacity)

	// 唤醒所有等待的生产者
	q.notFull.Broadcast()
}

func (q *Queue) shrink() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.expanded || len(q.items) > q.baseCapacity {
		return
	}

	oldCap := q.currentCapacity
	q.currentCapacity = q.baseCapacity
	q.expanded = false
	fmt.Printf("【缩容】%d → %d\n", oldCap, q.currentCapacity)
}

func (q *Queue) ProduceNormal(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 尝试扩容
	if len(q.items) >= q.currentCapacity {
		q.mu.Unlock()
		q.expand()
		q.mu.Lock()
	}

	// 再次检查
	for len(q.items) >= q.currentCapacity {
		fmt.Println("队列满，生产者等待")
		q.notFull.Wait()
	}

	q.items = append(q.items, Item{item, false})
	fmt.Printf("生产普通: %d (总数: %d/%d)\n", item, len(q.items), q.currentCapacity)

	// 确保至少一个消费者被唤醒
	q.notEmpty.Signal()
}

func (q *Queue) ProducePriority(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 紧急情况直接扩容
	if len(q.items) >= q.currentCapacity {
		q.mu.Unlock()
		q.expand()
		q.mu.Lock()
	}

	for len(q.items) >= q.currentCapacity {
		fmt.Println("队列满，紧急生产者等待")
		q.notFull.Wait()
	}

	q.items = append([]Item{{item, true}}, q.items...)
	fmt.Printf("‼️生产紧急: %d (总数: %d/%d)\n", item, len(q.items), q.currentCapacity)

	// 紧急情况广播通知
	q.notEmpty.Broadcast()
}

func (q *Queue) Consume() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 {
		fmt.Println("队列空，消费者等待")
		q.notEmpty.Wait()
	}

	var item Item
	if q.items[0].priority {
		item = q.items[0]
		q.items = q.items[1:]
		fmt.Printf("取出紧急: %d (剩余: %d/%d)\n", item.value, len(q.items), q.currentCapacity)
	} else {
		// 检查是否有紧急项
		for i, it := range q.items {
			if it.priority {
				item = it
				q.items = append(q.items[:i], q.items[i+1:]...)
				fmt.Printf("取出紧急: %d (剩余: %d/%d)\n", item.value, len(q.items), q.currentCapacity)
				goto notify
			}
		}

		item = q.items[0]
		q.items = q.items[1:]
		fmt.Printf("取出普通: %d (剩余: %d/%d)\n", item.value, len(q.items), q.currentCapacity)
	}

notify:
	q.notFull.Signal()

	// 异步缩容检查
	if len(q.items) < q.baseCapacity/2 {
		go q.shrink()
	}

	return item.value
}

func main() {
	queue := NewQueue(2)

	var wg sync.WaitGroup

	// 生产者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			queue.ProduceNormal(100 + i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 紧急生产者
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(150 * time.Millisecond)
		queue.ProducePriority(999)
		queue.ProducePriority(888)
	}()

	// 消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 7; i++ {
			item := queue.Consume()
			fmt.Printf("处理: %d\n", item)
			time.Sleep(250 * time.Millisecond)
		}
	}()

	wg.Wait()
	fmt.Println("=== 运行结束 ===")
}
