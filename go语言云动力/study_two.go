package main

import (
	"fmt"
	"sort"
)

//猜数字游戏
func Guess_num()  {
	min, max := 0, 100
	fmt.Printf("请想一个%d~%d的整数。\n", min, max)

	for min < max  {
		i := (min + max) / 2
		fmt.Printf("该数小于或者等于%d吗？（y/n）", i)
		var  s string
		fmt.Scanf("%s", &s)
		if s != "" && s[0] == 'y' {
			max = i
		}else {
			min = i + 1
		}
	}
	fmt.Printf("该数是%d\n", max)
}

//二分查找
func Binarysearch()  {

	fmt.Println("Pick a number from 0 to 100")
	fmt.Printf("Your number is %d\n", sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d?", i)
		var s string
		fmt.Scanf("%s\n", &s)
		return s != "" && s[0] == 'y'
	}))
}

func main()  {

	//猜数字游戏
	Guess_num()

	//go 语言  提供的二分查找功能
	Binarysearch()
}