package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "hello"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "world"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("Received from c1:", msg1)
		case msg2 := <-c2:
			fmt.Println("Received from c2:", msg2)
		}
	}
}
