package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

/* 函数 */
//1.函数式一组语句序列集合的单元
//2.函数科多次调用使用
//3.函数隐藏其实现的细节部分，至关重要的部分


//函数声明
//1.函数声明:函数名(形式参数列表）)(返回值列表(可忽略)){函数体}
//2.形式参数是局部变量
//3.形参是实参的拷贝
//4.形参一般不会影响实参，如果实参是指针、slice、map、func、channel类型，实参可能会被修改
//5.如果没有函数体的函数，则代表此函数不是Go实现的
//比如：func Sin(x float64) float //implemented in assembly language
func test_func()  {
	fmt.Println(hypot(3, 4))

	fmt.Printf("%T\n", add)   // "func(int, int) int"
	fmt.Printf("%T\n", sub)   // "func(int, int) int"
	fmt.Printf("%T\n", first) // "func(int, int) int"
	fmt.Printf("%T\n", zero)  // "func(int, int) int"
}
func hypot(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func add(x int, y int) int   {return x + y}
func sub(x, y int) (z int)   { z = x - y; return}
func first(x int, _ int) int { return x }
func zero(int, int) int      { return 0 }

//递归
//1.函数可以直接或者间接的调用自身，因此函数式可以递归的
//2.golang.org/x/... 目录下存储了一些由Go团队设计、维护，对网络编程、国际化文件处理、移动平台、图像处理、加密解密、开发者工具提供支持的扩展包
//3.未将扩展包加入标准库原因:1>部分包仍在开发中，2>扩展包提供的功能很少被使用
/*
	关于非标准包 golang.org/x 安装方法
	git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
	git clone https://github.com/golang/sys.git $GOPATH/src/golang.org/x/sys
	git clone https://github.com/golang/tools.git $GOPATH/src/golang.org/x/tools
*/
func test_recursive() {

	resp, err := http.Get("https://golang.org")
	if err != nil {
		fmt.Printf("fetch: %v\n", err)
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()

	if err != nil {
		fmt.Printf("findlinks: %v\n", err)
		return
	}

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
	fmt.Println("--------------------------分割线--------------------------")

	//练习 5.1： 修改findlinks代码中遍历n.FirstChild链表的部分，将循环调用visit，改成递归调用。
	for _, link := range visit_recursive(nil, doc) {
		fmt.Println(link)
	}
	fmt.Println("--------------------------分割线--------------------------")

	outline(nil, doc)
	fmt.Println("--------------------------分割线--------------------------")

	//练习 5.2： 编写函数，记录在HTML树中出现的同名元素的次数。
	dict := make(map[string]int)
	statistical_html(doc, &dict)
	for k, v := range dict {
		fmt.Printf("%8s---%d\n",k, v)
	}
	fmt.Println("--------------------------分割线--------------------------")

	//练习 5.3： 编写函数输出所有text结点的内容。
	//注意不要访问<script>和<style>元素,因为这些元素对浏览者是不可见的。
	text := text_all(nil, doc)
	fmt.Println(text)
	fmt.Println("--------------------------分割线--------------------------")


	//练习 5.4： 扩展visit函数，使其能够处理其他类型的结点，如images、scripts和style sheets。
	link := visit_extion(nil, doc)
	fmt.Println(link)
	fmt.Println("--------------------------分割线--------------------------")

}
//练习 5.4： 扩展visit函数，使其能够处理其他类型的结点，如images、scripts和style sheets。
func visit_extion(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		if n.Data == "a" ||
			n.Data == "img" ||
			n.Data == "link" ||
			n.Data == "scripts" ||
			n.Data == "sheets" ||
			n.Data == "style" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit_extion(links, c)
	}
	return links
}
//练习 5.3： 编写函数输出所有text结点的内容。
//注意不要访问<script>和<style>元素,因为这些元素对浏览者是不可见的。
func text_all(text []string, n *html.Node) []string {
	if n.Type == html.TextNode {
		text = append(text, n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data != "script" && c.Data != "style" {
			text = text_all(text, c)
		}
	}
	return text
}
func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		outline(stack, c)
	}
}
//练习 5.2： 编写函数，记录在HTML树中出现的同名元素的次数。
func statistical_html(n *html.Node, dict *map[string]int)  {

	if n.Type == html.ElementNode {
		(*dict)[n.Data] += 2//头尾各一个，因此是一次两个
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		statistical_html(c, dict)
	}
}
// visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
//练习 5.1： 修改findlinks代码中遍历n.FirstChild链表的部分，将循环调用visit，改成递归调用。
func visit_recursive(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	if n.FirstChild != nil {
		links = visit_recursive(links, n.FirstChild)
	}
	if n.NextSibling != nil {
		links = visit_recursive(links, n.NextSibling)
	}
	return links
}

//多返回值
//1.一个函数可以返回多个值
//2.调用多返回值函数时，必须显式的将值分配给变量，或者分配给_(blank identifier)
func test_more()  {
	links, err := findLinks("https://golang.org")
	if err != nil {
		fmt.Printf("findlinks: %v\n", err)
	}
	for _, link := range links {
		fmt.Println(link)
	}
	fmt.Println("--------------------------分割线--------------------------")

	//练习 5.5： 实现countWordsAndImages。（参考练习4.9如何分词）
	words, images, err := CountWordsAndImages("https://golang.org")
	fmt.Println(words, images, err)
}
// findLinks performs an HTTP GET request for url, parses the
// response as HTML, and extracts and returns the links.
func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}
//3.如果函数的返回值都是现实的变量名，则该函数的return可以省略操作数, 这称之为bare return
// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}
//练习 5.5： 实现countWordsAndImages。（参考练习4.9如何分词）
func countWordsAndImages(n *html.Node) (words, images int) {

	txts, images := visit_words(nil, n, 0)
	for _, word := range txts {
		wordss := strings.Fields(word)
		if len(wordss) > 0 {
			words += len(wordss)
		}
	}
	return
}

func visit_words(links []string, n *html.Node, img int) ([]string, int)  {
	if n.Type == html.TextNode {
		links = append(links, n.Data)
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		img++
	}

	if n.FirstChild != nil {
		if n.FirstChild.Data != "script" && n.FirstChild.Data != "style" {
			links, img = visit_words(links, n.FirstChild, img)
		}
	}
	if n.NextSibling != nil {
		if n.NextSibling.Data != "script" && n.NextSibling.Data != "style" {
			links, img = visit_words(links, n.NextSibling, img)
		}
	}
	return links, img
}


//错误
//1.错误是软件包API和应用程序用户界面的一个重要组成部分
//2.错误处理策略
//	1>>传播错误
//	2>>重新尝试失败的操作，但要限制时间和次数，防止无限制的重试
//	3>>输出错误信息并结束程序，这种策略只应在main中执行
//	4>>标准错误流输出错误信息
//	5>>可以直接忽略掉错误
//3.在Go中，错误处理有一套独特的编码风格
//4.文件结尾错误（EOF），io.EOF有固定的错误信息——“EOF”
func test_err() error  {

	url := "https://golan"

	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			fmt.Println("success")
			return nil
		}
		log.Printf("server not responding (%s);retrying…", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}

	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}


//函数值
//1.在Go中，函数被看作第一类值

func test_funcV()  {

	f := square
	fmt.Println(f(3)) // "9"

	f = negative
	fmt.Println(f(3))     // "-3"
	fmt.Printf("%T\n", f) // "func(int) int"

	//f = product//cannot use product (type func(int, int) int) as type func(int) int in assignment

	//2.函数类型的零值是nil，调用值为nil的函数值会引起panic错误，但是函数能和nil比较
	//var f1 func(int) int
	//f1(3)//panic: runtime error: invalid memory address or nil pointer dereference

	//3.函数值之间是不能比较的，因此也不鞥作为map的key
	//4.函数值不仅可以通过数据参数化函数，同样可以函数作为参数化数据
	fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
	fmt.Println(strings.Map(add1, "VMS"))      // "WNT"
	fmt.Println(strings.Map(add1, "Admix"))    // "Benjy"


	resp, err := http.Get("https://golang.org")
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	forEachNode(doc, startElement, endElement)

	//练习 5.8： 修改pre和post函数，使其返回布尔类型的返回值。
	//返回false时，中止forEachNoded的遍历。
	//使用修改后的代码编写ElementByID函数，
	//根据用户输入的id查找第一个拥有该id元素的HTML元素，
	//查找成功后，停止遍历。
	node := ElementByID(doc, "div")
	fmt.Printf("%#v\n",node)


	//练习 5.9： 编写函数expand，将s中的"foo"替换为f("foo")的返回值。
	fmt.Println(expand("footer", foo))

	//练习 5.7： 完善startElement和endElement函数，使其成为通用的HTML输出器。
	//要求：输出注释结点，文本结点以及每个元素的属性（< a href='...'>）。
	//使用简略格式输出没有孩子结点的元素（即用<img/>代替<img></img>）。
	//编写测试，验证程序输出的格式正确。（详见11章）
	forEachNode(doc, startElement_507, endElement_507)
}
func square(n int) int { return n * n }
func negative(n int) int { return -n }
func product(m, n int) int { return m * n }
func add1(r rune) rune { return r + 1 }

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
var depth int
func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		depth++
	}
}
func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}
//5.7
func startElement_507(n *html.Node) {
	if n.Type == html.ElementNode {
		s := ""
		for _, v := range n.Attr{
			s += "  " + v.Key + "=\"" + v.Val + "\"  "
		}
		fmt.Printf("%*s<%s%s", depth * 2, "", n.Data, s)
		depth++
	}
	if n.Type == html.ElementNode  {
		if n.FirstChild == nil && n.Data != "script" {
			fmt.Printf("/")
		}
		fmt.Printf(">\n")
	}
	if n.Type == html.TextNode {
		fmt.Printf("%*s %s\n", depth * 2, "", n.Data)
	}
}
//5.7
func endElement_507(n *html.Node) {
	depth--
	if n.Type == html.ElementNode {
		if n.FirstChild == nil && n.Data != "script" {
			fmt.Printf("\n")
		}
		fmt.Printf("%*s</%s>\n", depth * 2, "", n.Data)
	}
}
//5.8
func ElementByID(doc *html.Node, id string) *html.Node {

	if startElement_508(doc, id) == false{
		return doc
	}

	if doc.FirstChild != nil {
		ElementByID(doc.FirstChild, id)
	}
	if doc.NextSibling != nil {
		ElementByID(doc.NextSibling, id)
	}

	if endElement_508(doc, id) == false{
		return doc
	}
	return doc
}
//5.8
func startElement_508(n *html.Node, id string) (ok bool) {
	ok = true
	if n.Type == html.ElementNode && n.Data == id {
		ok = false
	}
	return
}
//5.8
func endElement_508(n *html.Node, id string) (ok bool) {
	ok = true
	if n.Type == html.ElementNode && n.Data == id {
		ok = false
	}
	return
}
//5.9
func expand(s string, f func(string) string) string {
	return strings.Replace(s, "foo", f("foo"), -1)
}
//5.9
func foo(s string) string  {
	return s + "+header+"
}

//匿名函数
//1.匿名函数字面量的语法和函数声明相似，区别在于func关键字后没有函数名
//2.匿名函数不仅仅是一串代码，还记录状态，因此属于引用类型，且函数值不可比较
//3.Go使用闭包（closures）技术实现函数值，Go程序员也把函数值叫做闭包
func test_anonymous()  {
	f := squares()
	fmt.Println(f()) // "1"
	fmt.Println(f()) // "4"
	fmt.Println(f()) // "9"
	fmt.Println(f()) // "16"

	// prereqs记录了每个课程的前置课程
	var prereqs = map[string][]string{
		"algorithms": {"data structures"},
		"calculus": {"linear algebra"},
		"linear algebra": {"calculus"},
		"compilers": {
			"data structures",
			"formal languages",
			"computer organization",
		},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
	}

	for i, course := range topoSort(prereqs){
		fmt.Printf("%d:\t%s\n", i + 1, course)
	}

	breadthFirst(crawl, []string{"https://golang.org"})

	//练习5.10： 重写topoSort函数，用map代替切片并移除对key的排序代码。
	//验证结果的正确性（结果不唯一）。
	for k, v := range topoSort_Map(prereqs) {
		fmt.Printf("%d:\t%s\n", k, v)
	}

	//练习5.11： 现在线性代数的老师把微积分设为了前置课程。
	//完善topSort，使其能检测有向图中的环。
	//能力有限，暂时不知道怎么搞，尝试却没得到结果

	//练习5.12： 
	//gopl.io/ch5/outline2（5.5节）的startElement和endElement共用了全局变量depth，
	//将它们修改为匿名函数，使其共享outline中的局部变量。
	outline2()
}
//练习5.13： 修改crawl，使其能保存发现的页面，
//必要时，可以创建目录来保存这些页面。
//只保存来自原始域名下的页面。
//假设初始页面在golang.org下，就不要保存vimeo.com下的页面。
func save_url(u string)  {

	resp, err := http.Get(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("此资源不存在\n")
		resp.Body.Close()
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Println("read data failed:", u, err)
		return
	}

	urll, _ := url.Parse(u)
	tmp := strings.TrimLeft(urll.Path, "/")
	filename := strings.ToLower(strings.Replace(tmp, "/", "-", -1))
	filename = "crawl/"+ filename
	if !strings.Contains(filename, "html") {
		filename += ".html"//确保保存下来的都是科打开的html文件
	}
	newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println("create file failed:", filename, err)
		return
	}
	newFile.Write(data)
	newFile.Close()
	data = nil
	return
}
func outline2()  {
	resp, err := http.Get("https://golang.org")
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}

	var depth int
	var startElementl func(n *html.Node)
	startElementl = func(n *html.Node) {
		if n.Type == html.ElementNode {
			s := ""
			for _, v := range n.Attr{
				s += "  " + v.Key + "=\"" + v.Val + "\"  "
			}
			fmt.Printf("%*s<%s%s", depth * 2, "", n.Data, s)
			depth++
		}
		if n.Type == html.ElementNode  {
			if n.FirstChild == nil && n.Data != "script" {
				fmt.Printf("/")
			}
			fmt.Printf(">\n")
		}
		if n.Type == html.TextNode {
			fmt.Printf("%*s %s\n", depth * 2, "", n.Data)
		}
	}
	var endElementl func(n *html.Node)
	endElementl = func(n *html.Node){
		depth--
		if n.Type == html.ElementNode {
			if n.FirstChild == nil && n.Data != "script" {
				fmt.Printf("\n")
			}
			fmt.Printf("%*s</%s>\n", depth * 2, "", n.Data)
		}
	}

	forEachNode(doc, startElementl, endElementl)
}
func topoSort_Map(m map[string][]string) map[int]string  {
	seen := make(map[string]bool)

	order := make(map[int]string)
	index := 1
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, key := range items{
			if !seen[key] {
				seen[key] = true
				visitAll(m[key])
				order[index] = key
				index++
			}
		}
	}
	keys := []string{}
	for key := range m{
		keys = append(keys, key)
	}
	visitAll(keys)

	return order
}
func crawl(url string) []string {

	//练习5.13
	go save_url(url)
	fmt.Println(url)
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
// 广度优先遍历
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}
// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}
// 深度优先遍历
func topoSort(m map[string][]string) (ss []string)  {
	seen := make(map[string]bool)

	//4.匿名函数需要递归时，必须提前声明，否则无法递归绑定
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items{
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				ss = append(ss, item)
			}
		}
	}

	var keys []string
	for key := range m{
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return
}
// squares返回一个匿名函数。
// 该匿名函数每次被调用时都会返回下一个数的平方。
func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}
func main() {

	//函数声明
	//test_func()

	//递归
	//test_recursive()

	//多返回值
	//test_more()

	//错误
	//if err := test_err(); err != nil {
	//	fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
	//	os.Exit(1)
	//}

	//函数值
	//test_funcV()

	//匿名函数
	test_anonymous()
}
