package main

import (
	"20180408/bank"
	"fmt"
	"sync"
)

/*
	第九章　基于共享变量的并发
*/


//一， 竞争条件
//1.在一个线性(就是说只有一个goroutine的)的程序中，程序的执行顺序只由程序的逻辑来决定
//2.一个特定类型的一些方法和操作函数，对于某个类型来说，如果其所有可访问的方法和操作都是并发安全的话，那么类型便是并发安全的
//3.在一个程序中有非并发安全的类型的情况下，我们依然可以使这个程序并发安全
//4.导出包级别的函数一般情况下都是并发安全的
//5.竞争条件指的是程序在多个goroutine交叉执行操作时，没有给出正确的结果
//6.避免数据竞争：第一种方法是不要去写变量
//7.避免数据竞争：避免从多个goroutine访问变量
//8.避免数据竞争：允许很多goroutine去访问变量，但是在同一个时刻最多只有一个goroutine在访问

func test_competition()  {

	//-----------

	bank.Init()

	var wg sync.WaitGroup
	wg.Add(1)
	// Alice:
	go func() {
		defer wg.Done()
		bank.Deposit(200)                // A1
	}()
	wg.Add(1)
	// Bob:
	go func(){
		defer wg.Done()
		bank.Deposit(100)
	}()

	wg.Add(1)
	go func(){
		defer wg.Done()
		res := bank.Withdraw(200)
		if res{
			fmt.Println("取款成功")
		}else {
			fmt.Println("取款失败")
		}
	}()
	wg.Wait()

	res := bank.Balance()
	fmt.Printf("剩余:%d" , res)



}

func main() {

	//一， 竞争条件
	test_competition()
}
