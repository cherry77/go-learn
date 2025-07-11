package bank

import (
	"fmt"
	"sync"
)

type Account struct {
	balance   int            // 账户余额
	opsChan   chan func()    // 操作通道
	closeChan chan func()    // 关闭信号
	wg        sync.WaitGroup // 用于优雅关闭
}

func NewAccount(balance int) *Account {
	a := &Account{
		balance: balance,
		opsChan: make(chan func()),
	}
	a.wg.Add(1)
	go a.loop()
	return a
}

func (a *Account) loop() {
	defer a.wg.Done()

	for {
		select {
		case op := <-a.opsChan:
			op() // 执行操作
		case <-a.closeChan:
			return // 退出goroutine
		}
	}
}

func (a *Account) Deposit(amount int) error {
	a.opsChan <- func() {
		if amount <= 0 {
			panic(fmt.Sprintf("deposit amount must be positive, got %d", amount))
		}
		a.balance += amount
	}

	return nil
}

func (a *Account) Withdraw(amount int) error {
	a.opsChan <- func() {
		if amount <= 0 || amount > a.balance {
			panic(fmt.Sprintf("withdrawal amount must be positive and less than or equal to balance, got %d", amount))
		}
		a.balance -= amount
	}
	return nil
}

func (a *Account) Balance() int {
	return a.balance
}
