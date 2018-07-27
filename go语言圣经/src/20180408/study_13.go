package main

import (
	"fmt"
	"unsafe"
	"20180408/myequal"
	"strings"
	"reflect"
	"20180408/bzip2"
	"os"
	"io"
	"log"
)

/*
	第十三章　底层编程

	1.Go语言的设计包含了诸多安全策略，限制了可能导致程序运行出错的用法
	2.字符串、map、slice和chan等所有的内置类型，都有严格的类型转换规则
	3.Go程序相比较低级的C语言来说更容易预测和理解，程序也不容易崩溃
	4.使用unsafe包来摆脱Go语言规则带来的限制，讲述如何创建C语言函数库的绑定，以及如何进行系统调用
	5.unsafe包被广泛地用于比较低级的包, 例如runtime、os、syscall还有net包等，因为它们需要和操作系统密切配合
*/

//一，unsafe.Sizeof, Alignof 和 Offsetof
//1.unsafe.Sizeof函数可以对任意表达式计算内存大小
//2.计算机在加载和保存数据时，如果内存地址合理地对齐的将会更有效率
//3.内存空洞是编译器自动添加的没有被使用的内存空间
//4.内存空洞可能会存在一些随机数据，可能会对用unsafe包直接操作内存的处理产生影响
/*

	bool							1个字节
	intN, uintN, floatN, complexN	N/8个字节(例如float64是8个字节)
	int, uint, uintptr				1个机器字
	*T								1个机器字
	string							2个机器字(data,len)
	[]T								3个机器字(data,len,cap)
	map								1个机器字
	func							1个机器字
	chan							1个机器字
	interface						2个机器字(type,value)
*/
//5.unsafe.Alignof 函数返回对应参数的类型需要对齐的倍数
//6.unsafe.Offsetof 函数的参数必须是一个字段 x.f,
// 然后返回 f 字段相对于 x 起始地址的偏移量, 包括可能的空洞
//备注：	虽然这几个函数在不安全的unsafe包，
// 		但是这几个函数调用并不是真的不安全，
// 		特别在需要优化内存空间时它们返回的结果对于理解原生的内存布局很有帮助
func test_one_unsafe()  {
	fmt.Println("==========unsafe.Sizeof===========")
	fmt.Println(unsafe.Sizeof(float64(0)))

	type TestOne struct {
		A       bool
		B      	float64
		C   	int16
	}

	type TestTwo struct {
		A       float64
		B      	int16
		C   	bool
	}

	type TestThi struct {
		A       bool
		B      	int16
		C   	float64
	}

	var one TestOne
	var two TestTwo
	var thi TestThi
	fmt.Println("=========unsafe.Sizeof============")

	fmt.Println(unsafe.Sizeof(one))
	fmt.Println(unsafe.Sizeof(two))
	fmt.Println(unsafe.Sizeof(thi))

	var x struct {
		a bool
		b int16
		c []int
	}
	fmt.Println("=========Sizeof============")

	fmt.Println(unsafe.Sizeof(x))
	fmt.Println(unsafe.Sizeof(x.a))
	fmt.Println(unsafe.Sizeof(x.b))
	fmt.Println(unsafe.Sizeof(x.c))

	fmt.Println("=========Alignof============")

	fmt.Println(unsafe.Alignof(x))
	fmt.Println(unsafe.Alignof(x.a))
	fmt.Println(unsafe.Alignof(x.b))
	fmt.Println(unsafe.Alignof(x.c))

	fmt.Println("=========Offsetof============")

	fmt.Println(unsafe.Offsetof(x.a))
	fmt.Println(unsafe.Offsetof(x.b))
	fmt.Println(unsafe.Offsetof(x.c))
}

//二，unsafe.Pointer
//1.unsafe.Pointer是特别定义的一种指针类型（译注：类似C语言中的void*类型的指针）,可以包含任意类型变量的地址
//2.unsafe.Pointer指针也是可以比较的，并且支持和nil常量比较判断是否为空指针
//3.许多将unsafe.Pointer指针转为原生数字，然后再转回为unsafe.Pointer类型指针的操作也是不安全的

func test_two_Pointer()  {

	fmt.Println("=========Pointer============")

	fmt.Printf("%#016x\n", Float64bits(1.0))


	fmt.Println("=========Pointer============")
	var x struct {
		a bool
		b int16
		c []int
	}

	// 和 pb := &x.b 等价
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b) // "42"

	fmt.Println("=========Pointer============")

	// NOTE: subtly incorrect!
	tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	pbb := (*int16)(unsafe.Pointer(tmp))
	*pbb = 42
	fmt.Println(x.b) // "42"


}
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}


//三，示例: 深度相等判断
//1.DeepEqual函数使用内建的==比较操作符对基础类型进行相等判断，
// 对于复合类型则递归该变量的每个基础类型然后做类似的比较判断,
// 甚至对于一些不支持==操作运算符的类型也可以工作，
// 因此在一些测试代码中广泛地使用该函数
//2.

func test_ex_ch()  {

	got := strings.Split("a:b:c", ":")
	want := []string{"a", "b", "c"};
	if reflect.DeepEqual(got, want) {
		fmt.Println("===========000============")
	}

	fmt.Println("===========111============")
	var aa, bb []string = nil, []string{}
	fmt.Println(reflect.DeepEqual(aa, bb)) // "false"

	var cc, dd map[string]int = nil, make(map[string]int)
	fmt.Println(reflect.DeepEqual(cc, dd)) // "false"


	fmt.Println("===========222============")
	fmt.Println(myequal.Equal([]int{1, 2, 3}, []int{1, 2, 3}))        // "true"
	fmt.Println(myequal.Equal([]string{"foo"}, []string{"bar"}))      // "false"
	fmt.Println(myequal.Equal([]string(nil), []string{}))             // "true"
	fmt.Println(myequal.Equal(map[string]int(nil), map[string]int{})) // "true"


	// Circular linked lists a -> b -> a and c -> c.
	type link struct {
		value string
		tail *link
	}
	a, b, c := &link{value: "a"}, &link{value: "b"}, &link{value: "c"}
	a.tail, b.tail, c.tail = b, a, c

	fmt.Println("===========333============")
	fmt.Println(myequal.Equal(a, a)) // "true"
	fmt.Println(myequal.Equal(b, b)) // "true"
	fmt.Println(myequal.Equal(c, c)) // "true"
	fmt.Println(myequal.Equal(a, b)) // "false"
	fmt.Println(myequal.Equal(a, c)) // "false"

	fmt.Println("===========444============")
	//练习 13.1： 定义一个深比较函数，对于十亿以内的数字比较，忽略类型差异。
	x := int(999999999)
	y := float64(999999999)
	//y := "1"  //panic: 1  is not num type
	fmt.Println(myequal.NumDeepEqual(x, y)) // "true"

	fmt.Println("===========555============")
	//练习 13.2： 编写一个函数，报告其参数是否为循环数据结构。
	ax, bx := &link{value: "a"}, &link{value: "b"}
	ax.tail, bx.tail = bx, ax
	//Cycle linked lists ax -> bx -> ax
	fmt.Println(myequal.CycleCheck(ax))// "true"

	cx, dx := &link{value: "c"}, &link{value: "d"}
	cx.tail = dx
	//Cycle linked lists cx -> dx
	fmt.Println(myequal.CycleCheck(cx))// "false"

	qx, wx, ex := &link{value: "q"}, &link{value: "w"}, &link{value: "e"}
	qx.tail, wx.tail, ex.tail = wx, ex, wx
	//Cycle linked lists qx -> wx -> ex -> wx
	fmt.Println(myequal.CycleCheck(qx))// "true"
}

//四，通过cgo调用C代码
//1.Go程序可能会遇到要访问C语言的某些硬件驱动函数的场景，
// 或者是从一个C++语言实现的嵌入式数据库查询记录的场景，
// 或者是使用Fortran语言实现的一些线性代数库的场景
//2.将Go编译为静态库然后链接到C程序，或者将Go程序编译为动态库然后在C程序中动态加载也都是可行的
func test_cgo()  {

	file,_:=os.OpenFile("words.txt",os.O_RDWR|os.O_CREATE,0777)
	defer file.Close()
	bs := []byte{97,98,99,100,101,102} // a-f
	file.Write(bs[0:4]) //0123
	file.WriteString("面朝大海")
	file.WriteAt([]byte{'A','B','C','D','E'},13)

	w := bzip2.NewWriter(file)
	if _, err := io.Copy(w, file); err != nil {
		log.Fatalf("bzipper: %v\n", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("bzipper: close: %v\n", err)
	}
}


//五，几点忠告
//1.警告要谨慎使用reflect包
//2.警告同样适用于本章的unsafe包

func main() {

	//一，unsafe.Sizeof, Alignof 和 Offsetof
	//test_one_unsafe()

	//二，unsafe.Pointer
	//test_two_Pointer()

	//三，示例: 深度相等判断
	//test_ex_ch()

	//四，通过cgo调用C代码
	test_cgo()
}
