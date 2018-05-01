package main

import (
	"fmt"
	"20180408/tempconv"
)

/* 程序结构 */

//一，命名
//二，声明
const boilingF  = 212.0

func test_one02()  {
	var f  = boilingF
	var c  = (f - 32) * 5 / 9
	fmt.Printf("boiling point = %g°F or %g°C\n", f, c)
}
func fToC(f float64) float64 {
	return (f - 32) * 5 / 9
}
func test_two02()  {
	const freezingF, boilingFF  = 32.0, 212.0
	fmt.Printf("%g°F = %g°C\n", freezingF, fToC(freezingF))
	fmt.Printf("%g°F = %g°C\n", boilingFF, fToC(boilingFF))
}

//三，变量
func test_thr02()  {
	var s string
	fmt.Println(s)

	var i, j, k int
	var b, f, c = true, 2.3, "four"

	fmt.Println(i, j ,k ,b, f, c)
}

//四，赋值
func test_fou02()  {

	var x int
	var pp bool
	var p *bool = &pp
	type Person struct {
		name string
	}
	var person Person
	var count02 = []int {1, 2, 3}

	x = 1
	*p = true
	person.name = "bob"
	scale := 100
	count02[1] = count02[1] * scale

	fmt.Println(x, p, person, count02)
}

//五，类型
//六，包和文件   tempconv
func test_fiv02()  {
	fmt.Printf("%g°C\n", tempconv.BoilingCCC - tempconv.FreezingCCC) // "100" °C
	boilingF := tempconv.CToF(tempconv.BoilingCCC)
	fmt.Printf("%s\n",boilingF)
	fmt.Printf("%g°F\n", boilingF - tempconv.CToF(tempconv.FreezingCCC)) // "180" °F
	//fmt.Printf("%g\n", boilingF - FreezingCCC) // compile error: type mismatch

	var cc tempconv.Celsius
	var f tempconv.Fahrenheit
	fmt.Println(cc == 0)	// "true"
	fmt.Println(f >= 0)	// "true"
	//fmt.Println(cc == f) // compile error: type mismatch
	fmt.Println(cc == tempconv.Celsius(f)) // "true"!

	c := tempconv.FToC(212.0)
	fmt.Println(c.String()) // "100°C"
	fmt.Printf("%v\n", c) // "100°C"; no need to call String explicitly
	fmt.Printf("%s\n", c) // "100°C"
	fmt.Println(c) // "100°C"
	fmt.Printf("%g\n", c) // "100"; does not call String
	fmt.Println(float64(c)) // "100"; does not call String

	k := tempconv.FToK(212.0)
	fmt.Println(k.String())
	fmt.Printf("%s\n", k)
	fmt.Println(k)
	fmt.Printf("%g\n", k)

	zero := tempconv.CToK(tempconv.AbsoluteZeroC)
	fmt.Println(zero)
}

//七，作用域
func ff() {}
var g = "g"
func test_six02()  {
	f := "f"
	fmt.Println(f)// "f"; local var f shadows package-level func f
	fmt.Println(g)// "g"; package-level var
	//fmt.Println(h)// compile error: undefined: h

	x := "hello"
	for _, x := range x{
		x := x + 'A' - 'a'
		fmt.Printf("%c", x)
	}
}

func main() {

	//一，命名
	//二，声明
	test_one02()
	test_two02()

	//三，变量
	test_thr02()

	//四，赋值
	test_fou02()

	//五，类型
	//六，包和文件   tempconv
	test_fiv02()

	//七，作用域
	test_six02()
}
