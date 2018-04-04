package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"time"
)

//映射
const (
	China int = iota
	USA
	Japan
	Total
)

type GDP struct {
	k string
	v float64
}

func mapping()  {
	GDPof := make([]GDP, Total)
	GDPof[China] = GDP{"China", 5.92}
	GDPof[Japan] = GDP{"Japan", 5.45}
	GDPof[USA]   = GDP{"USA", 14.58}

	for k, v := range GDPof{
		fmt.Printf("%d. %s:%g\n", k, v.k, v.v)
	}
}

func mapping_test()  {
	GDPof := map[string]float64 {
		"USA": 14.58,
		"Japan": 5.45,
	}
	GDPof["China"] = 5.92

	for k, v := range GDPof{
		fmt.Printf("%s:%g\n", k, v)
	}
}

type Lockmap struct {
	sync.Mutex
	m map[string]int
}

var lockmap = make([]Lockmap, 10)

func init()  {
	for i := range lockmap{
		lockmap[i].m = make(map[string]int)
	}
}

func counter(s string)  {
	i := int(s[0] - '0')
	time.Sleep(time.Duration(i) * time.Millisecond)

	lockmap[i].Lock()
	defer  lockmap[i].Unlock()
	lockmap[i].m[s]++
}
func mapping_test_test()  {
	for i := 0; i < 20; i++ {
		r := fmt.Sprintf("%d", rand.Uint32())
		go  counter(r)
	}
	time.Sleep(time.Second)
	for i := range lockmap {
		for k, v := range  lockmap[i].m{
			fmt.Printf("[%d] %s = %d\n", i, k, v)
		}
	}
}

//界面类型
//界面值
type Duck float32
func (Duck) Quack() string  {
	return "嘎"
}

type Educk complex128
func (Educk) Quack() string  {
	return "叮咚"
}

type Quacker interface {
	Quack() string
}

func interface_test()  {
	var d Duck = 0.
	var e Educk = 0i
	var q Quacker

	fmt.Println(d, e, q)
	q = d
	fmt.Println(q.Quack())
	q = e
	fmt.Println(q.Quack())
}

//error界面
type Err struct {}

func (_ *Err) Error()string  {
	return  "To err is human"
}

func ToErr(ok bool) error  {
	var e *Err = nil
	if ok {
		e = &Err{}
	}
	return e
}

func NoErr(ok bool) error  {
	if !ok {
		return &Err{}
	}
	return nil
}

func error_test()  {
	fmt.Println(ToErr(true))
	fmt.Println(ToErr(false))
	fmt.Println(NoErr(true))
	fmt.Println(NoErr(false))

}

//有界无类
type Goose struct {Duck}

func Goose_test()  {
	var d Duck = 0.
	var q Quacker
	g := Goose{d}
	g.Quack() //g.Quack() undefined
	q = g //Goose does not implement Quacker
	fmt.Println(q.Quack())
}

//排序
func sort_test()  {
	is := sort.IntSlice{3, 1, 4, 1, 5, 9, 2, 6}
	ss := sort.StringSlice{"士", "农", "工", "商"}
	fs := sort.Float64Slice{math.Inf(-1), math.Inf(+1), math.NaN()}

	sort.Sort(is)
	sort.Sort(ss)
	sort.Sort(fs)

	fmt.Println(is, ss, fs)
}

type Vec struct {
	x float64
	y float64
	z float64
}

func (v Vec) lenlen() float64  {
	return math.Sqrt((v.x * v.x + v.y * v.y + v.z * v.z))
}

type VecSlice []Vec

func (v VecSlice) Len() int  {
	return len(v)
}

func (v VecSlice) Less(i, j int) bool  {
	return v[i].lenlen() < v[j].lenlen()
}

func (v VecSlice) Swap(i, j int)  {
	v[i], v[j] = v[j], v[i]
}

type Rev struct {
	sort.Interface
}

func (r Rev) Less(i, j int)bool  {
	return r.Interface.Less(j, i)
}

func vec_test()  {
	v0 := Vec{0, 0, 0}
	v1 := Vec{1, 0, 0}
	v2 := Vec{1, 1, 1}
	vs := VecSlice{v1, v2, v0}
	fmt.Println(vs)

	sort.Sort(vs)
	fmt.Println(vs)

	sort.Sort(Rev{vs})
	fmt.Println(vs)
}

//类型断言
type (
	T0 	[]string
	T00 []string
	T1 	struct {a, b int}
	T11 struct {a, b int}
	T2 	func(int, float64) *T0
	T22 func(int, float64) *[]string
)

var (
	t interface{}
	t0 T0
	t1 T1
	t2 T2
)

func tt_test()  {
	t = t0
	{
		v, ok := t.(T0)
		fmt.Println(v, ok)
	}
	{
		v, ok := t.([]string)
		fmt.Println(v, ok)
	}
	{
		v, ok := t.(T00)
		fmt.Println(v, ok)
	}

	t = t1
	{
		v, ok := t.(T1)
		fmt.Println(v, ok)
	}
	{
		v, ok := t.(struct {a, b int})
		fmt.Println(v, ok)
	}
	{
		v, ok := t.(T11)
		fmt.Println(v, ok)
	}

	t = t2
	{
		v, ok := t.(T2)
		fmt.Println(v, ok)
	}
	{
		v, ok := t.(T22)
		fmt.Println(v, ok)
	}
}

//类型分支
func class_test(x interface{})  {

	switch i := x.(type) {
	case nil:
		fmt.Println("x is nil")
	case int:
		fmt.Println("i is an int")
	case float64:
		fmt.Println("i is a float64")
	case func(int) float64:
		fmt.Println("i is a function")
	case bool, string:
		fmt.Println("type is bool or string i is an interface")
	default:
		fmt.Println("don't know the type of ", i)
	}
}

//反射
type X struct {
	Y byte
	Z complex128
}

func (x *X) String() string  {
	return fmt.Sprintf("%v", x)
}

func reflect_test()  {
	var x float64 = 3.14
	fmt.Println("type:", reflect.TypeOf(x))
	fmt.Println("value:", reflect.ValueOf(x))

	var xx X
	fmt.Println(xx)

	v := reflect.TypeOf(xx)
	fmt.Println("type:", v)
	fmt.Println("Align:", v.Align())
	fmt.Println("FieldAlign:", v.FieldAlign())

	for i := 0; i < v.NumMethod(); i++ {
		fmt.Println("Method", i, v.Method(i).Name)
	}
	fmt.Println("Name", v.Name())
	fmt.Println("PkgPath", v.PkgPath())
	fmt.Println("Kind", v.Kind())

	for i := 0; i < v.NumMethod(); i++ {
		fmt.Println("Field", i, v.Field(i).Name)
	}
	fmt.Println("Size", v.Size())

	vv := reflect.ValueOf(x)
	fmt.Println("type:", vv.Type())
	fmt.Println("kind is float64:", vv.Kind() == reflect.Float64)
	fmt.Println("value:", vv.Float())
}

func reflect_test_one()  {
	var x float64 = 3.14
	fmt.Println("x:", x)
	v := reflect.ValueOf(&x)
	fmt.Println("v.CanSet:", v.CanSet())
	i := v.Elem()
	fmt.Println("i.CanSet:", i.CanSet())
	i.SetFloat(3.1415)
	fmt.Println("i:", i.Float())
	fmt.Println("x:", x)
}

func reflect_test_two()  {
	var x X
	fmt.Println(x)

	v := reflect.ValueOf(&x)
	e := v.Elem()
	e.Field(0).SetUint('e')

	π := math.Pi
	c := complex(math.Cos(π), math.Sin(π))
	e.Field(1).SetComplex(c)

	fmt.Println("最美的公式")
	fmt.Printf("%c^iπ + 1 = %g\n", x.Y, 1+real(x.Z))
	fmt.Println("可惜虚部有点精度误差：", imag(x.Z))
}

func main() {

	//映射
	mapping()
	mapping_test()
	mapping_test_test()

	//界面类型
	//界面值
	interface_test()

	//error界面
	error_test()

	//有界无类
	Goose_test()

	//排序
	sort_test()
	vec_test()

	//类型断言
	tt_test()

	//类型分支
	class_test("1")

	//反射
	reflect_test()
	reflect_test_one()
	reflect_test_two()
}

