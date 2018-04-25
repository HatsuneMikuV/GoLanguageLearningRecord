package geometry

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"sync"
)

type Point struct{ X, Y float64 }

//传统的函数
func Distance(p, q Point) float64  {
	return math.Hypot(q.X - p.X, q.Y - p.Y)
}

//同样的事情，将其变成一个方法
//p， 作为一个参数， 作为此方法的接收器, 可理解为，向p对象发送消息Distance

//3.在Go语言中，可以随意定制接收器的名字，不像oc那样使用self
func (p Point)Distance(q Point) float64 {
	return math.Hypot(q.X - p.X, q.Y - p.Y)
}


// A Path is a journey connecting the points with straight lines.
type Path []Point
// Distance returns the distance traveled along the path.
//能够给任意类型定义方法，这是Go和其他语言不一样的地方
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}
func PathDistance(path Path) float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

//当你根据一个变量来决定调用同一个类型的哪个函数时，方法表达式就显得很有用了。
//你可以根据选择来调用接收器各不相同的方法
func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }

func (path Path) TranslateBy(offset Point, add bool) {
	var op func(p, q Point) Point
	if add {
		op = Point.Add
	} else {
		op = Point.Sub
	}
	for i := range path {
		// Call either path[i].Add(offset) or path[i].Sub(offset).
		path[i] = op(path[i], offset)
	}
}

//方法的名字是(*Point).ScaleBy。
//这里的括号是必须的，没有括号的话这个表达式可能会被理解为*(Point.ScaleBy)
func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

type ColoredPoint struct {
	Point
	Color color.RGBA
}

type ColoredPointP struct {
	*Point
	Color color.RGBA
}

var cache = struct {
	sync.Mutex
	mapping map[string]string
}{
	mapping: make(map[string]string),
}
//因为sync.Mutex字段也被嵌入到了这个struct里，其Lock和Unlock方法也就都被引入到了这个匿名结构中了，
//这让我们能够以一个简单明了的语法来对其进行加锁解锁操作
func Lookup(key string) string {
	cache.Lock()
	v := cache.mapping[key]
	cache.Unlock()
	return v
}


/*
	练习 6.5：
	我们这章定义的IntSet里的每个字都是用的uint64类型，
	但是64位的数值可能在32位的平台上不高效。
	修改程序，使其使用uint类型，这种类型对于32位平台来说更合适。
	当然了，这里我们可以不用简单粗暴地除64，可以定义一个常量来决定是用32还是64，
	这里你可能会用到平台的自动判断的一个智能表达式：32 << (^uint(0) >> 63)
*/
var ComputerNumber = (32 << (^uint(0) >> 63))
// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// 判断集合中是否含有指定元素
func (s *IntSet) Has(x int) bool {
	word, bit := x/ComputerNumber, uint(x%ComputerNumber)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// 添加指定元素到集合中
func (s *IntSet) Add(x int) {
	word, bit := x/ComputerNumber, uint(x%ComputerNumber)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

//并集：将B中A没有的元素放到A中，新A相当于A ∪ B
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}
//并集：元素在A集合或B集合出现：A | B    数学公式：A ∪ B
func (s *IntSet) AndWith(t *IntSet) *IntSet {
	newS := IntSet{}
	for i, word := range s.words {
		if i >= len(t.words) {
			break
		}
		newS.words = append(newS.words, word|t.words[i])
	}
	return &newS
}
//交集：元素在A集合B集合均出现：A & B    数学公式：A ∩ B
func (s *IntSet) IntersectWith(t *IntSet) *IntSet {
	newS := IntSet{}
	for i, word := range s.words {
		if i >= len(t.words) {
			break
		}
		newS.words = append(newS.words, word&t.words[i])
	}
	return &newS
}
//差集：元素出现在A集合，未出现在B集合：A &^ B   数学公式：A - B
func (s *IntSet) DifferenceWith(t *IntSet) *IntSet {
	newS := IntSet{}
	for i, word := range s.words {
		if i >= len(t.words) {
			break
		}else {
			newS.words = append(newS.words, word&^t.words[i])
		}
	}
	return &newS
}
//并差集：元素出现在A但没有出现在B，或者出现在B没有出现在A：A ^ B   数学公式：￢A ∪ ￢B
func (s *IntSet) SymmetricDifference(t *IntSet) *IntSet {
	newS := IntSet{}
	for i, word := range s.words {
		if i >= len(t.words) {
			break
		}else {
			newS.words = append(newS.words, word^t.words[i])
		}
	}
	return &newS
}

// 将集合转变成字符串
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < ComputerNumber; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", ComputerNumber*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
//练习6.1: 为bit数组实现下面这些方法
// 返回集合元素的个数
func (s *IntSet) Len() int {
	len := 0
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < ComputerNumber; j++ {
			if word&(1<<uint(j)) != 0 {
				len++
			}
		}
	}
	return len
}
// 从集合里面删除指定元素
func (s *IntSet) Remove(x int) {
	ok := s.Has(x)
	if ok {
		word, bit := x/ComputerNumber, uint(x%ComputerNumber)
		s.words[word] &^= 1 << bit
	}
}
// 清空集合里面的元素
func (s *IntSet) Clear()  {
	if s.Len() > 0 {
		for i, word := range s.words {
			if word == 0 {
				continue
			}
			for j := 0; j < ComputerNumber; j++ {
				if word&(1<<uint(j)) != 0 {
					s.words[i] ^= 1 << uint(j)
				}
			}
		}
	}
}
// 复制一个新的集合
func (s *IntSet) Copy() *IntSet {
	var newS IntSet
	for _, word := range s.words {
		newS.words = append(newS.words, word)
	}
	return &newS
}
//练习 6.2： 定义一个变参方法(*IntSet).AddAll(...int)，
//这个方法可以添加一组IntSet，比如s.AddAll(1,2,3)。
func (s *IntSet) AddAll(words...int) {
	for _, word := range words {
		s.Add(word)
	}
}

//练习6.4: 实现一个Elems方法，返回集合中的所有元素，用于做一些range之类的遍历操作。
func (s *IntSet) Elems()[]int {
	elems := []int{}
	if s.Len() > 0 {
		for i, word := range s.words {
			if word == 0 {
				continue
			}
			for j := 0; j < ComputerNumber; j++ {
				if word&(1<<uint(j)) != 0 {
					elems = append(elems, ComputerNumber * i + j)
				}
			}
		}
	}
	return elems
}