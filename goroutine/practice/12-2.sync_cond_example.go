package main

import (
	"fmt"
	"sync"
	"time"
)

// Queue 代表餐厅的餐盘架
type Queue struct {
	items    []int      // 餐盘架上的菜品（用数字表示不同菜品）
	capacity int        // 餐盘架容量
	mu       sync.Mutex // 厨房门锁（一次只允许一个人操作餐盘架）
	notFull  *sync.Cond // "餐盘架未满"信号铃（厨师用）
	notEmpty *sync.Cond // "餐盘架不空"信号铃（服务员用）
}

// NewQueue 初始化餐盘架
func NewQueue(capacity int) *Queue {
	q := &Queue{
		items:    make([]int, 0, capacity),
		capacity: capacity,
	}
	// 初始化两个信号铃，都绑定厨房门锁
	q.notFull = sync.NewCond(&q.mu)  // 厨师等待的铃铛
	q.notEmpty = sync.NewCond(&q.mu) // 服务员等待的铃铛
	return q
}

// Produce 厨师往餐盘架放菜品
func (q *Queue) Produce(item int) {
	q.mu.Lock()         // 进入厨房先拿钥匙
	defer q.mu.Unlock() // 离开厨房时归还钥匙

	// 检查餐盘架是否已满
	for len(q.items) == q.capacity {
		fmt.Printf("厨师发现餐盘架已满（当前%d道菜），坐下休息\n", len(q.items))
		q.notFull.Wait() // 1.放下钥匙 2.去休息室睡觉 3.被唤醒后重新拿钥匙
	}

	// 放置菜品
	q.items = append(q.items, item)
	fmt.Printf("厨师放入菜品%d (当前餐盘架: %d/%d)\n", item, len(q.items), q.capacity)

	// 摇铃通知服务员可以取菜了
	q.notEmpty.Signal()
}

// Consume 服务员从餐盘架取菜品
func (q *Queue) Consume() int {
	q.mu.Lock()         // 进入厨房先拿钥匙
	defer q.mu.Unlock() // 离开厨房时归还钥匙

	// 检查餐盘架是否为空
	for len(q.items) == 0 {
		fmt.Println("服务员发现餐盘架空了，玩手机等待")
		q.notEmpty.Wait() // 1.放下钥匙 2.去玩手机 3.被唤醒后重新拿钥匙
	}

	// 取走最旧的菜品（队列先进先出）
	item := q.items[0]
	q.items = q.items[1:]
	fmt.Printf("服务员取走菜品%d (剩余%d道菜)\n", item, len(q.items))

	// 摇铃通知厨师可以继续做菜
	q.notFull.Signal()
	return item
}

func main() {
	// 初始化一个容量为3的餐盘架
	queue := NewQueue(3)

	var wg sync.WaitGroup // 用于等待所有厨师和服务员下班

	// 雇佣2个服务员
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(waiterID int) {
			defer wg.Done()
			for j := 0; j < 5; j++ { // 每个服务员需要完成5次上菜
				// 从餐盘架取菜
				dish := queue.Consume()

				// 模拟上菜耗时（资深服务员动作更快）
				serveTime := 200 + waiterID*100 // 服务员1:300ms, 服务员2:400ms
				time.Sleep(time.Duration(serveTime) * time.Millisecond)
				fmt.Printf("服务员%d上菜完成: %d号菜品\n", waiterID, dish)
			}
		}(i)
	}

	// 雇佣3个厨师
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(chefID int) {
			defer wg.Done()
			for j := 0; j < 4; j++ { // 每个厨师要做4道菜
				// 生成菜品编号（厨师ID*100 + 序号）
				dish := chefID*100 + j

				// 将菜品放到餐盘架
				queue.Produce(dish)

				// 模拟做菜时间（主厨动作更快）
				cookTime := 100 + chefID*50 // 厨师1:150ms, 厨师2:200ms, 厨师3:250ms
				time.Sleep(time.Duration(cookTime) * time.Millisecond)
			}
		}(i)
	}

	// 等待所有厨师和服务员完成工作
	wg.Wait()
	fmt.Println("== 营业结束，所有菜品生产和服务完成 ==")
}
