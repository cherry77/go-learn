package main

import (
	"fmt"
	"sync"
)

//使用 sync.RWMutex
//sync.RWMutex：通过读写锁实现线程安全，适合需要更多控制的场景。
// 需要更高的写性能和更复杂的操作控制，sync.RWMutex 更适合。

type SafeMap struct {
	mu sync.RWMutex           // 读写锁
	m  map[string]interface{} // 键值对存储
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]interface{}),
	}
}

func (sm *SafeMap) Store(key string, value interface{}) {
	sm.mu.Lock()         // 获取写锁
	defer sm.mu.Unlock() // 确保函数结束时释放锁
	sm.m[key] = value    // 存储键值对
}

func (sm *SafeMap) Load(key string) (interface{}, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, ok := sm.m[key]
	return value, ok
}
func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()         // 获取写锁
	defer sm.mu.Unlock() // 确保函数结束时释放锁
	delete(sm.m, key)    // 删除键值对
}

func (sm *SafeMap) Range(f func(key string, value interface{}) bool) {
	sm.mu.RLock()         // 获取读锁
	defer sm.mu.RUnlock() // 确保函数结束时释放锁

	for k, v := range sm.m {
		if !f(k, v) {
			break
		}
	}
}

func main() {
	//sm := NewSafeMap()
	//
	//sm.Store("key1", "value1")
	//
	//if value, ok := sm.Load("key1"); ok {
	//	fmt.Println("Found:", value)
	//}
	//
	//sm.Delete("key1")
	//
	//sm.Range(func(key string, value interface{}) bool {
	//	fmt.Println(key, value)
	//	return true // 返回 false 可以停止迭代
	//})

	m := NewSafeMap()

	// 并发写入
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Store(fmt.Sprintf("key%d", i), i)
		}(i)
	}

	// 并发读取
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if val, ok := m.Load(fmt.Sprintf("key%d", i)); ok {
				fmt.Printf("key%d: %v\n", i, val)
			}
		}(i)
	}

	wg.Wait()
}
