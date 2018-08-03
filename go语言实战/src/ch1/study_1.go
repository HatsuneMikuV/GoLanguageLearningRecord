package main

import (
	"fmt"
	"time"
)

func log(msg string)  {
	fmt.Println("log: ", msg)
}


func main() {

	//1.Go 语言还自带垃圾回收器，不需要用户自己管理内存
	//2.Go 语言使用接口作为代码复用的基础模块

	fmt.Println("Hello world!")

	go log("日志记录")

	//目的不让主线程直接结束
	time.Sleep(time.Second)
}
