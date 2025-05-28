package main

import (
	"errors"
	"fmt"
)

type argError struct {
	arg     int
	message error
}

func (a *argError) Error() string {
	return fmt.Sprintf("arg %d: %s", a.arg, a.message)
}

func f(arg int) (int, error) {
	if arg == 42 {
		return -1, &argError{arg, fmt.Errorf("can't work with %d", arg)}
	}
	return arg + 3, nil
}

// errors.Is(err, target)
// → 判断 err 是否等于某个特定错误（如 io.EOF），译为 "错误匹配"。
//
// errors.As(err, &target)
// → 判断 err 是否属于某种类型，译为 "错误类型解析" 或 "错误类型匹配"。
func main() {
	_, err := f(42)
	var ae *argError
	if errors.As(err, &ae) {
		//如果 err 是 *argError 类型，ae 会被赋值并进入此逻辑块
		fmt.Println(ae.arg)
		fmt.Println(ae.message)
	} else {
		fmt.Println("err doesn't match argError")
	}
}
