package main

import "fmt"

//15. 管道模式
//实现一个多阶段的管道处理：第一 stage 生成数字，第二 stage 计算平方，第三 stage 过滤偶数，最后收集结果。

// 第一阶段：生成数字
func generateNumbers(done <-chan struct{}, numbers ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range numbers {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// 第二阶段：计算平方
func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// 第三阶段：过滤偶数
func filterEven(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 != 0 {
				select {
				case out <- n:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

// 收集结果
func collect(done <-chan struct{}, in <-chan int) []int {
	var results []int
	for n := range in {
		select {
		case <-done:
			return nil
		default:
			results = append(results, n)
		}
	}
	return results
}

func main() {
	done := make(chan struct{})
	defer close(done)

	numbers := generateNumbers(done, 1, 2, 3, 4, 5)
	squares := square(done, numbers)
	odds := filterEven(done, squares)
	results := collect(done, odds)

	fmt.Println("结果:", results)
}
