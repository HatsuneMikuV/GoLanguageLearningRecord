package main

import (
	"20180408/geometry"
	"fmt"
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

func main() {

	//方法声明
	test_fun()
}
