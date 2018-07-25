package main

import (
	"fmt"
	"unsafe"
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

func main() {

	//一，unsafe.Sizeof, Alignof 和 Offsetof
	test_one_unsafe()

	//二，unsafe.Pointer
	test_two_Pointer()
}
