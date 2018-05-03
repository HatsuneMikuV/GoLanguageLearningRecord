package main

import (
	"20180408/byteCounter"
	"20180408/tempconv"
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

/*
	第七章 接口

	1.接口类型是对其它类型行为的抽象和概括
	2.Go语言中接口类型的独特之处在于它是满足隐式实现的
*/


//一，接口是合约
//1.接口合约就像ObjC的继承，继承父类的方法并重新实现，因此可以被合约原函数调用
func test_convention()  {

	fmt.Println("test")

	//因为*ByteCounter满足io.Writer的约定，我们可以把它传入Fprintf函数中
	//
	var c byteCounter.ByteCounter
	c.Write([]byte("hello"))
	fmt.Println(c) // "5", = len("hello")
	c = 0          // reset the counter
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c) // "12", = len("hello, Dolly")

	fmt.Println("================================")

	//练习 7.1： 使用来自ByteCounter的思路，实现一个针对对单词和行数的计数器。
	ss := []byte("hello\nI like you\nI like you too")

	var w byteCounter.WordCounter
	w.Write(ss)
	fmt.Println(w)
	w = 0
	fmt.Fprintf(&w, "hi, %s", string(ss))
	fmt.Println(w)
	fmt.Println("================================")

	var l byteCounter.LineCounter
	l.Write(ss)
	fmt.Println(l)
	l = 0
	fmt.Fprintf(&l, "hi\n %s", string(ss))
	fmt.Println(l)
	fmt.Println("================================")

	//练习 7.2： 写一个带有如下函数签名的函数CountingWriter，
	//传入一个io.Writer接口类型，
	//返回一个新的Writer类型把原来的Writer封装在里面和一个表示写入新的Writer字节数的int64类型指针
	var b byteCounter.ByteCounter
	b.Write([]byte("hello"))
	cc, nbytes := byteCounter.CountingWriter(&b)
	fmt.Println(cc, nbytes)
	fmt.Println("================================")

	//练习 7.3： 为在gopl.io/ch4/treesort (§4.4)的*tree类型实现一个String方法去展示tree类型的值序列。
	arr := []int{9, 8, 3, 4, 5, 6, 7}
	tree1 := byteCounter.Sort(arr)
	fmt.Println(tree1)
}

//二，接口类型
//1.接口类型具体描述了一系列方法的集合, 而实现这个方法的具体类型是这个接口类型的实例
//2.io.Writer类型是用的最广泛的接口之一，因为它提供了所有的类型写入bytes的抽象
//3.接口类型可以组合定义，成为一个集合方法
func test_interface_type()  {

	//练习7.4
	node,_ := NewReader("<html>111111</html>")
	fmt.Println(node)
	fmt.Println("================================")

	//练习 7.5
	ss := []byte("11112222")
	rr := bytes.NewReader(ss)
	reader := LimitReader(rr, 6)
	s, _ := ioutil.ReadAll(reader)
	fmt.Println(string(s))

}
type HtmlReader struct {
	r io.Reader
}
func (reader *HtmlReader) Read(p []byte) (n int, err error) {
	n, err = reader.r.Read(p)
	return
}
func creatReader(r io.Reader) io.Reader {
	return &HtmlReader{r:r}
}
func NewReader(s string) (*html.Node, error)  {
	hr := creatReader(strings.NewReader(s))
	node, err := html.Parse(hr)
	return node, err
}
/*
   练习 7.5： io包里面的LimitReader函数接收一个io.Reader接口类型的r和字节数n，
   并且返回另一个从r中读取字节,但是当读完n个字节后就表示读到文件结束的Reader。
   实现这个LimitReader函数.
*/
type MyLimitReader struct {
	R io.Reader
	N int64
}

func (myLimitReader *MyLimitReader) Read(p []byte) (n int, err error) {
	if myLimitReader.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p))  > myLimitReader.N {
		p = p[0:myLimitReader.N]
	}
	n, err = myLimitReader.R.Read(p)
	myLimitReader.N -= int64(n)
	return
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &MyLimitReader{r, n}
}


//三，实现接口的条件
//1.表达一个类型属于某个接口只要这个类型实现这个接口
//2.即使具体类型有其它的方法也只有接口类型暴露出来的方法会被调用到
//3.因为接口类型被称为空接口类型，因此可以将任意值赋给接口类型
func test_interface_condition()  {

	os.Stdout.Write([]byte("hello")) // OK: *os.File has Write method
	//os.Stdout.Close()                // OK: *os.File has Close method
	fmt.Println("\n================================")

	var w io.Writer
	w = os.Stdout
	w.Write([]byte("hello")) // OK: io.Writer has Write method
	//w.Close()//w.Close undefined (type io.Writer has no field or method Close)

	fmt.Println("\n================================")

	var any interface{}
	any = true
	any = 12.34
	any = "hello"
	any = map[string]int{"one": 1}
	any = new(bytes.Buffer)
	fmt.Println(any)
	fmt.Println("================================")

	type Artifact interface {
		Title() string
		Creators() []string
		Created() time.Time
	}
	type Text interface {
		Pages() int
		Words() int
		PageSize() int
	}
	type Audio interface {
		Stream() (io.ReadCloser, error)
		RunningTime() time.Duration
		Format() string // e.g., "MP3", "WAV"
	}
	type Video interface {
		Stream() (io.ReadCloser, error)
		RunningTime() time.Duration
		Format() string // e.g., "MP4", "WMV"
		Resolution() (x, y int)
	}

	type Streamer interface {
		Stream() (io.ReadCloser, error)
		RunningTime() time.Duration
		Format() string
	}
}
//四，flag.Value接口
func test_flag()  {
	//var period = flag.Duration("period", 1*time.Second, "sleep period")
	var period = flag.Duration("period", 50 * time.Millisecond, "sleep period")
	//var period = flag.Duration("period", 2 * time.Minute + 30 *time.Second, "sleep period")
	//var period = flag.Duration("period", 1 * time.Hour + 30 * time.Minute, "sleep period")
	//var period = flag.Duration("period", 24 * time.Hour, "sleep period")

	flag.Parse()
	fmt.Printf("Sleeping for %v...\n", *period)
	time.Sleep(*period)

	type Value interface {
		String() string
		Set(string) error
	}
	fmt.Println("================================")

	//var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")
	//var temp = tempconv.CelsiusFlag("temp", -18.0, "the temperature")
	//flag.Parse()
	//fmt.Println(*temp)

	//练习 7.6： 对tempFlag加入支持开尔文温度。
	var tempF = tempconv.FahrenheitFlag("temp", -18.0, "the temperature")

	flag.Parse()
	fmt.Println(*tempF)

	//练习 7.7： 解释为什么帮助信息在它的默认值是20.0没有包含°C的情况下输出了°C。
	//因为CelsiusFlag实现了set接口，一个*Celsius类型赋给了flag，flag实现的stringter接口
	//最终使Celsius调用了自身实现的string方法，从而将Celsius的值转成带°C的字符串
}
//五，接口值
//1.有两部分组成:一个具体的类型，一个此类型的值
//2.也被称为动态类型和动态值
//3.一个接口值可以持有任意大的动态值
//4.一个接口上的调用必须使用动态分配
//5.接口值得动态类型如果是可以比较的，即可以作为map的key或者switch的语句操作数
//6.一个不包含任何值的nil接口值和一个刚好包含nil指针的接口值是不同的
func test_interface_value()  {


	var w io.Writer
	fmt.Println(w,"---0")
	w = os.Stdout
	fmt.Println(w,"---1")
	w = new(bytes.Buffer)
	fmt.Println(w,"---2")
	w = nil
	fmt.Println(w,"---3")

	const debug = false

	var buf io.Writer//*bytes.Buffer
	if debug {
		buf = new(bytes.Buffer) // enable collection of output
	}
	f_nil(buf) // NOTE: subtly incorrect!
	if debug {
		// ...use buf...
	}
	//f_nil(buf):违反了(*bytes.Buffer).Write方法的接收者非空的隐含先觉条件
	//debug -- false : panic: runtime error: invalid memory address or nil pointer dereference
	fmt.Println(buf)
	//7.对于一些如*os.File的类型，nil是一个有效的接收者，但是*bytes.Buffer类型不在这些类型中
}
//如果输出的是非nil，输出江北写入Write函数
func f_nil(out io.Writer) {
	// ...do something...
	if out != nil {
		out.Write([]byte("done!\n"))
	}
}

//六，sort.Interface接口
//1.sort包内置的提供了根据一些排序函数来对任何序列排序的功能
//2.Go语言的sort.Sort函数不会对具体的序列和它的元素做任何假设
//3.Go使用了一个接口类型sort.Interface来指定通用的排序算法和可能被排序到的序列类型之间的约定
//4.一个内置的排序算法需要三个东西：序列的长度，表示两个元素比较的结果，一种交换两个元素的方式
func test_sort_Interface()  {

	names := []string{"1", "2", "5", "3", "4"}
	sort.Sort(StringSlice(names))
	fmt.Println(names)

	//sort.Sort(byArtist(tracks))
	//printTracks(tracks)

	fmt.Println("111===================================================================")
	//sort.Sort(sort.Reverse(byArtist(tracks)))
	//printTracks(tracks)

	fmt.Println("222===================================================================")
	//sort.Sort(byYear(tracks))
	//printTracks(tracks)

	fmt.Println("333===================================================================")
	//sort.Sort(customSort{tracks, func(x, y *Track) bool {
	//	if x.Title != y.Title {
	//		return x.Title < y.Title
	//	}
	//	if x.Year != y.Year {
	//		return x.Year < y.Year
	//	}
	//	if x.Length != y.Length {
	//		return x.Length < y.Length
	//	}
	//	return false
	//}})
	//printTracks(tracks)

	fmt.Println("444===================================================================")
	values := []int{3, 1, 4, 1}
	fmt.Println(sort.IntsAreSorted(values)) // "false"

	sort.Ints(values)
	fmt.Println(values)                     // "[1 1 3 4]"
	fmt.Println(sort.IntsAreSorted(values)) // "true"

	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	fmt.Println(values)                     // "[4 3 1 1]"
	fmt.Println(sort.IntsAreSorted(values)) // "false"

	fmt.Println("555===================================================================")
	//练习 7.8
	// 模拟点击记录
	var trackss = []*Track{
		{"Go", "Delilah", "Form Roots up", 2012, length("3m37s")},
		{"Go", "Bob", "Form Roots down", 2012, length("3m37s")},
		{"Ready 2 Go", "Moby", "Moby", 1992, length("3m37s")},
		{"Go", "Bob", "As I Am", 2012, length("3m37s")},
		{"Go", "Bob", "As I Am", 2012, length("3m20s")},
		{"Go", "Martin Solveing", "Smash", 2011, length("3m37s")},
	}
	clickRecords = append(clickRecords, "Title")
	clickRecords = append(clickRecords, "Year")
	clickRecords = append(clickRecords, "Artist")
	clickRecords = append(clickRecords, "Album")
	clickRecords = append(clickRecords, "Length")
	printTracks(trackss)

	fmt.Println("666===================================================================")
	sort.Sort(multiSort{trackss, less})
	printTracks(trackss)

	//练习 7.9
	//handl := func(w http.ResponseWriter, r *http.Request) {
	//	keyWords := r.URL.Path[1:]
	//	fmt.Println(keyWords)
	//
	//	if len(keyWords) <= 0 {
	//		return
	//	}
	//	searchKey := []string{keyWords}
	//
	//	result, err := github.SearchIssues(searchKey)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	if err := issueList.Execute(w, result); err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//http.HandleFunc("/", handl)
	//log.Println(http.ListenAndServe("localhost:8080", nil))

	http.HandleFunc("/", home)
	http.HandleFunc("/post", post)
	log.Println(http.ListenAndServe("localhost:1234", nil))
}
type User struct {
	Name string
	Date time.Time
}

var users []string

func home(w http.ResponseWriter, _ *http.Request)  {
	if err := homePage.Execute(w, users); err != nil {
		log.Printf("%v", err)
	}
}

func post(w http.ResponseWriter, r *http.Request)  {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
		return
	}
	fmt.Println(r.PostFormValue("but"))
	fmt.Println(r.PostFormValue("Title"))
	fmt.Println(r.PostFormValue("name"))
}
var homePage = template.Must(template.New("home").Parse(
	`<html><body>
<form action="/post"method="post"><br/>
<input type='button'value='Title'id='but1'/>
<input type='button'value='Artist'id='but2'/>
<input type='button'value='Album'id='but3'/>
<input type='button'value='Year'id='but4'/>
<input type='button'value='Length'id='but5'/>
</form></body></html>
`))

type StringSlice []string
func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table
}
type byArtist []*Track
func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byYear []*Track
func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

/*
	练习 7.8： 很多图形界面提供了一个有状态的多重排序表格插件：
		主要的排序键是最近一次点击过列头的列，
		第二个排序键是第二最近点击过列头的列，等等。
		定义一个sort.Interface的实现用在这样的表格中。
		比较这个实现方式和重复使用sort.Stable来排序的方式。
*/
type multiSortTrack []*Track

type multiSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x multiSort) Len() int {
	return len(x.t)
}

func (x multiSort) Less(i, j int) bool {
	return x.less(x.t[i], x.t[j])
}

func (x multiSort) Swap(i, j int) {
	x.t[i], x.t[j] = x.t[j], x.t[i]
}

var clickRecords []string

func less(x, y *Track) bool {
	for _, click := range clickRecords {
		if click == "Title" {
			if x.Title == y.Title {
				continue
			}
			return x.Title < y.Title
		}
		if click == "Year" {
			if x.Year == y.Year {
				continue
			}
			return x.Year < y.Year
		}
		if click == "Artist" {
			if x.Artist == y.Artist {
				continue
			}
			return x.Artist < y.Artist
		}
		if click == "Album" {
			if x.Album == y.Album {
				continue
			}
			return x.Album < y.Album
		}
		if click == "Length" {
			if x.Length == y.Length {
				continue
			}
			return x.Length < y.Length
		}
	}
	return false
}

//七，http.Handler接口
//八，error接口
//九，示例: 表达式求值
//十，类型断言
//十一，基于类型断言识别错误类型
//十二，通过类型断言查询接口
//十三，类型分支
//十四，示例: 基于标记的XML解码
//十五，补充几点


func main() {

	//一，接口是合约
	//test_convention()

	//二，接口类型
	//test_interface_type()

	//三，实现接口的条件
	//test_interface_condition()

	//四，flag.Value接口
	//test_flag()

	//五，接口值
	//test_interface_value()

	//六，sort.Interface接口
	test_sort_Interface()
	//七，http.Handler接口

	//八，error接口

	//九，示例: 表达式求值

	//十，类型断言

	//十一，基于类型断言识别错误类型

	//十二，通过类型断言查询接口

	//十三，类型分支

	//十四，示例: 基于标记的XML解码

	//十五，补充几点

}

