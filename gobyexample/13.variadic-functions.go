package main

import "fmt"

func sum(nums ...int) {
	fmt.Println("nums", nums)
	total := 0

	for _, num := range nums {
		total += num
	}
	fmt.Println("total", total)
}

func main() {
	sum(1, 3)
	sum(1, 2, 3)

	nums := []int{1, 2, 3, 4}
	sum(nums...)
}
