package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var ops atomic.Uint64

	var wg sync.WaitGroup

	for range 50 {
		wg.Add(1)

		go func() {
			for range 1000 {
				ops.Add(1)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("ops", ops.Load())
}
