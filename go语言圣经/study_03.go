package main

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"os"
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

func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	// Compute surface height z.
	z := f(x, y)
	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}
func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
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
}
func main()  {

	//整型
	test_one03()

	//浮点数
	//test_two03()

	//复数
	test_thr03()
}
