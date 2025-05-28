package main

import (
	"fmt"
	"sync"
)

//14. 并发 Map
//实现一个线程安全的 Map，支持并发读写（可以使用 sync.Map 或自己用 mutex 实现）。

//使用 sync.Map
//sync.Map：提供了线程安全的存储，适合简单的并发读写操作。
// 应用场景是简单的、读多写少的，sync.Map 是一个不错的选择

func main() {
	var sm sync.Map
	// 存储键值对
	sm.Store("key1", "value1")

	value, ok := sm.Load("key1")
	if ok {
		fmt.Println("Found", value)
	}

	sm.Delete("key1")

	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true // 返回 false 可以停止迭代
	})
}
