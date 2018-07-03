package main

import (
	"20180408/format"
	"fmt"
	"time"
)

/*
	第十二章　反射

	Go语言提供了一种反射机制，
	能够在运行时更新变量和检查它们的值、调用它们的方法和它们支持的内在操作，
	而不需要在编译时就知道这些变量的具体类型。

	反射是一个复杂的内省技术，
	不应该随意使用

	fmt包提供的字符串格式功能，encoding/json和encoding/xml，text/template和html/template包
	这些包内部都是用反射技术实现的，但是它们自己的API都没有公开反射相关的接口
*/


//一，为何需要反射?
//1.为了处理不满足公共接口的类型的值，或者设计函数还不存在的类型


//二，reflect.Type和reflect.Value
//1.reflect包提供了反射功能，定义两个类型Type和Value
//2.Type表示一个Go类型. 它是一个接口
//3.一个Value，有很多方法来检查其内容, 无论具体类型是什么

func test_format_reflect()  {

	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(format.Any(x))                  // "1"
	fmt.Println(format.Any(d))                  // "1"
	fmt.Println(format.Any([]int64{x}))         // "[]int64 0x8202b87b0"
	fmt.Println(format.Any([]time.Duration{d})) // "[]time.Duration 0x8202b87e0"
}

func main() {
	//二，reflect.Type和reflect.Value
	test_format_reflect()
}
