package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)
/*
	配合study_8使用的客户端
*/
func main() {

	//8.2
	//netcat01()

	//练习 8.3： 在netcat3例子中，conn虽然是一个interface类型的值，
	// 但是其底层真实类型是*net.TCPConn，代表一个TCP连接。
	// 一个TCP连接有读和写两个部分，可以使用CloseRead和CloseWrite方法分别关闭它们。
	// 修改netcat3的主goroutine代码，只关闭网络连接中写的部分，
	// 这样的话后台goroutine可以在标准输入被关闭后继续打印从reverb1服务器传回的数据
	netcat02()

	//练习 8.1
	//homework()

	//8.4
	//channel_84()

}
func netcat01()  {
	//1.
	conn, err := net.Dial("tcp", "localhost:8000")
	//2.Eastern
	//conn, err := net.Dial("tcp", "localhost:8010")
	//3.Tokyo
	//conn, err := net.Dial("tcp", "localhost:8020")
	//4.London
	//conn, err := net.Dial("tcp", "localhost:8030")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	mustCopy(os.Stdout, conn)
}
func netcat02()  {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	cw := conn.(*net.TCPConn)
	cw.CloseWrite()
	<-done // wait for background goroutine to finish
	conn.Close()
}
func homework()  {

	hosts := []string{"localhost:8000", "localhost:8010", "localhost:8020", "localhost:8030"}

	for {
		for _, host := range hosts {
			conn, err := net.Dial("tcp", host)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			go mustCopy(os.Stdout, conn)
		}
		time.Sleep(10 * time.Second)
	}

}

func channel_84()  {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}
func mustCopy(dst io.Writer, src io.Reader)  {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}
