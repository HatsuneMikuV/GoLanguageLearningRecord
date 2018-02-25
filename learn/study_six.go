package main

import (
	"fmt"
	//"go/ast"
	"log"
	//"runtime"
	//"sync"
	//"unicode"
	//
	//"go/ast"
	//"time"
)

var pr = fmt.Println

var f = func() {
	pr("Hello")
}

var g func()

func main() {
	//g = f
	//g()
	//
	//say("Hi")
	//
	//t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	//
	//fmt.Printf("GO 发布于 %s \n", t.Local())
	//
	//test()

	//add := func(base int) func(int)int {
	//	return func(n int) int {
	//		return base + n
	//	}
	//}
	//
	//add5 := add(5)
	//
	//fmt.Println(add5(10))

	//var s []byte
	//protect(func() { s[0] = 0 })
	//protect(func() { panic(42) })
	//
	//s[0] = 42

	areaTest()
}

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

func protect(g func())  {
	defer func () {
		log.Panicln("done")
		if x := recover(); x != nil {
			log.Printf("run time panic: %v", x)
		}
	}()
	log.Println("start")
	g()
}

func say(s string)  {
	fmt.Println(s)
}

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
