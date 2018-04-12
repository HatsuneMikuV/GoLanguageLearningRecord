package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

/* 入门 */

//Hello World
func test_one()  {
	//C语言是直接影响Go语言设计的语言之一
	//Go是一门编译型语言
	//Go语言原生支持Unicode，它可以处理全世界任何语言的文本

	//Go语言的代码通过包（package）组织
	//代码会放在$GOPATH/src/...
	//main 是整个程序执行时的入口(C系语言差不多都这样)

	//缺少了必要的包或者导入了不需要的包，程序都无法编译通过

	//gofmt工具把代码格式化为标准格式
	//goimports，可以根据代码需要, 自动地添加或删除import声明
	fmt.Println("Hello, 世界")

}

//命令行参数

//os包以跨平台的方式，提供了一些与操作系统交互的函数和变量
//os.Args变量是一个字符串（string）的切片（slice）
//Go语言提供了常规的数值和逻辑运算符
func test_two()  {

	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}
//Go语言中这种情况的解决方法是用空标识符（blank identifier），即_（也就是下划线）。
//空标识符可用于任何语法需要变量名但程序逻辑不需要的时候
func test_thi()  {
	s, sep := "", ""
	for _, arg := range os.Args[1:]{
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)

	fmt.Println(os.Args[1:])
}

//查找重复的行

//Go的map类似于Java语言中的HashMap，Python语言中的dict，Lua语言中的table，通常使用hash实现
func test_fou()  {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		counts[input.Text()]++
	}
	//note：ignoring potential errors from input.Err()
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
func countLines(f *os.File, counts map[string]int)  {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	//note：ignoring potential errors from input.Err()
}
func test_fiv()  {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	}else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2:%v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
func test_six()  {
	counts := make(map[string]int)

	for _, fileName := range os.Args[1:] {
		data , err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3:%v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n"){
			counts[line]++
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

//GIF 动画
//Go语言标准库里的image这个package的用法，
//我们会用这个包来生成一系列的bit-mapped图，
//然后将这些图片编码为一个GIF动画
var palette = []color.Color{
	color.White,
	color.Black,
	color.RGBA{0xff,0xc0,0xcb, 1.0},
	color.RGBA{0x00,0x80,0x00, 1.0}}

const (
	whiteIndex = 0//first color in palette
	blackIndex = 1//next color in palette
	redIndex = 2//next color in palette
	greenIndex = 3//next color in palette
)
/*
	练习 1.6： 修改Lissajous程序，修改其调色板来生成更丰富的颜色，
	然后修改SetColorIndex的第三个参数，看看显示结果吧。

 	练习 1.5： 修改前面的Lissajous程序里的调色板，由黑色改为绿色。
 	我们可以用color.RGBA{0xRR, 0xGG, 0xBB, 0xff}来得到#RRGGBB这个色值，
 	三个十六进制的字符串分别代表红、绿、蓝像素。
 */

func test_seven(name string, index uint8)  {
	const  (
		cycles = 5 //number of complete x oscillator revolutions
		res = 0.001 //angular resolution
		size = 100 //image canvas covers [-size..size]
		nframes = 64 // number of animation frames
		delay = 8 // delay between frames in 10ms units
	)

	freq := rand.Float64() * 3.0 //relative frequency of y oscillator
	anim := gif.GIF{LoopCount:nframes}
	phase:= 0.9 //phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2 * size + 1, 2 * size + 1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles * 2 * math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t * freq + phase)
			img.SetColorIndex(size + int(x * size + 0.5), size + int(y * size + 0.5), index)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	gif.EncodeAll(file, &anim)//note:ignorning encoding errors
	file.Close()
}
func test_seven_url(out io.Writer, index uint8)  {
	const  (
		cycles = 5 //number of complete x oscillator revolutions
		res = 0.001 //angular resolution
		size = 100 //image canvas covers [-size..size]
		nframes = 64 // number of animation frames
		delay = 8 // delay between frames in 10ms units
	)

	freq := rand.Float64() * 3.0 //relative frequency of y oscillator
	anim := gif.GIF{LoopCount:nframes}
	phase:= 0.9 //phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2 * size + 1, 2 * size + 1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles * 2 * math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t * freq + phase)
			img.SetColorIndex(size + int(x * size + 0.5), size + int(y * size + 0.5), index)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim)//note:ignorning encoding errors
}

//获取URL
func test_eight(url string)  {

	if strings.HasPrefix(url, "http://") == false {
		url = "http://" + url
		fmt.Println("http://被添加")
	}
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("fetch:", err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("fetch: reading", url, err)
		return
	}
	fmt.Println("http状态码", resp.Status)
	fmt.Printf("\n%s\n", b)
}

//并发获取多个URL
func test_nine()  {
	start := time.Now()
	ch := make(chan  string)

	var url = "https://golang.org"
	go fetch(url, ch)
	fmt.Println(<-ch)

	url = "http://gopl.io"
	go fetch(url, ch)
	fmt.Println(<-ch)

	url = "https://godoc.org"
	go fetch(url, ch)
	fmt.Println(<-ch)

	//测试是否缓存
	url = "https://godoc.org"
	go fetch(url, ch)
	fmt.Println(<-ch)

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
func fetch(url string, ch chan <- string)  {
	start := time.Now()
	resp, err := http.Get(url)

	if err != nil {
		ch<- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() //don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s:%v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}

//web服务
var mu sync.Mutex
var count int
func handler(w http.ResponseWriter, r *http.Request)  {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path[1:])
}
func counter(w http.ResponseWriter, r *http.Request)  {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}
func test_ten()  {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler3(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v) }
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err) }
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}
func test_elev()  {
	http.HandleFunc("/", handler3)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func test_twe()  {
	handl := func(w http.ResponseWriter, r *http.Request) {

		s := r.URL.Path[1:]

		if strings.EqualFold(s, "red") == true{
			test_seven_url(w, redIndex)
			return
		}else if strings.EqualFold(s, "green") == true {
			test_seven_url(w, greenIndex)
			return
		}else if strings.EqualFold(s, "black") == true {
			test_seven_url(w, blackIndex)
			return
		}

	}
	http.HandleFunc("/", handl)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func main() {

	//Hello World
	test_one()

	//命令行参数
	//test_two()
	//test_thi()

	//查找重复的行
	//test_fou()
	//test_fiv()
	//test_six()

	//GIF 动画
	//test_seven("test1.gif", blackIndex)
	//test_seven("test2.gif", redIndex)
	//test_seven("test3.gif", greenIndex)

	//获取URL
	//test_eight("http://gopl.io")
	//test_eight("gopl.io")

	//并发获取多个URL
	//test_nine()

	//web服务
	//test_ten()
	//test_elev()
	//test_twe()
}
