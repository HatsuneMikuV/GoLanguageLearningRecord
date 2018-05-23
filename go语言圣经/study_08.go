package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

/*
	Goroutines和Channels
	1.并发程序指同时进行多个任务的程序
	2.goroutine
	3.channel
*/

//一，Goroutines
//1.在Go语言中，每一个并发的执行单元叫作一个goroutine
//2.新的goroutine会用go语句来创建,go语句会使其语句中的函数在一个新创建的goroutine中运行
//3.主函数返回时，所有的goroutine都会被直接打断，程序退出
func test_goroutine()  {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

//二，示例: 并发的Clock服务
//1.网络编程是并发大显身手的一个领域
//2.go语言的net包，提供编写一个网络客户端或者服务器程序的基本组件
func test_clock()  {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(c net.Conn)  {
	defer c.Close()

	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(time.Second * 1.0)
	}
}

//三，示例: 并发的Echo服务
//1.

func test_echo()  {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go echo_handleConn(conn)
	}
}
func echo_handleConn(c net.Conn)  {
	go echo(c, "Hello", 1*time.Second)
	c.Close()
}
func echo(c net.Conn, shut string, delay time.Duration)  {
	fmt.Fprintf(c, "\t", strings.ToUpper(shut))
	time.Sleep(delay)
	fmt.Fprintf(c, "\t", shut)
	time.Sleep(delay)
	fmt.Fprintf(c, "\t", strings.ToLower(shut))
}
func main() {
	//一，Goroutines
	//test_goroutine()

	//二，示例: 并发的Clock服务
	//test_clock()

	//三，示例: 并发的Echo服务
	test_echo()
	
}
