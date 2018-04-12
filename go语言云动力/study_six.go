package main


/*
	函数
*/
import (
	"fmt"
	"log"
	"math"
	"time"
)
//签名
var pr = fmt.Println

var f = func() {
	pr("Hello")
}

var g func()

//参数
func say(s string)  {
	fmt.Println(s)
}

//返回语句
func tempFile_test(s string, a  int) (result string)  {

	resu := "result"

	if a <= 0 {

		return resu
	}

	resu = s + "_test 123"

	return resu
}

//函数调用
func goo(oo int) int  {
	oo--
	return oo
}

func test()  {
	o := 42
	c := goo(o)
	c++
	fmt.Println(c, o)
}

//闭包
func closure()  {
	add := func(base int) func(int) int {
		return func(n int) int {
			return base + n
		}
	}

	add5 := add(5)
	fmt.Println(add5(10))
}

//压后
func after()  {
	for i := 0; i <= 3; i++ {
		defer fmt.Print(i)
	}
}

//派错和恢复
func protect(gg func())  {
	defer func () {
		log.Panicln("done")
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()
	log.Println("start")
	gg()
}

//方法
type Point struct {
	x, y float64
}
type Rect struct {
	min, max Point
}

func (r *Rect) Area() float64  {
	w := r.max.x - r.min.x
	h := r.max.y - r.min.y
	return w * h
}

func (r Point) Area() float64  {
	return 0
}

func areaTest()  {
	var r Rect
	r.max = Point{2, 2}
	fmt.Println(r.Area())
	fmt.Println(r.Area(), (&r).Area(), (*Rect).Area(&r))
	p := &r.min
	fmt.Println(p.Area(), (*p).Area(), Point.Area(r.min))
}

//主程序入口
func main() {
	//签名
	g = f
	g()

	//参数
	say("Hi")
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Printf("GO 发布于 %s \n", t.Local())

	//返回语句
	s := tempFile_test("shishi", 3)
	fmt.Println(s)

	//函数调用
	test()

	//闭包
	closure()

	//压后
	after()

	//派错和恢复
	//var s []byte
	//protect(func() { s[0] = 0 })
	//protect(func() { panic(42) })
	//s[0] = 42

	//方法
	areaTest()

	//包 + 导入
	abs := math.Abs(-34)
	fmt.Print(abs)

}



