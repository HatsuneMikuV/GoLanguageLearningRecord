package main

import (
	"bufio"
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
//1.让服务使用并发不只是处理多个客户端的请求，甚至在处理单个连接时也可能会用到

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
	input := bufio.NewScanner(c)
	for input.Scan() {
		str := input.Text()
		if str == "EOF" {
			break
		}
		if len(str) > 0 {
			go echo(c, str, 1*time.Second)
		}
	}
	c.Close()
}
func echo(c net.Conn, shut string, delay time.Duration)  {
	fmt.Fprintf(c, "\t", strings.ToUpper(shut))
	time.Sleep(delay)
	fmt.Fprintf(c, "\t", shut)
	time.Sleep(delay)
	fmt.Fprintf(c, "\t", strings.ToLower(shut))
}

//四，Channels
//1.goroutine是Go语言程序的并发体的话，channels则是它们之间的通信机制
//2.一个channel是一个通信机制,它可以让一个goroutine通过它给另一个goroutine发送值信息
//3.channel的零值也是nil
//4.两个相同类型的channel可以使用==运算符比较，同样可以和nil比较
//5.channel有发送和接受两个主要操作，都是通信行为
//6.<-  channel在左侧是发送，在右侧是接收，不使用接收数据操作符也是合法的
//7.如果channel被关闭，则任何发送操作都会导致panic异常，
// 	但是依然能接收到成功发送的数据，如果没有数据，则产生零值数据
//8.channel的make创建第二个参数，对应其容量

//9.无缓存Channels的发送和接收操作将导致两个goroutine做一次同步操作，因此被称为同步Channels
//10.happens before
func test_Channels()  {

	//演示在netcat的channel_84
	test_echo()
}

func main() {
	//一，Goroutines
	//test_goroutine()

	//二，示例: 并发的Clock服务
	//test_clock()

	//三，示例: 并发的Echo服务
	//test_echo()

	//四，Channels
	test_Channels()
}
