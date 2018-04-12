package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

/* 复合数据类型 */

//数组
func test_array()  {

	//数组的每个元素可以通过下标来访问，从0到len - 1
	//默认情况下，数组元素初始化为元素类型的零值，比如数字元素为0
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

	//数组初始化时如果未给长度，会根据元素个数来计算
	qq := [...]int{1, 2, 3, 4}
	fmt.Printf("%T\n", qq)

	//数组的长度在编译阶段就确定了
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


func main() {

	//数组
	test_array()
}
