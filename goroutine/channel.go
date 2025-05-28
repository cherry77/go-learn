package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		fmt.Printf("Hello, %d\n", <-ch)
	}()

	ch <- 42

	close(ch)

	fmt.Println(<-ch)

}
