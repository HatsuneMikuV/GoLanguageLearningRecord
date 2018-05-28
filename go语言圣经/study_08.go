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
//11.一个Channel的输出作为下一个Channel的输入, 这种串联的Channels就是所谓的管道（pipeline）
func test_Channels()  {

	//演示在netcat的channel_84
	//test_echo()

	//串行channel
	//test_pipeline()

	//使用单向channel实现串行
	//单向操作的channel中的只接受方是不需要关闭的，发送者关闭即可
	//任何双向channel向单向channel变量的赋值操作都将导致该隐式转换
	//naturals := make(chan int)
	//squares := make(chan int)
	//go counter_chan(naturals)
	//go squarer(squares, naturals)
	//printer(squares)

	//带缓存channel
	test_channel_cache()
}

//12.串行的channel是无法知道另一个channel是否关闭的，
// 但是关闭的channel会产生第二个值，为bool类型，ture代表接收到值，false表示channel已被关闭没有可接收的值
func test_pipeline()  {
	naturals := make(chan int)
	squares := make(chan int)

	max := 100

	// Counter
	go func() {
		for x := 0; ; x++  {
			if x > max {
				close(naturals)
				break
			}
			naturals <- x
		}
	}()

	// Squarer
	go func() {
		for {
			x, ok := <- naturals
			if !ok {
				fmt.Println("naturals close")
				break
			}
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	for {
		x, ok := <-squares
		if !ok {
			fmt.Println("squares close")
			break
		}
		fmt.Println(x)
	}
}
func counter_chan(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}
//13.带缓存的channel内部有一个元素队列，容量在创建时通过第二个参数指定
//14.带缓存channel会因为队列满导致发送等待，队列空接收等待
//15.发送的内容添加到队列尾部，接收内容是从队列头部取出并删除
func test_channel_cache()  {
	chan_cache := make(chan string, 3)
	//channel缓存容量
	fmt.Println(cap(chan_cache))

	chan_cache <- "A"
	chan_cache <- "B"
	chan_cache <- "C"
	//有效元素个数
	fmt.Println(len(chan_cache))

	fmt.Println(<-chan_cache) // "A"
	//有效元素个数
	fmt.Println(len(chan_cache))

	chan_cache <- "AA"
	//有效元素个数
	fmt.Println(len(chan_cache))

	fmt.Println(<-chan_cache) // "B"
	fmt.Println(<-chan_cache) // "C"
	//有效元素个数
	fmt.Println(len(chan_cache))

	result := mirroredQuery()
	fmt.Println(result)
}
func mirroredQuery() string {
	responses := make(chan string, 3)
	go func() { responses <- request("asia.gopl.io") }()
	go func() { responses <- request("europe.gopl.io") }()
	go func() { responses <- request("americas.gopl.io") }()
	return <-responses // return the quickest response
}
func request(hostname string) (response string) {
	if hostname == "asia.gopl.io" {
		time.Sleep(1 * time.Second)
		return "111"
	}else if hostname == "europe.gopl.io" {
		time.Sleep(2 * time.Second)
		return "222"
	}else if hostname == "americas.gopl.io" {
		time.Sleep(500 * time.Millisecond)
		return "333"
	}
	return "000"
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
