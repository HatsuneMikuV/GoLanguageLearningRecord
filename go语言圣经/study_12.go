package main

import (
	"20180408/display"
	"20180408/eval"
	"20180408/format"
	"fmt"
	"os"
	"reflect"
	"time"
)

/*
	第十二章　反射

	Go语言提供了一种反射机制，
	能够在运行时更新变量和检查它们的值、调用它们的方法和它们支持的内在操作，
	而不需要在编译时就知道这些变量的具体类型。

	反射是一个复杂的内省技术，
	不应该随意使用

	fmt包提供的字符串格式功能，encoding/json和encoding/xml，text/template和html/template包
	这些包内部都是用反射技术实现的，但是它们自己的API都没有公开反射相关的接口
*/


//一，为何需要反射?
//1.为了处理不满足公共接口的类型的值，或者设计函数还不存在的类型


//二，reflect.Type和reflect.Value
//1.reflect包提供了反射功能，定义两个类型Type和Value
//2.Type表示一个Go类型. 它是一个接口
//3.一个Value，有很多方法来检查其内容, 无论具体类型是什么

func test_format_reflect()  {

	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(format.Any(x))                  // "1"
	fmt.Println(format.Any(d))                  // "1"
	fmt.Println(format.Any([]int64{x}))         // "[]int64 0x8202b87b0"
	fmt.Println(format.Any([]time.Duration{d})) // "[]time.Duration 0x8202b87e0"
}

//三，Display，一个递归的值打印器
//1.reflect包提供了反射功能，定义两个类型Type和Value
//2.Type表示一个Go类型. 它是一个接口
//3.一个Value，有很多方法来检查其内容, 无论具体类型是什么

func test_Display()  {

	e, _ := eval.Parse("sqrt(A / pi)")
	display.Display("e", e)


	fmt.Print("\n\n=====================\n\n")

	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	display.Display("strangelove", strangelove)

	fmt.Print("\n\n=====================\n\n")

	display.Display("os.Stderr", os.Stderr)

	fmt.Print("\n\n=====================\n\n")

	display.Display("rV", reflect.ValueOf(os.Stderr))

	fmt.Print("\n\n=====================\n\n")

	var i interface{} = 3

	display.Display("i", i)
	display.Display("&i", &i)

	fmt.Print("\n\n=====================\n\n")
	//4.对象图中含有回环，Display将会陷入死循环
	// a struct that points to itself
	type Cycle struct{ Value int; Tail *Cycle }
	var c Cycle
	c = Cycle{42, &c}
	display.Display("c", c)
}

//四，示例: 编码为S表达式
//1.S表达式格式，采用Lisp语言的语法

func main() {
	//二，reflect.Type和reflect.Value
	//test_format_reflect()

	//三，Display，一个递归的值打印器
	test_Display()
}
