package geometry

import "math"

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