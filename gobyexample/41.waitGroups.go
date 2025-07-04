package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			worker(i)
		}()
	}
	wg.Wait()
}

func worker(id int) {
	fmt.Printf("Worker %d is starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d has finished\n", id)
}
