package bank

import "sync"

// Package bank provides a concurrency-safe bank with one account.

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

var (
	mu     sync.RWMutex
	balance int
)
func deposit(amount int) { balance += amount }

func Deposit(amount int) {
	mu.Lock()
	defer 	mu.Unlock()
	deposit(amount)
}
func Balance() int       {
	mu.Lock()
	defer mu.Unlock()
	return balance
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

//取款用函数
func Withdraw(amount int)bool{
	mu.Lock()
	defer mu.Unlock()
	Deposit(-amount)
	if Balance() < 0 {
		Deposit(amount)
		return false // insufficient funds
	}
	return true
}

func Init() {
	go teller() // start the monitor goroutine
}