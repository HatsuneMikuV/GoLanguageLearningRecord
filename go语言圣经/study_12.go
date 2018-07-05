package main

import (
	"20180408/display"
	"20180408/eval"
	"20180408/format"
	"20180408/sexpr"
	"encoding/json"
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
//2.Go语言自带的标准库并不支持S表达式，主要是因为它没有一个公认的标准规范
/*
	42          integer
	"hello"     string (带有Go风格的引号)
	foo         symbol (未用引号括起来的名字)
	(1 2 3)     list   (括号包起来的0个或多个元素)
*/
func test_s_fun()  {

	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
		//CC				complex128
		Inface			interface{}
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		//CC:       complex(1, 2),
		Inface:   []int{1, 2, 3},
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

	// Encode it
	data, err := sexpr.Marshal(strangelove)
	if err != nil {
		fmt.Printf("Marshal failed: %v\n", err)
		fmt.Print("\n\n=====================\n\n")
	}
	fmt.Printf("Marshal() = %s\n", data)
	fmt.Print("\n\n=====================\n\n")


	/*
		练习 12.5： 修改encode函数，用JSON格式代替S表达式格式。
		然后使用标准库提供的json.Unmarshal解码器来验证函数是正确的。

		基本上是一致，结果是正确的---------------------
	*/
	// Encode_json it
	data_json, err_json := sexpr.Marshal_Json(strangelove)
	if err_json != nil {
		fmt.Printf("Marshal_Json failed: %v\n", err_json)
		fmt.Print("\n\n=====================\n\n")
	}
	fmt.Printf("Marshal_Json() = %s\n", data_json)
	fmt.Print("\n\n=====================\n\n")

	// json it
	data_j, err_j := json.Marshal(strangelove)
	if err_j != nil {
		fmt.Printf("json.Marshal failed: %v\n", err_j)
		fmt.Print("\n\n=====================\n\n")
	}
	fmt.Printf("json.Marshal() = %s\n", data_j)
	fmt.Print("\n\n=====================\n\n")
	//---------------------


	// Decode it
	var movie Movie
	if err := sexpr.Unmarshal(data, &movie); err != nil {
		fmt.Printf("Unmarshal failed: %v\n", err)
		fmt.Print("\n\n=====================\n\n")
	}
	fmt.Printf("Unmarshal() = %+v\n", movie)
	fmt.Print("\n\n=====================\n\n")


	// Check equality.
	if !reflect.DeepEqual(movie, strangelove) {
		fmt.Printf("not equal\n")
		fmt.Print("\n\n=====================\n\n")
	}

	// Pretty-print it:
	data, err = sexpr.MarshalIndent(strangelove)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("MarshalIdent() = \n%s\n", data)
	fmt.Print("\n\n=====================\n\n")
}

//五，通过reflect.Value修改值
//1.变量可以通过寻址空间更新存储的值
//2.通过调用reflect.Value的CanAddr方法来判断其是否可以被取地址
//3.知道变量的类型，可以使用类型的断言机制将得到的interface{}类型的接口强制转为普通的类型指针，即可改值
//4.通过调用可取地址的reflect.Value的reflect.Value.Set方法来更新值
//5.一个可取地址的reflect.Value会记录一个结构体成员是否是未导出成员，
// 如果是的话则拒绝修改操作。因此，CanAddr方法并不能正确反映一个变量是否是可以被修改的
//6.另一个相关的方法CanSet是用于检查对应的reflect.Value是否是可取地址并可被修改的

func test_reflect_Value()  {

	x := 2                   // value   type    variable?
	a := reflect.ValueOf(2)  // 2       int     no
	b := reflect.ValueOf(x)  // 2       int     no
	c := reflect.ValueOf(&x) // &x      *int    no
	d := c.Elem()            // 2       int     yes (x)

	fmt.Print(x, "\n", a, "\n", b, "\n", c, "\n", d, "\n")

	fmt.Println(a.CanAddr()) // "false"
	fmt.Println(b.CanAddr()) // "false"
	fmt.Println(c.CanAddr()) // "false"
	fmt.Println(d.CanAddr()) // "true"

	px := d.Addr().Interface().(*int) // px := &x
	*px = 3                           // x = 3
	fmt.Println(x)                    // "3"

	d.Set(reflect.ValueOf(4))
	fmt.Println(x) // "4"

	//d.Set(reflect.ValueOf(int64(5))) // panic: reflect.Set: value of type int64 is not assignable to type int
	fmt.Print(x, "\n")

	//b.Set(reflect.ValueOf(3)) // panic: reflect: reflect.Value.Set using unaddressable value
	fmt.Print(x, "\n")

	d.SetInt(3)
	fmt.Println(x, "\n") // "3"

	xx := 1
	rx := reflect.ValueOf(&xx).Elem()
	rx.SetInt(2)                     // OK, x = 2
	rx.Set(reflect.ValueOf(3))       // OK, x = 3
	//rx.SetString("hello")            // panic: reflect: call of reflect.Value.SetString on int Value
	//rx.Set(reflect.ValueOf("hello")) // panic: reflect.Set: value of type string is not assignable to type int


	var y interface{}
	ry := reflect.ValueOf(&y).Elem()
	//ry.SetInt(2)                     // panic: reflect: call of reflect.Value.SetInt on interface Value
	ry.Set(reflect.ValueOf(3))       // OK, y = int(3)
	//ry.SetString("hello")            // panic: reflect: call of reflect.Value.SetString on interface Value
	ry.Set(reflect.ValueOf("hello")) // OK, y = "hello"
}


func main() {
	//二，reflect.Type和reflect.Value
	//test_format_reflect()

	//三，Display，一个递归的值打印器
	//test_Display()

	//四，示例: 编码为S表达式
	//test_s_fun()

	//五，通过reflect.Value修改值
	test_reflect_Value()
}
