package main

import (
	"20180408/byteCounter"
	"fmt"
)

/*
	第七章 接口

	1.接口类型是对其它类型行为的抽象和概括
	2.Go语言中接口类型的独特之处在于它是满足隐式实现的
*/


//一，接口是合约
//1.接口合约就像ObjC的继承，继承父类的方法并重新实现，因此可以被合约原函数调用
func test_convention()  {

	fmt.Println("test")

	//因为*ByteCounter满足io.Writer的约定，我们可以把它传入Fprintf函数中
	//
	var c byteCounter.ByteCounter
	c.Write([]byte("hello"))
	fmt.Println(c) // "5", = len("hello")
	c = 0          // reset the counter
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c) // "12", = len("hello, Dolly")

	fmt.Println("================================")

	//练习 7.1： 使用来自ByteCounter的思路，实现一个针对对单词和行数的计数器。
	ss := []byte("hello\nI like you\nI like you too")

	var w byteCounter.WordCounter
	w.Write(ss)
	fmt.Println(w)
	w = 0
	fmt.Fprintf(&w, "hi, %s", string(ss))
	fmt.Println(w)
	fmt.Println("================================")

	var l byteCounter.LineCounter
	l.Write(ss)
	fmt.Println(l)
	l = 0
	fmt.Fprintf(&l, "hi\n %s", string(ss))
	fmt.Println(l)
	fmt.Println("================================")

	//练习 7.2： 写一个带有如下函数签名的函数CountingWriter，
	//传入一个io.Writer接口类型，
	//返回一个新的Writer类型把原来的Writer封装在里面和一个表示写入新的Writer字节数的int64类型指针
	var b byteCounter.ByteCounter
	b.Write([]byte("hello"))
	cc, nbytes := byteCounter.CountingWriter(&b)
	fmt.Println(cc, nbytes)
	fmt.Println("================================")

	//练习 7.3： 为在gopl.io/ch4/treesort (§4.4)的*tree类型实现一个String方法去展示tree类型的值序列。
	arr := []int{9, 8, 3, 4, 5, 6, 7}
	tree1 := byteCounter.Sort(arr)
	fmt.Println(tree1)
}


//二，接口类型
//三，实现接口的条件
//四，flag.Value接口
//五，接口值
//六，sort.Interface接口
//七，http.Handler接口
//八，error接口
//九，示例: 表达式求值
//十，类型断言
//十一，基于类型断言识别错误类型
//十二，通过类型断言查询接口
//十三，类型分支
//十四，示例: 基于标记的XML解码
//十五，补充几点


func main() {

	//一，接口是合约
	test_convention()
}
