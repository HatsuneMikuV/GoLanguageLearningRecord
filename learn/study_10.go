package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

/* 标准包 */


//格式包

func test_one()  {
	var i = 42
	var s = "Answer"
	fmt.Print(s, "is", i, 3.14, '\n', "\n")

	fmt.Println(s, "is", i, 3.14, '\n', "\n")

	fmt.Printf("%s is %d\n %f %c %x\n", s, i, 3.14, '\n', '\n')

	fmt.Printf("Answer is %06d and %0.4f\n", 42, 3.14159)

	fmt.Printf("Say %6.4s!\n", "hello")

	fmt.Printf("Say %*.*s!\n", 4, 2, "hello")

	fmt.Printf("% d\n% d\n% x\n", -42, 42, "你好")
}

//字节包
func test_two()  {
	s := "你Go了吗？"
	fmt.Println(bytes.Count([]byte(s), nil))
	fmt.Println(strings.Count(s, ""))
}

//模板包
const tpl  = `
How many roads must {{.}} walk down
Before they call him {{.}}
`
func test_thi()  {
	tmpl := template.New("M2")
	tmpl.Parse(tpl)
	tmpl.Execute(os.Stdout, "a man")
}

func test_fou() {
	tpll := `
	知止{{range.}}而后能{{.}}，{{.}}{{end}}而后能得
	`

	var 大学= []string{"定", "静", "安", "虑"}

	tmpll := template.New("")
	tmpll.Parse(tpll)
	tmpll.Execute(os.Stdout, 大学)
	fmt.Println("\n")
}


//正则表达式包
var pat = `\s*[a-hj-y][a-hj-z]`

func test_fiv()  {
	var text = "qwrqb~@#~@$~@$feje.fjas\\gdfsgbns\\sdgs kdg\\klbfasbrganskfabkbn/sdns;r.masgknsf"

	reg, _ := regexp.Compile(pat)

	for i, v := range reg.FindAllString(text, -1){
		fmt.Errorf("%d : %s\n", i, v)
	}
}

//时间包
func test_six()  {
	now := time.Now()
	fmt.Println(now)
	fmt.Println(now.Format(time.RFC3339))
	fmt.Println(now.Format("01月02日03时04分05秒06年07区"))

	now = time.Now()

	yr, mo, _ := now.Date()
	d1 := time.Date(yr, mo, 1, 0,0,0,0,time.UTC)
	d2 := time.Date(yr, mo+1, 1, 0,0,0,0,time.UTC)
	d2 = d2.Add(-24 * time.Hour)

	w := d1.Weekday()
	_, ww := d1.ISOWeek()
	fmt.Println("\n周\n 日 一 二 三 四 五 六 ")
	fmt.Println("%d\n%*d", ww, int((w+1) * 3), 1)
	for i := 2; i <= d2.Day(); i++ {
		if w++; w%7 == 0 {
			ww++
			fmt.Println("\n%d\n", ww)
		}
		fmt.Printf("%3d", i)
	}

	fmt.Print("\n\n")
	var days = []rune("日一二三四五六")

	now1 := time.Now()
	w2 := now1.Weekday()
	fmt.Print("星期", string(days[w2]), "\n")
}


//超链接包
func test_seven()  {

	//http.HandleFunc("/", func(w  http.ResponseWriter, _ *http.Request) {
	//	io.WriteString(w, "Hi there!")
	//})
	//http.ListenAndServe(":1234", nil)

	http.HandleFunc("/", hi)
	http.ListenAndServe(":1234", nil)
}
func hi(w  http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, r.URL.Path[1:] + "你好")

	var s = r.URL.RawQuery
	if len(s) > 0 {
		fmt.Fprintf(w, "\nRawQuery: %s\n", r.URL.RawQuery)
		fmt.Fprintf(w, "Name: %s\n", r.URL.Query()["Name"])
	}
}

func main() {

	//格式包
	test_one()

	//字节包
	test_two()

	//模板包
	test_thi()
	test_fou()

	//正则表达式包
	test_fiv()

	//时间包
	test_six()

	//超链接包
	test_seven()
}
