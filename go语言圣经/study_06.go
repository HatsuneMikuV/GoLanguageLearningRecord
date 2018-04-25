package main

import (
	"20180408/geometry"
	"fmt"
	"image/color"
	"net/url"
)

/* 方法 */


//一，方法声明
//1.函数声明的时候，在其前面加上一个变量，即是一个方法
//2.这个变量将函数附加到这个类型上，即方法成了此变量的独占方法

func test_fun()  {

	p := geometry.Point{1, 2}
	q := geometry.Point{4, 6}
	fmt.Println(geometry.Distance(p, q))//函数调用
	fmt.Println(p.Distance(q))//方法调用

	//4.可以给同一个包内的任意命名类型定义方法，只要这个命名类型的底层类型
	perim := geometry.Path{
		{1, 1},
		{5, 1},
		{5, 4},
		{1, 1},
	}
	fmt.Println(perim.Distance())

	//5.对于既定的类型，内部的方法名必须唯一，而不同类型，则可出现同名方法
	//6.使用方法，能让我们减少输入，也不必带上包名这么麻烦的事
	perimm := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
	fmt.Println(geometry.PathDistance(perimm)) // "12", standalone function
	fmt.Println(perim.Distance())
}
//二，基于指针对象的方法
//1.当变量比较大时，为了节省内存是可以选择指针变量，这样就不会复制变量本身
func test_point()  {

	r := &geometry.Point{1, 2}
	r.ScaleBy(2)
	fmt.Println(*r)

	p := geometry.Point{1, 2}
	pptr := &p
	pptr.ScaleBy(2)
	fmt.Println(p)

	pp := geometry.Point{1, 2}
	(&pp).ScaleBy(2)
	fmt.Println(pp)

	//2.编译器会隐式地帮我们用&p去调用ScaleBy这个方法
	//3.不管你的method的receiver是指针类型还是非指针类型，都是可以通过指针/非指针类型进行调用的，编译器会帮你做类型转换
	//4.声明方法是否为指针类型取决于变量本身的大小，因为非指针类型变量会拷贝一份，而指针类型要注意内存地址，永远指向一份
	pp.ScaleBy(3)
	fmt.Println(pp)

	//5.Nil也是一个合法的接收器类型
	//6.传入的是存储了内存地址的变量，你改变这个变量是影响不了原始的变量的
	m := url.Values{"lang": {"en"}}
	m.Add("item", "1")
	m.Add("item", "2")

	fmt.Println("------------------")
	fmt.Println(m.Get("lang"))
	fmt.Println(m.Get("q"))
	fmt.Println(m.Get("item"))
	fmt.Println(m["item"])

	fmt.Println("------------------")
	m = nil
	fmt.Println(m.Get("item"))
	//m.Add("item", "3")//panic: assignment to entry in nil map
}


//三，通过嵌入结构体来扩展类型
//1.内嵌可以使定义时的简写形式，并使其包含嵌入类型所具有的一切字段，类似OC的继承
func test_extion()  {

	var cp geometry.ColoredPoint
	cp.X = 1
	fmt.Println(cp.Point.X) // "1"
	cp.Point.Y = 2
	fmt.Println(cp.Y) // "2"
	fmt.Println("------------------")

	//2.内嵌可以使多字段类型定义时，划分成多个小类型，然后再定义小类型方法，方便于读写
	//	可以把多字段类型认为是小类型的继承或者子类，就像OC那样，但是3...
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	var p = geometry.ColoredPoint{geometry.Point{1, 1}, red}
	var q = geometry.ColoredPoint{geometry.Point{5, 4}, blue}
	fmt.Println(p.Distance(q.Point)) // "5"
	p.ScaleBy(2)
	q.ScaleBy(2)
	fmt.Println(p.Distance(q.Point))
	fmt.Println("------------------")
	//3.但是并不能利用隐式的方式调用小类型的方法，因为类型不同，必须显式的调用
	//p.Distance(q) // compile error: cannot use q (ColoredPoint) as Point

	p1 := geometry.ColoredPointP{&geometry.Point{1, 1}, red}
	q1 := geometry.ColoredPointP{&geometry.Point{5, 4}, blue}
	fmt.Println(p1.Distance(*q1.Point)) // "5"
	q1.Point = p1.Point                 // p and q now share the same Point
	p1.ScaleBy(2)
	fmt.Println(*p1.Point, *q1.Point)
}

//四，方法值和方法表达式
func test_func_value()  {

	p := geometry.Point{1, 2}
	q := geometry.Point{4, 6}

	//p.Distance作为选择器，distanceFromP一个方法"值"
	//一个将方法(Point.Distance)绑定到特定接收器变量的函数
	distanceFromP := p.Distance        // method value
	fmt.Println(distanceFromP(q))      // "5"
	var origin geometry.Point                   // {0, 0}
	fmt.Println(distanceFromP(origin)) // "2.23606797749979", sqrt(5)

	scaleP := p.ScaleBy // method value
	scaleP(2)           // p becomes (2, 4)
	scaleP(3)           //      then (6, 12)
	scaleP(10)          //      then (60, 120)

	{//{}目的和上面的变量区分，防止冲突

		// 译注：这个Distance实际上是指定了Point对象为接收器的一个方法func (p Point) Distance()，
		// 但通过Point.Distance得到的函数需要比实际的Distance方法多一个参数，
		// 即其需要用第一个额外参数指定接收器，后面排列Distance方法的参数。
		// 看起来本书中函数和方法的区别是指有没有接收器，而不像其他语言那样是指有没有返回值。
		p := geometry.Point{1, 2}
		q := geometry.Point{4, 6}

		distance := geometry.Point.Distance   // method expression
		fmt.Println(distance(p, q))  // "5"
		fmt.Printf("%T\n", distance) // "func(Point, Point) float64"

		scale := (*geometry.Point).ScaleBy
		scale(&p, 2)
		fmt.Println(p)            // "{2 4}"
		fmt.Printf("%T\n", scale) // "func(*Point, float64)"
	}
}

//五，示例: Bit数组
//1.在数据流分析领域，集合元素通常是一个非负整数，集合会包含很多元素，
//并且集合会经常进行并集、交集操作，这种情况下，bit数组会比map表现更加理想
//2.一个bit数组通常会用一个无符号数或者称之为“字”的slice来表示，
//每一个元素的每一位都表示集合里的一个值
func test_bit()  {

	var x, y geometry.IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"

	//并集：将y中没有的元素放到x中，新x相当于x ∪ y
	x.UnionWith(&y)
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x.Has(8), x.Has(123))

	fmt.Println(&x)         // "{1 9 42 144}"
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x)          // "{[4398046511618 0 65536]}"

	//练习6.1: 为bit数组实现下面这些方法
	fmt.Println(x.Len())
	fmt.Println("======")

	x.Remove(9)
	fmt.Println(&x)
	fmt.Println("======")

	x.Clear()
	fmt.Println(&x)
	fmt.Println("======")

	xx := x.Copy()
	fmt.Println(xx)
	fmt.Println("======")

	//练习 6.2： 定义一个变参方法(*IntSet).AddAll(...int)
	x.AddAll(1, 2, 3, 42)
	fmt.Println(&x)
	fmt.Println("======")

	//练习 6.3： (*IntSet).UnionWith会用|操作符计算两个集合的交集，
	//我们再为IntSet实现另外的几个函数
	//IntersectWith(交集：元素在A集合B集合均出现),
	//AndWith(并集：元素在A集合或B集合中出现),
	//DifferenceWith(差集：元素出现在A集合，未出现在B集合),
	//SymmetricDifference(并差集：元素出现在A但没有出现在B，或者出现在B没有出现在A)。
	y.AddAll(1, 2)
	fmt.Println(&x)
	fmt.Println(&y)
	fmt.Println("======↑")

	//交集x & y    数学公式：x ∩ y
	fmt.Println(x.IntersectWith(&y))
	fmt.Println("======交集↑")

	//并集x | y    数学公式：x ∪ y
	fmt.Println(x.AndWith(&y))
	fmt.Println("======并集↑")

	//差集x &^ y   数学公式：x - y
	fmt.Println(x.DifferenceWith(&y))
	fmt.Println("======差集↑")

	//并差集x ^ y  数学公式：￢x ∪ ￢y
	fmt.Println(x.SymmetricDifference(&y))
	fmt.Println("======并差集↑")

	//练习6.4: 实现一个Elems方法，返回集合中的所有元素，用于做一些range之类的遍历操作。
	//返回集合中的所有元素
	fmt.Println(x.Elems())
	fmt.Println("======集合所有元素↑")

	//计算机位数平台的自动判断的一个智能表达式：32 << (^uint(0) >> 63)
	fmt.Println((32 << (^uint(0) >> 63)))

}

//六，封装

func main() {

	//一，方法声明
	//test_fun()

	//二，基于指针对象的方法
	//test_point()

	//三，通过嵌入结构体来扩展类型
	//test_extion()

	//四，方法值和方法表达式
	//test_func_value()

	//五，示例: Bit数组
	//test_bit()
}

