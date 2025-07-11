package main

import (
	"fmt"
)

type Counter struct {
	value   int
	addChan chan int
	getChan chan chan int
}

func NewCounter() *Counter {
	c := &Counter{
		addChan: make(chan int),
		getChan: make(chan chan int),
	}

	go c.run()

	return c
}

func (c *Counter) run() {
	for {
		select {
		case n := <-c.addChan:
			c.value += n
		case ch := <-c.getChan:
			ch <- c.value
		}
	}
}

func (c *Counter) Add(n int) {
	c.addChan <- n
}

func (c *Counter) Get() int {
	resultChan := make(chan int)
	c.getChan <- resultChan
	return <-resultChan
}

func main() {
	counter := NewCounter()

	counter.Add(10)
	fmt.Println(counter.Get()) // 输出: 10

	counter.Add(5)
	fmt.Println(counter.Get()) // 输出: 15
}
