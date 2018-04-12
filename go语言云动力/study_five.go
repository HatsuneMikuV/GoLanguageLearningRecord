package main

import (
	"fmt"
	"net/url"

	//"net/url"
)

//字符串转换
func string_change()  {
	var s = "Go 程"
	fmt.Printf("%c %c",s[0], s[2])
}

func string_change_one()  {
	var s = "Go 程"
	var r = []rune(s)

	fmt.Printf("%c %c",r[0], r[2])

	fmt.Printf("% x \n", r)
}

//数组
func array_one()  {
	type Point [2]float32
	type Line [2]Point

	a := Point{1, 2}
	b := Point{3, 4}
	l := Line{a, b}

	fmt.Printf("%g, %g, %v, %v \n", a[0], l[1][1], b, l)
}

func array_two()  {
	type Point [2]float32
	type Line [2]Point
	a := Point{1, 2}
	b := a
	c := Line{a, b}

	b[0] = 42

	s := "%v, \n%v, \n%v"

	fmt.Printf(s, a, b, c)
}

//切片
func slice_one()  {
	s := [4]int{0, 1, 2, 3}
	t := s[1:3]
	fmt.Println(s[0])
	fmt.Println(t)
	fmt.Println(s[:3])
	fmt.Println(t[1:])
	fmt.Println(len(s), cap(s))
	fmt.Println(len(t), cap(t))
}

//结构体
func struct_one()  {
	gurl, er := url.Parse("http://golang.org/pkg")

	fmt.Println(gurl.Host, er)
}

func struct_two()  {
	var d2 = D2{1, 2}
	d3 := D3{d2, 3}

	fmt.Println(d2.x, d3.D2, d3.x, d3.z)

	fmt.Println(d3.D2 == d2)
}

//指针
func point_one()  {
	var i int
	var p *int
	var pp **int

	i = 0
	p = &i
	pp = &p
	*p++

	fmt.Println(i, p, *p, pp, *pp)
}

func main()  {

	//字符串转换
	string_change()
	string_change_one()

	//数组
	array_one()
	array_two()

	//切片
	slice_one()

	//结构体
	struct_one()
	struct_two()

	//指针
	point_one()
}

type D2 struct {x, y float64}

type D3 struct {
	D2
	z float64
}