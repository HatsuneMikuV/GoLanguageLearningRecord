package main

import (
	"20180408/byteCounter"
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"strings"
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
//1.接口类型具体描述了一系列方法的集合, 而实现这个方法的具体类型是这个接口类型的实例
//2.io.Writer类型是用的最广泛的接口之一，因为它提供了所有的类型写入bytes的抽象
//3.接口类型可以组合定义，成为一个集合方法
func test_interface_type()  {

	//练习7.4
	node,_ := NewReader("<html>111111</html>")
	fmt.Println(node)
	fmt.Println("================================")

	//练习 7.5
	ss := []byte("11112222")
	rr := bytes.NewReader(ss)
	reader := LimitReader(rr, 6)
	s, _ := ioutil.ReadAll(reader)
	fmt.Println(string(s))

}
type HtmlReader struct {
	r io.Reader
}
func (reader *HtmlReader) Read(p []byte) (n int, err error) {
	n, err = reader.r.Read(p)
	return
}
func creatReader(r io.Reader) io.Reader {
	return &HtmlReader{r:r}
}
func NewReader(s string) (*html.Node, error)  {
	hr := creatReader(strings.NewReader(s))
	node, err := html.Parse(hr)
	return node, err
}
/*
   练习 7.5： io包里面的LimitReader函数接收一个io.Reader接口类型的r和字节数n，
   并且返回另一个从r中读取字节,但是当读完n个字节后就表示读到文件结束的Reader。
   实现这个LimitReader函数.
*/
type MyLimitReader struct {
	R io.Reader
	N int64
}

func (myLimitReader *MyLimitReader) Read(p []byte) (n int, err error) {
	if myLimitReader.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p))  > myLimitReader.N {
		p = p[0:myLimitReader.N]
	}
	n, err = myLimitReader.R.Read(p)
	myLimitReader.N -= int64(n)
	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &MyLimitReader{r, n}
}


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
	//test_convention()

	//二，接口类型
	test_interface_type()
}
