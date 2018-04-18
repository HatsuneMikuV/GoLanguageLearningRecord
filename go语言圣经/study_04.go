package main

import (
	"20180408/github"
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

/* 复合数据类型 */

//数组
//1.数组是另零个或者多个特定类型的元素组成的固定长度的序列
//2.数组的长度是固定的，因此Go语言很少直接使用数组，而使用可变数组Slice（切片）
func test_array()  {

	//3.数组的每个元素可以通过下标来访问，从0到len - 1
	//4.默认情况下，数组元素初始化为元素类型的零值，比如数字元素为0
	var a [3]int
	fmt.Println(a[0])
	fmt.Println(a[len(a) - 1])

	for i, v := range a {
		fmt.Printf("%d %d %T-%T\n", i, v, i, v)
	}

	for _, v := range a {
		fmt.Printf("%d %T\n", v, v)
	}

	var q [3]int = [3]int{1, 2, 3}
	var r [3]int = [3]int{1, 2}
	fmt.Println(q[2])
	fmt.Println(r[2])

	//5.数组初始化时如果未给长度，会根据元素个数来计算
	qq := [...]int{1, 2, 3, 4}
	fmt.Printf("%T\n", qq)

	//6.数组的长度在编译阶段就确定了
	q1 := [3]int{1, 2, 3}
	//q1 = [4]int{1, 2, 3, 4}//cannot use [4]int literal (type [4]int) as type [3]int in assignment
	fmt.Println(q1)

	type Currency  int
	const (
		USD Currency = iota //美元
		EUR 				//欧元
		GBP 				//英镑
		RMB 				//人民币
	)

	symbol := [...]string { USD:"$", EUR:"€", GBP:"£", RMB:"¥"}
	fmt.Println(RMB, symbol[RMB])

	r1 := [...]int{99: -1}
	fmt.Println(r1)

	//7.数组长度是类型的一部分
	a1 := [2]int{1, 2}
	b1 := [...]int{1, 2}
	c1 := [2]int{1, 3}
	fmt.Println(a1 == b1, a1 == c1, b1 == c1)
	//d1 := [3]int{1, 2}
	//fmt.Println(a1 == d1) //invalid operation: a1 == d1 (mismatched types [2]int and [3]int)

	c2 := sha256.Sum256([]byte("x"))
	c3 := sha256.Sum256([]byte("X"))
	fmt.Printf("%x\n%x\n%t\n%T\n", c2, c3, c2 == c3, c2)

	//练习 4.1： 编写一个函数，计算两个SHA256哈希码中不同bit的数目
	aaa := sha256.Sum256([]byte("123"))
	bbb := sha256.Sum256([]byte("124"))
	fmt.Printf("%x\n%x\n", aaa, bbb)
	fmt.Println(compareSHA256(aaa, bbb))

	//练习 4.2： 编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法
	sss := "test"
	fmt.Println(shaPrinf("sha512", sss) + "\n")
	fmt.Println(shaPrinf("sha384", sss) + "\n")
	fmt.Println(shaPrinf("sha256", sss) + "\n")
	fmt.Println(shaPrinf("", sss))

}
func compareSHA256(a [sha256.Size]byte, b [sha256.Size]byte) int {
	num := 0

	for i := 0; i < len(a); i++ {
		for j := 1; j <= 8; j++ {
			aa := a[i] >> uint8(j)
			bb := b[i] >> uint8(j)
			if aa != bb {
				num++
			}
		}
	}
	return num
}

func shaPrinf(flag string, s string) string  {

	if flag == "sha384" {
		return fmt.Sprintf("%x", sha512.Sum384([]byte(s)))
	}else if flag == "sha512" {
		return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

//Slice
//1.Slice是一组相同类型的可变序列，可理解为可变的数组
//2.Slice是一个轻量级的数据结构，可访问底层引用对象的全部元素
//3.Slice由指针、长度(len)、容量(cap)构成，长度和容量不一定相等
func test_slice()  {
	moths := [...]string{
		1:"Jannary",
		2:"February",
		3:"March",
		4:"April",
		5:"May",
		6:"June",
		7:"July",
		8:"August",
		9:"September",
		10:"October",
		11:"November",
		12:"December",
	}

	Q2 := moths[4:7]
	summer := moths[6:9]

	fmt.Println(Q2)
	fmt.Println(summer)

	//4.使用冒泡算法求相同月份(性能较低)
	for _, sm := range summer{
		for _, qm := range Q2{
			if sm == qm {
				fmt.Printf("%s appears in both\n", sm)
			}
		}
	}

	//5.slice 操作长度(len)和容量(cap)时，cap不能超出范围
	//  len超出意味着基于底层数组扩展，同样不能超出cap
	//fmt.Println(summer[:20]) // panic: runtime error: slice bounds out of range

	endlessSummer := summer[:5]
	fmt.Println(endlessSummer)

	a := []int{0, 1, 2, 3, 4, 5, 6}
	reverse(a)
	fmt.Println(a)

	s := []int{0, 1, 2, 3, 4, 5}

	//6.slice 不能使用==的操作判断，因为slice的元素是间接引用的，唯一能操作的是==nil
	reverse(s[:2])
	reverse(s[2:])
	reverse(s)
	fmt.Println(s)

	var runes []rune
	for _, r := range "Hello，世界"{
		runes = append(runes, r)
	}
	fmt.Printf("%q\n", runes)

	var x, y []int
	for i := 0; i < 10; i++{
		y = appendInt(x, i)
		fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
	}

	var x1 []int
	x1 = append(x1, 1)
	x1 = append(x1, 2, 3)
	x1 = append(x1, 4, 5, 6)
	x1 = append(x1, x1...) // append the slice x
	fmt.Println(x1)

	data := []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty(data))
	fmt.Printf("%q\n", data)

	s2 := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s2, 2))
	s3 := []int{5, 6, 7, 8, 9}
	fmt.Println(remove1(s3, 2))

	//练习 4.3： 重写reverse函数，使用数组指针代替slice。
	aP := &([]int{0, 1, 2, 3, 4, 5, 6})
	reversePoint(aP)
	fmt.Println("aP = ", aP)

	//练习 4.4： 编写一个rotate函数，通过一次循环完成旋转。
	aR := []int{0, 1, 2, 3, 4, 5, 6}
	rotate(aR)
	fmt.Println("aR = ", aR)

	//练习 4.5： 写一个函数在原地完成消除[]string中相邻重复的字符串的操作。
	aS := []string{"1", "2", "1", "1", "3", "2", "1", "1", "3"}
	aS = RemoveDuplicate(aS)
	fmt.Println(aS)

	//练习 4.6： 编写一个函数，原地将一个UTF-8编码的[]byte类型的slice中相邻的空格（参考unicode.IsSpace）替换成一个空格返回
	aSP := []rune{' ', '/', '1', '2', ' ', ' ', ' ', '3'}
	aSP = RemoveDuplicateSpace(aSP)
	fmt.Printf("%q\n", aSP)

	//练习 4.7： 修改reverse函数用于原地反转UTF-8编码的[]byte。是否可以不用分配额外的内存？
	aUT8 := []byte("abcde")
	reversUT8(aUT8)
	fmt.Printf("%q\n", aUT8)

}
func reversUT8(s []byte)  {
	for i, j := 0, len(s) - 1; i < j; i, j = i + 1, j - 1 {
		s[i], s[j] = s[j], s[i]
	}
}
func RemoveDuplicateSpace(s []rune) []rune {
	i := 0
	for {
		if unicode.IsSpace(s[i]) && unicode.IsSpace(s[i+1]) {
			s = append(s[:i], s[i+1:]...)
			i = 0
		}else  {
			i++
		}
		if i == len(s) - 1 {
			break
		}
	}
	return s
}
func RemoveDuplicate(s []string) []string {
	i := 0
	for {
		if s[i] == s[i + 1] {
			s = append(s[:i], s[i+1:]...)
			i = 0
		}else  {
			i++
		}
		if i == len(s) - 1 {
			break
		}
	}
	return s
}
func rotate(s []int)  {
	len := len(s)
	for i := 0; i < len / 2; i++ {
		temp := s[i]
		s[i] = s[len - i - 1]
		s[len - i - 1] = temp
	}
}
func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
func remove1(slice []int, i int) []int {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func nonempty(strings []string) []string {
	out := strings[:0] // zero-length slice of original
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		// There is room to grow.  Extend the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space.  Allocate a new array.
		// Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x) // a built-in function; see text
	}
	z[len(x)] = y
	return z
}
func reverse(s []int)  {
	for i, j := 0, len(s) - 1; i < j; i, j = i + 1, j - 1 {
		s[i], s[j] = s[j], s[i]
	}
}
func reversePoint(s *[]int)  {
	for i, j := 0, len(*s) - 1; i < j; i, j = i + 1, j - 1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

func equalInt(s, s1 []int) bool {
	if len(s) != len(s1) {
		return false
	}

	for i := range s{
		if s[i] != s1[i] {
			return false
		}
	}
	return true
}

//Map
//1.Map是一个无序的key-value对的集合，可理解为字典，go中为哈希表
//2.Map的key必须为同类型，最好不要用浮点型，value不受限制
func test_map()  {
	ages1 := make(map[string]int)
	ages1["alice"] = 31
	ages1["tom"] = 34
	fmt.Println("ages1 = ", ages1)

	ages2 := map[string]int{
		"alice":31,
		"tom":34,
	}
	fmt.Println("ages2 = ", ages2)

	delete(ages2, "alice")
	fmt.Println(ages2)

	ages2["bob"] = ages2["bob"] + 1
	fmt.Println(ages2)

	for name, age := range ages2{
		fmt.Printf("%s\t%d\n", name, age)
	}

	var names []string
	for name := range ages2 {
		names = append(names, name)
	}

	//3.map的key是无序的，需要对keys进行排序才能打印有序map
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages2[name])
	}

	var ages map[string]int
	fmt.Println(ages == nil)    // "true"
	fmt.Println(len(ages) == 0) // "true"

	//4.不能向一个为nil的map中添加key-value，因此在向map存数据前必须先创建map。
	//ages["carol"] = 21 // panic: assignment to entry in nil map
	//5.map的元素不能使用取地址操作，因为map会因为云阿苏增加而改变内存地址，重新分配内存空间
	//_ = &ages["bob"] // compile error: cannot take address of map element

	ok := equal(map[string]int{"A": 0}, map[string]int{"B": 42})
	fmt.Println(ok)

	seen := make(map[string]bool)
	input := []string{"a", "b", "c", "d", "a", "b", "d"}
	for _, line := range input{
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}
	fmt.Println(seen)

	fmt.Println(m)

	m[k(input)] = 123
	fmt.Println(m)

	//unicodeCount()

	//练习 4.9： 编写一个程序wordfreq程序，报告输入文本中每个单词出现的频率。
	//在第一次调用Scan前先调用input.Split(bufio.ScanWords)函数，这样可以按单词而不是按行输入。
	wordfreq()
}
func wordfreq()  {
	counts := make(map[string]int)

	scan := bufio.NewScanner(os.Stdin)
	scan.Split(bufio.ScanWords)

	fmt.Printf("|words        |count\n\n")
	for scan.Scan() {
		ss := scan.Text()

		if counts[ss] >= 0 {
			counts [ss]++
		}
		for key, value := range counts{

			fmt.Printf("|%8s        |%d\n\n", key, value)
		}
	}
}
//跟踪出现过字符的次数
func unicodeCount() {
	counts := make(map[rune]int)    // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	letters := make(map[rune]int)    // counts of Unicode characters
	numbers := make(map[rune]int)    // counts of Unicode characters

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error

		//回车结束输入
		if r == '\n' {
			break
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}

		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		//练习 4.8： 修改charcount程序，使用unicode.IsLetter等相关的函数，统计字母、数字等Unicode中不同的字符类别。
		if unicode.IsLetter(r) {
			letters[r]++
		}
		if unicode.IsNumber(r) {
			numbers[r]++
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for c, n := range letters {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tletter\n")
	for c, n := range numbers {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tnumber\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}
var m = make(map[string]int)
func k(list []string) string {
	return fmt.Sprintf("%q", list)
}
func Add(list []string)       {
	m[k(list)]++
}
func Count(list []string) int {
	return m[k(list)]
}
func equal(x, y map[string]int) bool  {
	if len(x) != len(y) {
		return false
	}
	//6.map和slice一样，不能使用==操作，唯一能用==操作的是和nil比较时
	//7.ok 是map在取值value时产生，标志着此元素是否真的存在
	for k, xvalue := range x{
		if yv, ok := y[k]; !ok || yv != xvalue {
			return false
		}
	}

	return true
}

//结构体
//1.结构体是一种聚合的数据类型，是由零个或多个任意类型的值聚合成的实体

func test_struct()  {
	type Employee struct {
		ID        int
		Name      string
		Address   string
		DoB       time.Time
		Position  string
		Salary    int
		ManagerID int
	}

	//2.结构体变量的成员可以通过点操作符访问和赋值
	var dilbert Employee
	dilbert.Salary -= 5000
	fmt.Println(dilbert)

	//3.也可以通过取地址->通过指针访问和赋值
	position := &dilbert.Position
	*position = "Senior " + *position
	fmt.Println(dilbert, *position)

	//4.点操作和指针操作可以一起工作
	var employeeOfTheMonth *Employee = &dilbert
	employeeOfTheMonth.Position += " (proactive team player)"
	fmt.Println(dilbert)

	//和上面相等
	(*employeeOfTheMonth).Position += " (proactive team player)"
	fmt.Println(dilbert)

	//二叉树
	arr := []int{9, 8, 3, 4, 5, 6, 7}
	tree1 := Sort(arr)
	treeStr := ""
	printTree(tree1, &treeStr)
	fmt.Println(treeStr)

	//5.结构体成员均为零值时，意味着结构体零值，也是最合理的默认值
	seen := make(map[string]struct{}) // set of strings
	if _, ok := seen["key"]; !ok {
		seen["key"] = struct{}{}
	}
	fmt.Println(seen)

	//结构体成员能比较的，结构体一样可以比较，因此结构体可以作为map的key使用
	type Point struct{ X, Y int }

	p := Point{1, 2}
	q := Point{2, 1}
	fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
	fmt.Println(p == q)                   // "false"

	type address struct {
		hostname string
		port     int
	}

	hits := make(map[address]int)
	hits[address{"golang.org", 443}]++
	fmt.Println(hits)

	//8.结构体的嵌入和匿名成员，能让点语操作法访问操作变得更简洁
	type Circle struct {
		Point
		Radius int
	}

	type Wheel struct {
		Circle
		Spokes int
	}

	var w Wheel
	w.X = 8            // equivalent to w.Circle.Point.X = 8
	w.Y = 8            // equivalent to w.Circle.Point.Y = 8
	w.Radius = 5       // equivalent to w.Circle.Radius = 5
	w.Spokes = 20
	fmt.Println(w)

	//但是结构体初始化不能使用匿名成员的特性
	w = Wheel{Circle{Point{8, 8}, 5}, 20}
	//#副词，打印出成员名称
	fmt.Printf("%#v\n", w)

}
type tree struct {
	value       int
	left, right *tree
}
//6.聚合值不能包含他自身（同样适用于数组），但是能包含自身类型的指针类型
// Sort sorts values in place.
func Sort(values []int) *tree {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
	return root
}
//打印二叉树
//7.如果在函数内部修改结构体成员的话，必须用指针传入
func printTree(t *tree, s *string) {
	if t != nil {
		*s = *s + fmt.Sprint(t.value)
		if t.left != nil || t.right != nil{
			*s = *s + "("
			printTree(t.left, s)
			if t.right != nil {
				*s = *s + ","
			}
			printTree(t.right, s)
			*s = *s + ")"
		}
	}
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

//JSON
//1.JSON是对JavaScript中各种类型的值——字符串、数字、布尔值和对象——Unicode本文编码
//2.不过JSON使用的是\Uhhhh转义数字来表示一个UTF-16编码，而不是Go语言的rune类型。
func test_json()  {
	type Movie struct {
		Title  string
		Year   int  `json:"released"`
		Color  bool `json:"color,omitempty"`
		Actors []string
	}

	var movies = []Movie{
		{Title: "Casablanca", Year: 1942, Color: false,
			Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
		{Title: "Cool Hand Luke", Year: 1967, Color: true,
			Actors: []string{"Paul Newman"}},
		{Title: "Bullitt", Year: 1968, Color: true,
			Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	}
	fmt.Println(movies)

	//3.将slice转为JSON的过程叫编组，MarshalIndent函数将产生整齐缩进的输出
	//编码操作
	data, err := json.MarshalIndent(movies, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)

	//4.编码的逆操作是解码，对应将JSON数据解码为Go语言的数据结构
	//解码操作
	var titles []struct{ Title string }
	if err := json.Unmarshal(data, &titles); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	fmt.Println(titles)

	//
	searchKey := []string{"is:open json decoder"}
	result, err := github.SearchIssues(searchKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}

	//练习 4.10： 修改issues程序，根据问题的时间进行分类，比如不到一个月的、不到一年的、超过一年。
	//打印结果
	for k, v := range result.Dict{
		fmt.Printf("%s:\n", k)
		for _, item := range v{
			fmt.Printf("    #%-5d %v %9.9s %.55s\n",
				item.Number, item.CreatedAt, item.User.Login, item.Title)
		}
	}


	//练习 4.12
	xkcd_json(false)

	//练习 4.13
	//poster()
	//以上练习只能执行一个

}
//练习 4.12： 流行的web漫画服务xkcd也提供了JSON接口。
//例如，一个 https://xkcd.com/571/info.0.json 请求将返回一个很多人喜爱的571编号的详细描述。
//下载每个链接（只下载一次）然后创建一个离线索引。
//编写一个xkcd工具，使用这些离线索引，打印和命令行输入的检索词相匹配的漫画的URL。
func xkcd_json(res bool)  {

	type caton struct {
		Num int64
		Month string
		Link string
		Year string
		News string
		Safe_title string
		Transcript string
		Alt string
		Img string
		Title string
		Day string
		Url string
	}
	//创建数组
	xkcds := []*caton{}

	//打开本地文件
	if res {
		fmt.Println("重新读取本地离线资源，请稍后...")
	}else {
		fmt.Println("正在读取本地离线资源，请稍后...")
	}
	name := "test.json"
	file, err := os.Open(name)
	defer file.Close()

	//是否已经存在，并解析成功
	ok := false

	if err != nil {
		fmt.Println("json Cannot open:", err)
	}

	dec := json.NewDecoder(file)
	if err := dec.Decode(&xkcds); err != nil {
		fmt.Println("json Cannot decode:", err)
		fmt.Println("读取本地离线资源失败，准备更新资源")
	}else {
		ok = true
	}

	//存在并解析成功，将不再网络下载
	if ok {
		fmt.Println("读取完成，请输入您要查找的关键字，回车结束...")
		scan := bufio.NewScanner(os.Stdin)
		scan.Split(bufio.ScanWords)

		for scan.Scan() {
			keyWord := scan.Text()

			searchResults := []string{}
			for _, item := range xkcds {
				//搜索含有关键字的itme
				if strings.Contains(item.Title, keyWord) {
					searchResults = append(searchResults, item.Url)
				}else if strings.Contains(item.Safe_title, keyWord) {
					searchResults = append(searchResults, item.Url)
				}else if strings.Contains(item.Alt, keyWord) {
					searchResults = append(searchResults, item.Url)
				}else if strings.Contains(item.Transcript, keyWord) {
					searchResults = append(searchResults, item.Url)
				}
			}
			fmt.Println("您要查找的相关漫画链接：")
			fmt.Println(searchResults)
		}
		return
	}


	//网络下载
	fmt.Println("正在更新本地离线资源，请稍后...")
	baseUrl := "https://xkcd.com/"
	suffixUrl := "/info.0.json"

	count := 1
	maxCount := 2000
	
	for  {
		newUrl := baseUrl + fmt.Sprint(count) + suffixUrl

		resp, err := http.Get(newUrl)
		if err != nil {
			fmt.Println(err)
			break
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Printf("%v 此资源不存在，跳过\n", newUrl)
			count++
		}else {
			var result caton
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				resp.Body.Close()
				fmt.Println(err)
				break
			}
			result.Url = newUrl[:len(newUrl) - len(suffixUrl) + 1]
			resp.Body.Close()

			xkcds = append(xkcds, &result)
		}

		count++
		if count == maxCount {
			break
		}
	}

	//保存本地
	newFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)

	defer newFile.Close()
	if err != nil {
		fmt.Println(err)
	}

	enc := json.NewEncoder(newFile)
	if err := enc.Encode(xkcds); err != nil {
		fmt.Println("json Cannot encode:", err)
	}
	fmt.Println("本地离线资源更新完成！！！")
	xkcd_json(true)
}
//练习 4.13： 使用开放电影数据库的JSON服务接口，允许你检索和下载 https://omdbapi.com/ 上电影的名字和对应的海报图像。
//编写一个poster工具，通过命令行输入的电影名字，下载对应的海报。
func poster()  {
	const (
		APIKEY = "1f0bbfee"
		BaseUrl = "http://www.omdbapi.com/?apikey=" + APIKEY + "&t="
	)
	type rating struct {
		Source string
		Value string
	} 
	
	type movie struct {
		Title string
		Year string
		Rated string
		Released string
		Runtime string
		Genre string
		Director string
		Writer string
		Actors string
		Plot string
		Language string
		Country string
		Awards string
		Poster string
		Ratings []*rating
		Metascore string
		ImdbRating string
		ImdbVotes string
		ImdbID string
		Type string
		DVD string
		BoxOffice string
		Production string
		Website string
		Response string
	}

	fmt.Println("请输入您要查找电影的关键字，回车结束")
	scan := bufio.NewScanner(os.Stdin)
	scan.Split(bufio.ScanWords)

	for scan.Scan() {
		keyWord := scan.Text()

		resp, err := http.Get(BaseUrl + keyWord)
		if err != nil {
			fmt.Println(err)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%v 此资源不存在\n", keyWord)
		}


		var result movie
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()

		fmt.Println(result.Poster)
		saveImages(result.Poster)
	}
}
//下载图片
func saveImages(img_url string){
	log.Println(img_url)
	u, err := url.Parse(img_url)
	if err != nil {
		log.Println("parse url failed:", img_url, err)
		return
	}
	fmt.Println("正在下载图片...")
	response, err := http.Get(img_url)
	if err != nil {
		log.Println("get img_url failed:", err)
		return
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("read data failed:", img_url, err)
		return
	}

	//保存本地
	//去掉最左边的'/'
	tmp := strings.TrimLeft(u.Path, "/")
	filename := strings.ToLower(strings.Replace(tmp, "/", "-", -1))

	newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	defer newFile.Close()
	if err != nil {
		fmt.Println("create file failed:", filename, err)
		return
	}

	newFile.Write(data)
	fmt.Println("图片保存完成")
}


//文本和HTML模板

func test_txt()  {
	const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
{{end}}`

	report, err := template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(report)

	report = template.Must(template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ))

	searchKey := []string{"is:open json decoder"}

	result, err := github.SearchIssues(searchKey)
	if err != nil {
		log.Fatal(err)
	}

	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}
var issueList = template.Must(template.New("issuelist").Parse(`
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .Items}}
<tr>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`))
func test_html(fileName string, searchKey []string)  {

	result, err := github.SearchIssues(searchKey)
	if err != nil {
		log.Fatal(err)
	}

	newFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	defer newFile.Close()
	if err != nil {
		fmt.Println("create file failed:", fileName, err)
		return
	}

	if err := issueList.Execute(newFile, result); err != nil {
		log.Fatal(err)
	}
}

func test_autoescape()  {
	const templ = `<p>A: {{.A}}</p><p>B: {{.B}}</p>`
	t := template.Must(template.New("escape").Parse(templ))
	var data struct {
		A string        // untrusted plain text
		B template.HTML // trusted HTML
	}
	data.A = "<b>Hello!</b>"
	data.B = "<b>Hello!</b>"

	newFile, err := os.OpenFile("test_autoescape.html", os.O_RDWR|os.O_CREATE, 0666)
	defer newFile.Close()
	if err != nil {
		fmt.Println("create file failed:", "test_autoescape.html", err)
		return
	}

	if err := t.Execute(newFile, data); err != nil {
		log.Fatal(err)
	}
}
//练习 4.14： 创建一个web服务器，查询一次GitHub，然后生成BUG报告、里程碑和对应的用户信息。
func test_GitHub()  {

	handl := func(w http.ResponseWriter, r *http.Request) {
		keyWords := r.URL.Path[1:]
		fmt.Println(keyWords)

		if len(keyWords) <= 0 {
			return
		}
		searchKey := []string{keyWords}

		result, err := github.SearchIssues(searchKey)
		if err != nil {
			log.Fatal(err)
		}
		if err := issueList.Execute(w, result); err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/", handl)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func main() {

	//数组
	//test_array()

	//Slice
	//test_slice()

	//Map
	//test_map()

	//结构体
	//test_struct()

	//JSON
	//test_json()

	//文本和HTML模板
	//打印
	test_txt()
	//生成本地html
	test_html("issues.html", []string{"commenter:gopherbot json encoder"})
	test_html("issues2.html", []string{"3133", "10535"})
	test_autoescape()
	//web服务器
	test_GitHub()
}
