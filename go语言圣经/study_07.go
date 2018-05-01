package main

import (
	"20180408/byteCounter"
	"20180408/tempconv"
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
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
//1.表达一个类型属于某个接口只要这个类型实现这个接口
//2.即使具体类型有其它的方法也只有接口类型暴露出来的方法会被调用到
//3.因为接口类型被称为空接口类型，因此可以将任意值赋给接口类型
func test_interface_condition()  {

	os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
	//os.Stdout.Close()                // OK: *os.File has Close method
	fmt.Println("\n================================")

	var w io.Writer
	w = os.Stdout
	w.Write([]byte("hello")) // OK: io.Writer has Write method
	//w.Close()//w.Close undefined (type io.Writer has no field or method Close)

	fmt.Println("\n================================")

	var any interface{}
	any = true
	any = 12.34
	any = "hello"
	any = map[string]int{"one": 1}
	any = new(bytes.Buffer)
	fmt.Println(any)
	fmt.Println("================================")

	type Artifact interface {
		Title() string
		Creators() []string
		Created() time.Time
	}
	type Text interface {
		Pages() int
		Words() int
		PageSize() int
	}
	type Audio interface {
		Stream() (io.ReadCloser, error)
		RunningTime() time.Duration
		Format() string // e.g., "MP3", "WAV"
	}
	type Video interface {
		Stream() (io.ReadCloser, error)
		RunningTime() time.Duration
		Format() string // e.g., "MP4", "WMV"
		Resolution() (x, y int)
	}

	type Streamer interface {
		Stream() (io.ReadCloser, error)
		RunningTime() time.Duration
		Format() string
	}
}
//四，flag.Value接口
func test_flag()  {
	//var period = flag.Duration("period", 1*time.Second, "sleep period")
	var period = flag.Duration("period", 50 * time.Millisecond, "sleep period")
	//var period = flag.Duration("period", 2 * time.Minute + 30 *time.Second, "sleep period")
	//var period = flag.Duration("period", 1 * time.Hour + 30 * time.Minute, "sleep period")
	//var period = flag.Duration("period", 24 * time.Hour, "sleep period")

	flag.Parse()
	fmt.Printf("Sleeping for %v...\n", *period)
	time.Sleep(*period)

	type Value interface {
		String() string
		Set(string) error
	}
	fmt.Println("================================")

	//var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")
	//var temp = tempconv.CelsiusFlag("temp", -18.0, "the temperature")
	//flag.Parse()
	//fmt.Println(*temp)

	//练习 7.6： 对tempFlag加入支持开尔文温度。
	var tempF = tempconv.FahrenheitFlag("temp", -18.0, "the temperature")

	flag.Parse()
	fmt.Println(*tempF)

	//练习 7.7： 解释为什么帮助信息在它的默认值是20.0没有包含°C的情况下输出了°C。
	//因为CelsiusFlag实现了set接口，一个*Celsius类型赋给了flag，flag实现的stringter接口
	//最终使Celsius调用了自身实现的string方法，从而将Celsius的值转成带°C的字符串
}
//五，接口值
//1.有两部分组成:一个具体的类型，一个此类型的值
//2.也被称为动态类型和动态值
//3.一个接口值可以持有任意大的动态值
//4.一个接口上的调用必须使用动态分配
//5.接口值得动态类型如果是可以比较的，即可以作为map的key或者switch的语句操作数
//6.
func test_interface_value()  {

	var w io.Writer
	fmt.Println(w,"---0")
	w = os.Stdout
	fmt.Println(w,"---1")
	w = new(bytes.Buffer)
	fmt.Println(w,"---2")
	w = nil
	fmt.Println(w,"---3")

}
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
	//test_interface_type()

	//三，实现接口的条件
	//test_interface_condition()

	//四，flag.Value接口
	//test_flag()

	//五，接口值
	test_interface_value()
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
}

