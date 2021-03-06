package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"
)

/* 函数 */
//1.函数式一组语句序列集合的单元
//2.函数科多次调用使用
//3.函数隐藏其实现的细节部分，至关重要的部分


//一，函数声明
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

//二，递归
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

//三，多返回值
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


//四，错误
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


//五，函数值
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

//六，匿名函数
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
//捕获迭代变量------Go词法作用域的一个陷阱
//循环变量中被声明的变量会被每次循环共享变量，为了解决这个问题，
//在循环的局部创建一个新的同类型变量，作为共享变量的副本，这杨就能避免共享变量被重复使用
//如下面的变量dir，虽然这看起来很奇怪，但却很有用
//	for _, dir := range tempDirs() {
//		dir := dir // declares inner dir, initialized to outer dir
//		...
//	}
//go或者defer  同样会导致这个问题的出现，解决的办法也是重新创建一个变量作为副本


//七，可变参数
//1.参数数量可变的函数被称为可变参数函数，例如fmt.Printf等类似的函数，参数不受限制
//2.声明可变参数时，需要在参数类型之前加上省略号…

func test_Variable()  {

	fmt.Println(sum())           // "0"
	fmt.Println(sum(3))          // "3"
	fmt.Println(sum(1, 2, 3, 4)) // "10"

	values := []int{1, 2, 3, 4}
	fmt.Println(sum(values...)) // "10"

	//可变参数函数和以切片作为参数的函数是不同的
	fmt.Printf("%T\n", fff) // "func(...int)"
	fmt.Printf("%T\n", ggg) // "func([]int)

	linenum, name := 12, "count"
	errorf(linenum, "undefined: %s", name)

	fmt.Println(sum_max())
	fmt.Println(sum_min(4, 3, 5))
	fmt.Println(sum_Ex2(4, 1, 2))

	fmt.Println(strings.Join([]string{"1", "2", "3"}, "4"))
	fmt.Println(strings_join([]string{"1", "2", "3"}, "4", "5", "6"))


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
	images := ElementsByTagName(doc, "img")
	headings := ElementsByTagName(doc, "div", "h2", "h3", "h4")

	fmt.Println(len(images))
	fmt.Println(len(headings))
}
//练习5.17：编写多参数版本的ElementsByTagName，
//函数接收一个HTML结点树以及任意数量的标签名，
//返回与这些标签名匹配的所有元素
func ElementsByTagName(doc *html.Node, name...string) []*html.Node {
	nodes := []*html.Node{}
	for _, nam := range name{
		visitNode := func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == nam {
				nodes = append(nodes, n)
			}
		}
		forEachNode(doc, visitNode, nil)
	}
	return nodes
}

//练习5.16：编写多参数版本的strings.Join。
func strings_join(a []string, sep string, args ...interface{}) string {

	ss := ""

	ss += fmt.Sprint(args...)

	switch len(a) {
	case 0:
		return ss
	case 1:
		return a[0] + ss
	case 2:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1] + ss
	case 3:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1] + sep + a[2] + ss
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b) + ss
}

//练习5.15： 编写类似sum的可变参数函数max和min。
//考虑不传参时，max和min该如何处理，
//再编写至少接收1个参数的版本。

//不传参
func sum_max(vals...int) int {
	max := math.MinInt64
	if len(vals) <= 0 {
		return 0
	}
	for _, val := range vals {
		if max < val {
			max = val
		}
	}
	return max
}
func sum_min(vals...int) int {
	min := math.MaxInt64
	if len(vals) <= 0 {
		return 0
	}
	for _, val := range vals {
		if min > val {
			min = val
		}
	}
	return min
}

//至少接收1个参数
func sum_Ex2(val int, vals...int) int {
	total := val
	for _, val := range vals {
		total += val
	}
	return total
}
//4.interfac{}表示函数的最后一个参数可以接收任意类型
func errorf(linenum int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}
func fff(...int) {}
func ggg([]int) {}
//3.隐式的创建一个数组，然后将参数复制到数组中，作为一个切片传给函数,
//如果参数本身就是切片，则需要在后面加上省略号
func sum(vals...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}


//八，Defer 函数
//1.defer语句被执行时，跟在defer后面的函数会被延迟执行
//2.直到包含该defer语句的函数执行完毕时，defer后的函数才会被执行
//3.函数中可执行多条defer语句，它们的执行顺序与声明顺序相反
//4.defer语句经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁
func test_defer()  {
	//fmt.Println(title("http://gopl.io"))
	//fmt.Println(title("https://golang.org/doc/effective_go.html"))
	//fmt.Println(title("https://golang.org/doc/gopher/frontpage.png"))

	//bigSlowOperation()

	//_ = double(4)

	fmt.Println(triple(4))

	fmt.Println(fetch_defer("https://golang.org"))
}
//练习5.18：不修改fetch的行为，重写fetch函数，要求使用defer机制关闭文件。
/*
	通过os.Create打开文件进行写入，在关闭文件时，我们没有对f.close采用defer机制，因为这会产生一些微妙的错误。
	许多文件系统，尤其是NFS，写入文件时发生的错误会被延迟到文件关闭时反馈。
	如果没有检查文件关闭时的反馈信息，可能会导致数据丢失，而我们还误以为写入操作成功。
	如果io.Copy和f.close都失败了，我们倾向于将io.Copy的错误信息反馈给调用者，
	因为它先于f.close发生，更有可能接近问题的本质
*/
func fetch_defer(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	//在下面return之后，会关闭网络请求连接
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

	defer f.Close()//在return之后关闭文件
	n, err = io.Copy(f, resp.Body)
	//在关闭文件之前，将copy产生的错误信息反馈出去
	return local, n, err
}
//7.被延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值
func triple(x int) (result int) {
	defer func() { result += x }()
	return double(x)
}
//6.因为defer的性质，会在return之后再执行，因此可以用来观测匿名函数的返回值
func double(x int) (result int) {
	defer func() { fmt.Printf("double(%d) = %d\n", x,result) }()
	return x + x
}
func bigSlowOperation() {
	//5.不要忘记defer语句后的圆括号，
	//否则本该在进入时执行的操作会在退出时执行，
	//而本该在退出时执行的，永远不会被执行
	defer trace("bigSlowOperation")() // don't forget the
	//extra parentheses
	// ...lots of work…
	time.Sleep(10 * time.Second) // simulate slow
	//operation by sleeping
}
func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() {
		log.Printf("exit %s (%s)", msg,time.Since(start))
	}
}
func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check Content-Type is HTML (e.g., "text/html;charset=utf-8").
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct,"text/html;") {
		return fmt.Errorf("%s has type %s, not text/html",url, ct)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url,err)
	}
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title"&&n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	forEachNode(doc, visitNode, nil)
	return nil
}

//九，Panic异常
//1.运行时的错误会引起Panic异常
//2.panic异常发生时，程序会中断运行，并立即执行在该goroutine中被延迟的函数（defer 机制)
//3.panic一般用于严重错误，在健壮的程序中，任何可以预料到的错误，最好的处理方式是使用Go的错误机制
func test_panic()  {

	s := ""
	switch s {
	case "Spades":                                // ...
	case "Hearts":                                // ...
	case "Diamonds":                              // ...
	case "Clubs":                                 // ...
	default:
		//panic(fmt.Sprintf("invalid s %q", s)) // Joker?
	}
	//Reset(nil)

	//在Go的panic机制中，延迟函数的调用在释放堆栈信息之前。
	defer printStack()
	f_panic(3)
}
func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}
func f_panic(x int) {
	fmt.Printf("f(%d)\n", x+0/x) // panics if x == 0
	defer fmt.Printf("defer %d\n", x)
	f_panic(x - 1)
}
//断言函数必须满足的前置条件是明智的做法，但这很容易被滥用
//除非你能提供更多的错误信息，或者能更快速的发现错误，
//否则不需要使用断言，编译器在运行时会帮你检查代码
func Reset(x *string) {
	if x == nil {
		panic("x is nil") // unnecessary!
	}
	x = nil
}

//十，Recover捕获异常
//1.panic异常，在一些情况下需要恢复正常，这时需要用到Recover，前提是defer调用它，并且定义该defer语句的函数发生了panic异常
//2.recover能使异常的长须回复正常，并且能返回oanic的value，如果并没有panic调用recover，则返回的是nil
//3.不应该恢复一个由他人开发的函数引起的panic
//4.公有的API应该将函数的运行失败作为error返回
//5.并不是所有的异常都能被恢复，比如内存不足时，go会终止程序的运行，是不能被恢复的
func test_Recove()  {

	//使用study_01的test_httpUrl来启动一个web服务器模拟这个情况的复现
	//resp, err := http.Get("http://localhost:8080/")
	//if err != nil {
	//	return
	//}
	//defer resp.Body.Close()
	//doc, err := html.Parse(resp.Body)
	//if err != nil {
	//	err = fmt.Errorf("parsing HTML: %s", err)
	//	return
	//}
	//title, err :=  soleTitle(doc)
	//if err != nil {
	//	fmt.Printf("soleTitle:%s", err)
	//	return
	//}
	//fmt.Println(title)

	x := string("")
	fmt.Println(test_panic_recover(&x))
}
//练习5.19： 使用panic和recover编写一个不包含return语句但能返回一个非零值的函数。
func test_panic_recover(xx *string) (err string)  {

	defer func() {
		p := recover()
		if p == nil {
			// no panic
		}else if p == "333" {
			err = fmt.Sprintf("panic recover : %s", p)
		}
	}()

	*xx  = "333"
	panic(*xx)
}
// soleTitle returns the text of the first non-empty title element
// in doc, and an error if there was not exactly one.
func soleTitle(doc *html.Node) (title string, err error) {
	type bailout struct{}
	defer func() {
		switch p := recover(); p {
		case nil:       // no panic
		case bailout{}: // "expected" panic
			err = fmt.Errorf("multiple title elements")
		default:
			panic(p) // unexpected panic; carry on panicking
		}
	}()

	// Bail out of recursion if we find more than one nonempty title.
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			if title != "" {
				panic(bailout{}) // multiple titleelements
			}
			title = n.FirstChild.Data
		}
	}, nil)
	if title == "" {
		return "", fmt.Errorf("no title element")
	}
	return title, nil
}

func main() {

	//一，函数声明
	//test_func()

	//二，递归
	//test_recursive()

	//三，多返回值
	//test_more()

	//四，错误
	//if err := test_err(); err != nil {
	//	fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
	//	os.Exit(1)
	//}

	//五，函数值
	//test_funcV()

	//六，匿名函数
	//test_anonymous()

	//七，可变参数
	//test_Variable()

	//八，Defer 函数
	//test_defer()

	//九，Panic异常
	//test_panic()

	//十，Recover捕获异常
	test_Recove()
}
