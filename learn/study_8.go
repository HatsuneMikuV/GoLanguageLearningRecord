package main

import "fmt"

/*
	面向对象
*/

//继承

type O struct {
	h int
}

func (a *O) Get() int  {
	return a.h
}

func (a *O) Set(h int)  {
	a.h = h
}

type OO struct {
	O
	i int
}

func (a *OO) Get() int  {
	a.i++
	return a.O.Get()
}

func oo_test()  {
	oo := new(OO)
	oo.Set(42)
	fmt.Println(oo, oo.Get())
}


func main() {

	//继承
	oo_test()
}
