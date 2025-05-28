package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int, 10)

	for i := 0; i < 10; i++ {
		ch <- i
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			egg := <-ch
			fmt.Printf("People: %d, Egg: %d\n", num, egg)
			select {
			case egg := <-ch:
				fmt.Printf("People: %d, Egg: %d\n", num, egg)
			default:
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
