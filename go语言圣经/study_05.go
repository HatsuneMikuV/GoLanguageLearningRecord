package main

import (
	"fmt"
	"golang.org/x/net/html"
	"math"
	"net/http"
	"os"
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
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
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
func main() {

	//函数声明
	test_func()

	//递归
	test_recursive()
}
