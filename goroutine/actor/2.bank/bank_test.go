package bank

import (
	"sync"
	"testing"
)

func TestAccount_Deposit(t *testing.T) {
	acc := NewAccount(100)

	// 测试正常存款
	acc.Deposit(50)
	if balance := acc.Balance(); balance != 150 {
		t.Errorf("Expected balance 150 after deposit, got %d", balance)
	}

	// 测试负数存款
	err := acc.Deposit(-10)
	if err == nil {
		t.Error("Expected error for negative deposit, got nil")
	}
	if balance := acc.Balance(); balance != 150 {
		t.Errorf("Balance should not change after failed deposit, got %d", balance)
	}
}

func TestAccount_Withdraw(t *testing.T) {
	acc := NewAccount(100)

	// 测试正常取款
	err := acc.Withdraw(30)
	if err != nil {
		t.Errorf("Unexpected error on valid withdrawal: %v", err)
	}
	if balance := acc.Balance(); balance != 70 {
		t.Errorf("Expected balance 70 after withdrawal, got %d", balance)
	}

	// 测试透支
	err = acc.Withdraw(100)
	if err == nil {
		t.Error("Expected error for overdraft, got nil")
	}
	if balance := acc.Balance(); balance != 70 {
		t.Errorf("Balance should not change after failed withdrawal, got %d", balance)
	}

	// 测试负数取款
	err = acc.Withdraw(-10)
	if err == nil {
		t.Error("Expected error for negative withdrawal, got nil")
	}
}

func TestAccount_ConcurrentAccess(t *testing.T) {
	acc := NewAccount(0)
	var wg sync.WaitGroup

	// 并发存款
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acc.Deposit(1)
		}()
	}

	// 并发取款
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acc.Withdraw(1)
		}()
	}

	wg.Wait()

	if balance := acc.Balance(); balance != 50 {
		t.Errorf("Expected final balance 50, got %d", balance)
	}
}

//func TestTransfer(t *testing.T) {
//	acc1 := NewAccount(100)
//	acc2 := NewAccount(50)
//
//	// 测试正常转账
//	err := Transfer(acc1, acc2, 30)
//	if err != nil {
//		t.Errorf("Unexpected error on valid transfer: %v", err)
//	}
//	if acc1.Balance() != 70 || acc2.Balance() != 80 {
//		t.Errorf("After transfer, expected acc1=70 acc2=80, got acc1=%d acc2=%d",
//			acc1.Balance(), acc2.Balance())
//	}
//
//	// 测试透支转账
//	err = Transfer(acc1, acc2, 100)
//	if err == nil {
//		t.Error("Expected error for overdraft transfer, got nil")
//	}
//	// 检查余额不应改变
//	if acc1.Balance() != 70 || acc2.Balance() != 80 {
//		t.Error("Balances should not change after failed transfer")
//	}
//
//	// 测试负数转账
//	err = Transfer(acc1, acc2, -10)
//	if err == nil {
//		t.Error("Expected error for negative transfer, got nil")
//	}
//}
//
//func TestAccount_Close(t *testing.T) {
//	acc := NewAccount(100)
//
//	// 测试关闭后操作
//	acc.Close()
//
//	err := acc.Deposit(10)
//	if err == nil {
//		t.Error("Expected error when depositing to closed account, got nil")
//	}
//
//	err = acc.Withdraw(10)
//	if err == nil {
//		t.Error("Expected error when withdrawing from closed account, got nil")
//	}
//
//	// 余额查询在关闭后应该仍然可以工作
//	if balance := acc.Balance(); balance != 100 {
//		t.Errorf("Expected balance 100 after closing, got %d", balance)
//	}
//}
//
//func TestAccount_ConcurrentTransfer(t *testing.T) {
//	acc1 := NewAccount(1000)
//	acc2 := NewAccount(1000)
//	var wg sync.WaitGroup
//
//	// 并发转账测试
//	for i := 0; i < 100; i++ {
//		wg.Add(2)
//		go func() {
//			defer wg.Done()
//			Transfer(acc1, acc2, 10)
//		}()
//		go func() {
//			defer wg.Done()
//			Transfer(acc2, acc1, 5)
//		}()
//	}
//
//	wg.Wait()
//
//	// 验证总金额不变
//	total := acc1.Balance() + acc2.Balance()
//	if total != 2000 {
//		t.Errorf("Total money should remain 2000, got %d", total)
//	}
//}
