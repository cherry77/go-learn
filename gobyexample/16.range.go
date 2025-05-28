package main

import "fmt"

func main() {
	// arrays and slices
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	println("sum:", sum)

	for i, num := range nums {
		if num == 3 {
			println("index:", i)
		}
	}

	// maps
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}

	for k := range kvs {
		println("key:", k)
	}

	// strings
	for i, c := range "go" {
		println(i, c)
	}
}
