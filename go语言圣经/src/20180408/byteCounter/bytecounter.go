package byteCounter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

//练习 7.1： 使用来自ByteCounter的思路，实现一个针对对单词和行数的计数器。
type WordCounter int
type LineCounter int


func (c *WordCounter) Write(p []byte) (int, error) {

	ss := strings.NewReader(string(p))
	scan := bufio.NewScanner(ss)
	scan.Split(bufio.ScanWords)
	wordC := 0
	for scan.Scan(){
		wordC++
	}
	*c += WordCounter(wordC) // convert int to ByteCounter
	return wordC, nil
}

func (c *LineCounter) Write(p []byte) (int, error) {

	ss := strings.NewReader(string(p))
	scan := bufio.NewScanner(ss)
	scan.Split(bufio.ScanLines)
	lineC := 0
	for scan.Scan(){
		lineC++
	}
	*c += LineCounter(lineC) // convert int to ByteCounter
	return lineC, nil
}

//练习 7.2： 写一个带有如下函数签名的函数CountingWriter

type CountWriter struct {
	W io.Writer
	C *int64
}
func (c *CountWriter) Write(p []byte) (int, error) {
	cc := int64(len(p))
	c.C = &cc
	return len(p), nil
}
func CountingWriter(w io.Writer) (io.Writer, *int64) {
	cc := &CountWriter{
		W: w,
	}
	cc.Write([]byte("Hello, world!"))
	return cc, cc.C
}

type Tree struct {
	value       int
	left, right *Tree
}
//练习 7.3： 为在gopl.io/ch4/treesort (§4.4)的*tree类型实现一个String方法去展示tree类型的值序列。
func (t *Tree)String() string  {
	s := ""
	if t != nil {
		s = s + fmt.Sprint(t.value)
		if t.left != nil || t.right != nil{
			s += "(left→" + t.left.String() + "," + t.right.String() + "←right)"
		}
	}
	return s
}
func appendValues(values []int, t *Tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func addTree(t *Tree, value int) *Tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(Tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = addTree(t.left, value)
	} else {
		t.right = addTree(t.right, value)
	}
	return t
}
func Sort(values []int) *Tree {
	var root *Tree
	for _, v := range values {
		root = addTree(root, v)
	}
	appendValues(values[:0], root)
	return root
}