package myequal

import (
	"reflect"
	"unsafe"
	"fmt"
)

//!+
func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() {
		return false
	}

	// ...cycle check omitted (shown later)...

	//!-
	//!+cyclecheck
	// cycle check
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true // identical references
		}
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}
	//!-cyclecheck
	//!+
	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()

	case reflect.String:
		return x.String() == y.String()

		// ...numeric cases omitted for brevity...

		//!-
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return x.Int() == y.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return x.Uint() == y.Uint()

	case reflect.Float32, reflect.Float64:
		return x.Float() == y.Float()

	case reflect.Complex64, reflect.Complex128:
		return x.Complex() == y.Complex()
		//!+
	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()

	case reflect.Ptr, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)

	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true

		// ...struct and map cases omitted for brevity...
		//!-
	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if !equal(x.Field(i), y.Field(i), seen) {
				return false
			}
		}
		return true

	case reflect.Map:
		if x.Len() != y.Len() {
			return false
		}
		for _, k := range x.MapKeys() {
			if !equal(x.MapIndex(k), y.MapIndex(k), seen) {
				return false
			}
		}
		return true
		//!+
	}
	panic("unreachable")
}

//!-

//!+comparison
// Equal reports whether x and y are deeply equal.
//!-comparison
//
// Map keys are always compared with ==, not deeply.
// (This matters for keys containing pointers or interfaces.)
//!+comparison
func Equal(x, y interface{}) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

//!-comparison



//练习 13.1： 定义一个深比较函数，对于十亿以内的数字比较，忽略类型差异。
/* 只限数字类型 */
func NumDeepEqual(x, y interface{}) bool {
	return numdeepequal(reflect.ValueOf(x), reflect.ValueOf(y))
}


func numdeepequal(x, y reflect.Value) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}

	if !isNum(x) {
		s := fmt.Sprint(x, "  is not num type")
		panic(s)
	}
	if !isNum(y) {
		s := fmt.Sprint(y, "  is not num type")
		panic(s)
	}

	if x.Type() != y.Type() {
		xx := numchanggefloat64(x)
		yy := numchanggefloat64(y)

		return numdeepequal(reflect.ValueOf(xx), reflect.ValueOf(yy))
	}

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return x.Int() == y.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return x.Uint() == y.Uint()

	case reflect.Float32, reflect.Float64:
		return x.Float() == y.Float()
	}
	return false
}

func numchanggefloat64(x reflect.Value) float64 {
	var xx float64

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
			xx = float64(x.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		xx = float64(x.Uint())

	case reflect.Float32, reflect.Float64:
		xx = float64(x.Float())
	}

	return xx
}

func isNum(x reflect.Value) bool {

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return true

	case reflect.Float32, reflect.Float64:
		return true
	}

	return false
}

//练习 13.2： 编写一个函数，报告其参数是否为循环数据结构。
func CycleCheck(x interface{}) bool {
	seen := make(map[compcycle]bool)
	return cycleCheck(reflect.ValueOf(x), seen)
}
type compcycle struct {
	x 	unsafe.Pointer
	t   reflect.Type
}
func cycleCheck(x reflect.Value, seen map[compcycle]bool) bool {
	if !x.IsValid(){
		return false
	}

	// ...cycle check omitted (shown later)...

	//!-
	//!+cyclecheck
	// cycle check

	if x.CanAddr(){
		xptr := unsafe.Pointer(x.UnsafeAddr())
		c := compcycle{xptr,x.Type()}
		if seen[c] {
			return true // cycle
		}
		seen[c] = true
	}

	switch x.Kind() {
	case reflect.Bool:
		return false

	case reflect.String:
		return false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return false

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return false

	case reflect.Float32, reflect.Float64:
		return false

	case reflect.Complex64, reflect.Complex128:
		return false

	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return false

	case reflect.Ptr, reflect.Interface:
		return cycleCheck(x.Elem(), seen)

	case reflect.Array, reflect.Slice:
		return false

	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if cycleCheck(x.Field(i), seen) {
				return true
			}
		}
		return false

	case reflect.Map:
		return false
	}
	return false
}