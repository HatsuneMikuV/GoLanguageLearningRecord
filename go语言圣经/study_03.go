package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/* 基础数据类型 */


//整型
func test_one03()  {
	var u uint8 = 255
	fmt.Println(u, u + 1, u * u)

	var i int8 = 127
	fmt.Println(i, i + 1, i * i)

	fmt.Printf("%08d\n",233)

	var x uint8 = 1<<1 | 1<<5
	var y uint8 = 1<<1 | 1<<2
	fmt.Printf("%08b\n", x) // "00100010", the set {1, 5}
	fmt.Printf("%08b\n", y) // "00000110", the set {1, 2}
	fmt.Printf("%08b\n", x&y) // "00000010", the intersection {1}
	fmt.Printf("%08b\n", x|y) // "00100110", the union {1, 2, 5}
	fmt.Printf("%08b\n", x^y) // "00100100", the symmetric difference {2, 5}
	fmt.Printf("%08b\n", x&^y) // "00100000", the difference {5}
	for i := uint(0); i < 8; i++ {
		if x&(1<<i) != 0 { // membership test
			fmt.Println(i) // "1", "5"
		}
	}
	fmt.Printf("%08b\n", x<<1) // "01000100", the set {2, 6}
	fmt.Printf("%08b\n", x>>1) // "00010001", the set {0, 4}

	medals := []string {"gold", "silver", "bronze"}
	for i := len(medals) - 1; i >= 0; i-- {
		fmt.Println(medals[i])
	}
	var apples int32 = 1
	var oranges int16 = 2
	var compote int = int(apples) + int(oranges)
	fmt.Println(compote)

	f := 3.14
	ii := int(f)
	fmt.Println(f, "\n", ii)
	f = 1.99
	fmt.Println(int(f))

	ff := 1e100
	iii := int(ff)
	fmt.Println(ff, iii)

	o := 0666
	fmt.Printf("%d %[1]o %#[1]o\n", o)
	xx := int64(0xdeadbeef)
	fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", xx)

	ascii := 'a'
	unicode := '国'
	newline := '\n'
	fmt.Printf("%d %[1]c %[1]q\n", ascii)
	fmt.Printf("%d %[1]c %[1]q\n", unicode)
	fmt.Printf("%d %[1]q\n", newline)
}


//浮点数
func test_two03()  {
	var f float32 = 16777216 //1 << 24
	fmt.Println(f == f + 1)

	const e = 2.71828
	const Avogadro = 6.02214129e23
	const Planck = 6.62606957e-34

	for x := 0; x < 8; x++  {
		fmt.Printf("x = %d e^x = %8.3f\n", x, math.Exp(float64(x)))
	}

	var z float64
	fmt.Println(z, -z, 1/z, -1/z, z/z)

	nan := math.NaN()
	fmt.Println(nan == nan, nan < nan, nan > nan)


	//保存本地的svg图
	s := getSvg()

	fileName := "SVG.svg"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	dstFile.WriteString(s)

	//网页打开的svg图 http://localhost:1234/
	handle := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		if err := r.ParseForm(); err != nil{
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, getSvg())
	}
	http.HandleFunc("/",handle)
	log.Fatal(http.ListenAndServe("localhost:1234", nil))
}
const (
	width, height = 600, 320 // canvas size in pixels
	cells = 100  // number of grid cells
	xyrange = 30.0 // axis ranges (-xyrange..+xyrange)
	xyscale = width / 2 / xyrange // pixels per x or y unit
	zscale = height * 0.4 // pixels per z unit
	angle = math.Pi / 6 // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

//练习 5.6： 修改gopl.io/ch3/surface (§3.2) 中的corner函数，将返回值命名，并使用bare return。
func corner(i, j int) (sx, sy float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i) / cells - 0.5)
	y := xyrange * (float64(j) / cells - 0.5)
	// Compute surface height z.
	z := f(x, y)
	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx = width / 2 + (x - y) * cos30 * xyscale
	sy = height / 2 + (x + y) * sin30 * xyscale - z * zscale

	return
}
func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	if math.IsNaN(r) || math.IsInf(r, 0) {
		return 0
	}
	return math.Sin(r) / r
}
func getSvg()(svg string)  {

	s := fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: #ff0000; fill: #0000ff; stroke-width: 0.3' "+
		"width='%d' height='%d'>" , width, height )

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)

			s += fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",ax, ay, bx, by, cx, cy, dx, dy)
		}
	}

	s += fmt.Sprintf("</svg>")

	return s
}

//复数
func test_thr03()  {
	var x complex128 = complex(1, 2)
	var y complex128 = complex(3, 4)
	fmt.Println(x * y)
	fmt.Println(real(x * y))
	fmt.Println(imag(x * y))

	fmt.Println(1i * 1i)

	xx := 1 + 2i
	yy := 3 + 4i
	fmt.Println(xx * yy)
	fmt.Println(real(xx * yy))
	fmt.Println(imag(xx * yy))

	fmt.Println(cmplx.Sqrt(-1))


	//绘制Mandelbrot图像
	getMandelbrot("mandelbrot1.png", 1)//黑白
	getMandelbrot("mandelbrot2.png", 2)//彩色
	getMandelbrot("mandelbrot3.png", 3)//网格
	getMandelbrot("mandelbrot4.png", 4)//彩色
	getMandelbrot("mandelbrot5.png", 5)//彩色
	getMandelbrot("mandelbrot6.png", 6)//彩色
	getMandelbrot("mandelbrot7.png", 7)//彩色
	getMandelbrot("mandelbrot8.png", 8)//彩色
	getMandelbrot("mandelbrot9.png", 9)//彩色
}
func getMandelbrot(fileName string, intType int)  {
	//绘制Mandelbrot图像
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		yyy := float64(py) / height * (ymax - ymin) + ymin
		for px := 0; px < width; px++ {
			xxx := float64(px) / width * (xmax - xmin) + xmin
			z := complex(xxx, yyy)
			if intType == 1 {
				img.Set(px, py, mandelbrot(z))
			}else if intType == 2 {
				img.Set(px, py, sqrt(z))
			}else if intType == 3 {
				img.Set(px, py, newton(z))
			}else if intType == 4 {
				img.Set(px, py, acos(z))
			}else if intType == 5 {
				img.Set(px, py, aCosh(z))
			}else if intType == 6 {
				img.Set(px, py, asin(z))
			}else if intType == 7 {
				img.Set(px, py, aSinh(z))
			}else if intType == 8 {
				img.Set(px, py, atan(z))
			}else if intType == 9 {
				img.Set(px, py, aTanh(z))
			}
		}
	}
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	png.Encode(file, img)
}
func mandelbrot(z complex128) color.Color  {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v * v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast * n}
		}
	}
	return color.Black
}
func acos(z complex128) color.Color {
	v := cmplx.Acos(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}
func aCosh(z complex128) color.Color {
	v := cmplx.Acosh(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}
func asin(z complex128) color.Color {
	v := cmplx.Asin(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}
func aSinh(z complex128) color.Color {
	v := cmplx.Asinh(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}
func atan(z complex128) color.Color {
	v := cmplx.Atan(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}
func aTanh(z complex128) color.Color {
	v := cmplx.Atanh(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}

func sqrt(z complex128) color.Color {
	v := cmplx.Sqrt(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{128, blue, red}
}

func newton(z complex128) color.Color {
	const iterations = 37
	const contrast = 7
	for i := uint8(0); i < iterations; i++ {
		z -= (z - 1/(z*z*z)) / 4
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			return color.Gray{255 - contrast*i}
		}
	}
	return color.Black
}

//布尔型
func test_fou03()  {
	var s string
	if s != "" && s[0] == 'x' {
		fmt.Println(s)
	}else {
		fmt.Println("s is nil")
	}

	var c int
	if 'a' <= c && c <= 'z' ||
		'A' <= c && c <= 'Z' ||
		'0' <= c && c <= '9'{
		fmt.Println(c)
	}else {
		fmt.Println("c is not in [a, z][A, Z][0, 9]")
	}
}

//字符串
func test_fiv03()  {
	s := "hello,world"
	fmt.Println(len(s))
	fmt.Println(s[0],s[7])

	fmt.Println(s[0:5])
	fmt.Println(s[:5])
	fmt.Println(s[7:])
	fmt.Println(s[7:len(s)])
	fmt.Println(s[:])
	fmt.Println("goodbye" + s[5:])

	ss := "left foot"
	t := ss
	ss += ",right foot"
	fmt.Println(t)
	fmt.Println(ss)

	sss := "Hello，世界"
	fmt.Println(len(sss))
	fmt.Println(utf8.RuneCountInString(sss))

	for i, r := range sss {
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}

	ssss := "プログラム"
	fmt.Printf("% x\n", ssss)
	r := []rune(ssss)
	fmt.Printf("%x\n", r)

	fmt.Println(string(r))

	fmt.Println(string(65))
	fmt.Println(string(0x4eac))
	fmt.Println(string(1234567))

	fmt.Println(basename("a/b/c.go"))
	fmt.Println(basename1("a/b/c.go"))

	fmt.Println(basename("c.d.go"))
	fmt.Println(basename1("c.d.go"))

	fmt.Println(basename("abc"))
	fmt.Println(basename1("abc"))

	fmt.Println(comma("1234567890"))//整型
	fmt.Println(comma("123456.7890"))//浮点型

	fmt.Println(intsToString([]int{1, 2, 3}))

	//练习3.10
	//练习3.11
	fmt.Println(commaBuf("1234567890"))
	fmt.Println(commaBuf("123456.7890"))
	fmt.Println(commaBuf(".1234567890"))
	fmt.Println(commaBuf("123.4567890"))
	fmt.Println(commaBuf("-123.4567890"))

	fmt.Println(compareTwoString("123","213"))//练习3.12

	x := 123
	y := fmt.Sprintf("%d", x)
	fmt.Println(y, strconv.Itoa(x))

	fmt.Println(strconv.FormatInt(int64(x), 2))
	fmt.Println(strconv.FormatInt(int64(x), 4))

	s1 := fmt.Sprintf("x=%b", x)
	fmt.Println(s1)

	x1, err := strconv.Atoi("123")
	fmt.Println(x1, err)

	y1, err := strconv.ParseInt("123", 10, 64)
	fmt.Println(y1, err)

}
func basename(s string) string  {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}
func basename1(s string) string  {
	slash := strings.LastIndex(s, "/")
	s = s[slash+1:]
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}
func comma(s string) string  {
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		return comma(s[:dot]) + s[dot:]
	}
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n - 3]) + "," + s[n-3:]
}
func commaBuf(s string) string  {
	var buf bytes.Buffer

	var sub string = ""
	if strings.Contains(s, "-") {
		sub += "-"
		s = s[1:]
	}

	dot := strings.LastIndex(s, ".")

	if dot < 0 {
		dot = len(s)
	}else if dot == 0 {
		return sub + "0" + s
	}else if dot <= 3 {
		return sub + s
	}

	buf.WriteString(sub)
	ss := s[:dot]
	if len(ss) > 3 {
		var ii int = 0
		for i := len(ss); i >= 0; i-- {
			if i%3 == 0 && i != 0 && i != len(ss) {
				buf.WriteString(ss[ii:len(ss) - i] + ",")
				ii = len(ss) - i
			}else if i == 0 || i == len(ss) {
				buf.WriteString(ss[ii:len(ss) - i])
			}
		}
	}
	if dot >= 0 {
		buf.WriteString(s[dot:])
	}

	return buf.String()
}
func intsToString(values []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}
func compareTwoString(s string, k string) bool  {
	ok := false

	if strings.EqualFold(s, k) {
		return ok
	}

	if len(s) != len(k) {
		return ok
	}

	for i := 0; i < len(s); i++{
		 ss := string(s[i])
		if strings.Contains(k, ss) {
			index := strings.Index(k, ss)
			k = k[:index] + k[index + 1:]
		}
	}

	if len(k) == 0 {
		ok = true
	}
	return ok
}
//常量
func test_six03()  {
	const  pi  = 3.14159
	const (
		e = 2.71828182845904523536028747135266249775724709369995957496696763
		ppi = 3.14159265358979323846264338327950288419716939937510582097494459
	)

	const noDelay time.Duration = 0
	const timeout = 5 * time.Minute
	fmt.Printf("%T %[1]v\n", noDelay)
	fmt.Printf("%T %[1]v\n", timeout)
	fmt.Printf("%T %[1]v\n", time.Minute)

	const (
		a = 1
		b
		c = 2
		d
	)
	fmt.Println(a, b, c, d)

	var v Flags = FlagMulticast | FlagUp
	fmt.Printf("%06b \n%06b \n%06b \n%06b \n%06b\n", FlagUp, FlagBroadcast, FlagLoopback, FlagPointToPoint, FlagMulticast)
	fmt.Printf("%b %t\n", v, IsUp(v)) // "10001 true"
	TurnDown(&v)
	fmt.Printf("%b %t\n", v, IsUp(v)) // "10000 false"
	SetBroadcast(&v)
	fmt.Printf("%b %t\n", v, IsUp(v))   // "10010 false"
	fmt.Printf("%b %t\n", v, IsCast(v)) // "10010 true"

	fmt.Println(B, KB, MB)

	fmt.Println(YiB / ZiB)

	var f float64 = 212
	fmt.Println((f - 32) * 5 / 9)
	fmt.Println(5 / 9 * (f - 32))
	fmt.Println(5.0 / 9.0 * (f - 32))

	fmt.Printf("%T\n", 0)
	fmt.Printf("%T\n", 0.0)
	fmt.Printf("%T\n", 0i)
	fmt.Printf("%T\n", '\000')
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	wednesday
	Thursday
	Friday
	Saturday
)
type Flags uint

const (
	FlagUp Flags = 1 << iota // is up
	FlagBroadcast            // supports broadcast access capability
	FlagLoopback             // is a loopback interface
	FlagPointToPoint         // belongs to a point-to-point link
	FlagMulticast            // supports multicast access capability
)
func IsUp(v Flags) bool     { return v&FlagUp == FlagUp }
func TurnDown(v *Flags)     { *v &^= FlagUp }
func SetBroadcast(v *Flags) { *v |= FlagBroadcast }
func IsCast(v Flags) bool   { return v&(FlagBroadcast|FlagMulticast) != 0 }
const (
	iB = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776             (exceeds 1 << 32)
	PiB // 1125899906842624
	EiB // 1152921504606846976
	ZiB // 1180591620717411303424    (exceeds 1 << 64)
	YiB // 1208925819614629174706176
)

const (//练习 3.13
	B  = 1e0
	KB = 1e3
	MB = 1e6
	GB = 1e9
	TB = 1e12
	PB = 1e15
	EB = 1e18
	ZB = 1e21
	YB = 1e24
)
func main()  {

	//整型
	//test_one03()

	//浮点数
	//test_two03()

	//复数
	//test_thr03()

	//布尔型
	//test_fou03()

	//字符串
	//test_fiv03()

	//常量
	test_six03()
}
