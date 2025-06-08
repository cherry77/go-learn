package main

import (
	"fmt"
	"sync"
)

//19. 并发安全的对象池
//使用 sync.Pool 实现一个并发安全的对象池，用于复用昂贵创建成本的对象。

/**
类比理解
把 sync.Pool 想象成一个物品租赁商店：

租用（Get）：拿到一个物品，可能是全新的，也可能是别人用过的（上面可能有之前的贴纸或磨损）。

使用：你可以直接使用它当前的状态，或者先清理/改装它。

归还（Put）：还回去时，商店不会自动帮你清理，下次别人可能拿到你归还时的状态。
*/

// ExpensiveObject 代表一个创建成本高的对象
type ExpensiveObject struct {
	ID   int
	Name string
	// 其他可能很耗资源的字段
}

// ObjectPool 是一个并发安全的对象池
type ObjectPool struct {
	pool sync.Pool
	// 可以添加其他池相关的字段，如统计信息等
}

func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		pool: sync.Pool{
			New: func() interface{} {
				// 当池中没有可用对象时，调用New函数创建新对象
				fmt.Println("Creating new expensive object")
				return &ExpensiveObject{}
			},
		},
	}
}

// Get 从池中获取一个对象
func (p *ObjectPool) Get() *ExpensiveObject {
	return p.pool.Get().(*ExpensiveObject)
}

// Put 将对象放回池中
func (p *ObjectPool) Put(obj *ExpensiveObject) {
	// 可以在这里重置对象状态
	obj.ID = 0 // 重置ID或其他字段
	p.pool.Put(obj)
}

func main() {
	//var myPool = sync.Pool{
	//	New: func() interface{} {
	//		// 当池中没有可用对象时，调用此函数创建新对象
	//		return &MyObject{}
	//	},
	//}
	//
	//obj := myPool.Get().(*MyObject) // 类型断言
	//fmt.Println(obj)
	//
	//myPool.Put(obj)
	//fmt.Println(obj)

	// 创建对象池
	pool := NewObjectPool()

	// 获取对象
	obj1 := pool.Get()
	obj1.ID = 1
	fmt.Printf("Object 1: %+v\n", obj1)

	obj2 := pool.Get()
	obj2.ID = 2
	fmt.Printf("Object 2: %+v\n", obj2)

	// 将对象放回池中
	pool.Put(obj1)
	pool.Put(obj2)

	// 再次获取对象，应该会复用之前放回的对象
	obj3 := pool.Get()
	obj3.ID = 3
	fmt.Printf("Object 3 (should be reused): %+v\n", obj3)

	obj4 := pool.Get()
	fmt.Printf("Object 4 (should be reused): %+v\n", obj4)

	// 获取第三个对象，此时池中没有可复用对象，会创建新对象
	obj5 := pool.Get()
	fmt.Printf("Object 5 (should be new): %+v\n", obj5)

}
