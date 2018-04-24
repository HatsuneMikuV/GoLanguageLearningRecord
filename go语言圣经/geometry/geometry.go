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


// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
//练习6.1: 为bit数组实现下面这些方法
// return the number of elements
func (s *IntSet) Len() int {
	len := 0
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				len++
			}
		}
	}
	return len
}
// remove x from the set
func (s *IntSet) Remove(x int) {
	ok := s.Has(x)
	if ok {
		word, bit := x/64, uint(x%64)
		s.words[word] &^= 1 << bit
	}
}
// remove all elements from the set
func (s *IntSet) Clear()  {
	if s.Len() > 0 {
		for i, word := range s.words {
			if word == 0 {
				continue
			}
			for j := 0; j < 64; j++ {
				if word&(1<<uint(j)) != 0 {
					s.words[i] ^= 1 << uint(j)
				}
			}
		}
	}
}
// return a copy of the set
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