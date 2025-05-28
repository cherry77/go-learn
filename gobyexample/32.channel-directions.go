package main

import "fmt"

// 向只写通道（chan<-）发送消息
func ping(pings chan<- string, msg string) {
	pings <- msg // 只能发送
}

// 从只读通道（<-chan）接收消息，并转发到另一个只写通道。
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings // 只能接收
	pongs <- msg   // 只能发送
}

//通道方向的作用
//提高类型安全：明确函数对通道的操作权限（只读或只写），避免误用。
//
//例如，ping 函数无法从 pings 读取，pong 函数无法向 pings 写入。
//
//代码可读性：清晰表达设计意图（如“生产者-消费者”模式）。

//类比现实场景
//ping：像快递员，只能将包裹（消息）投递到指定信箱（pings）。
//
//pong：像邮局中转站，从信箱取包裹，再投递到另一个信箱（pongs）。
//
//主函数：像收件人，最终从 pongs 取包裹。

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message") // 向 pings 发送消息
	pong(pings, pongs)            // 从 pings 接收，并转发到 pongs
	fmt.Println(<-pongs)          // 从 pongs 接收并打印
}
