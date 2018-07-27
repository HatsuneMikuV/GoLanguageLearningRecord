package main

import (
	"io"
	"log"
	"net"
	"time"
)


/*
	练习 8.1： 修改clock2来支持传入参数作为端口号，
	然后写一个clockwall的程序，
	这个程序可以同时与多个clock服务器通信，
	从多服务器中读取时间，
	并且在一个表格中一次显示所有服务传回的结果，
	类似于你在某些办公室里看到的时钟墙。
*/

//const area  = ""
//const area  = "Eastern"
//const area  = "Tokyo"
const area  = "London"

func main() {

	//1.
	//listener, err := net.Listen("tcp", "localhost:8000")
	//2.Eastern
	//listener, err := net.Listen("tcp", "localhost:8010")
	//3.Tokyo
	//listener, err := net.Listen("tcp", "localhost:8020")
	//4.London
	listener, err := net.Listen("tcp", "localhost:8030")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConnWall(conn)
	}
}


func handleConnWall(c net.Conn)  {
	defer c.Close()

	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05--") + area + "\n" )
		if err != nil {
			return
		}
		time.Sleep(time.Second * 1.0)
	}
}