package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"unicode"
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

	//使用冒泡算法求相同月份(性能较低)
	for _, sm := range summer{
		for _, qm := range Q2{
			if sm == qm {
				fmt.Printf("%s appears in both\n", sm)
			}
		}
	}

	//slice 操作长度(len)和容量(cap)时，cap不能超出范围
	//len超出意味着基于底层数组扩展，同样不能超出cap
	//fmt.Println(summer[:20]) // panic: runtime error: slice bounds out of range

	endlessSummer := summer[:5]
	fmt.Println(endlessSummer)

	a := []int{0, 1, 2, 3, 4, 5, 6}
	reverse(a)
	fmt.Println(a)

	s := []int{0, 1, 2, 3, 4, 5}

	//slice 不能使用==的操作判断，因为slice的元素是间接引用的，唯一能操作的是==nil
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

func main() {

	//数组
	//test_array()

	//Slice
	test_slice()
}
