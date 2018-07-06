package params

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Unpack populates the fields of the struct pointed to by ptr
// from the HTTP request parameters in req.
func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	// Build map of fields keyed by effective name.
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	// Update struct field for each parameter in the request.
	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue // ignore unrecognized HTTP parameters
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String, reflect.Interface:
		v.SetString(value)

	case reflect.Ptr:
		populate(v.Elem(), value)

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := populate(v.Index(i), value); err != nil {
				return err
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	case reflect.Struct: // ((name value) ...)
		for i := 0; i < v.NumField(); i++ {
			v.SetString(v.Type().Field(i).Name)
			if err := populate(v.Field(i), value); err != nil {
				return err
			}
		}

	case reflect.Float32, reflect.Float64:
		b, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		v.SetFloat(b)

	case reflect.Map: // ((key value) ...)
		for _, key := range v.MapKeys() {
			if err := populate(key, value); err != nil {
				return err
			}
			if err := populate(v.MapIndex(key), value); err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

/*
		练习 12.11：
		编写相应的Pack函数，
		给定一个结构体值，
		Pack函数将返回合并了所有结构体成员和值的URL。
*/
func Pack(reqUrl string, ptr interface{}) (Url string, err error) {
	Url = reqUrl
	if Url == "" {
		err = fmt.Errorf("reqUrl is Invalid")
		return
	}
	Url += "?"
	// Build map of fields keyed by effective name.
	fields := make(map[string]bool)
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		f := v.Field(i)

		if f.IsValid() && fields[name] {
			continue // ignore unrecognized HTTP parameters
		}
		fields[name] = true
		if i > 0{
			Url += "&"
		}
		if f.Kind() == reflect.Slice {
			for n := 0; n < f.Len(); n++ {
				ss := f.Index(n)
				Url += "&"
				Url += name + "=" +fmt.Sprint(ss)
			}
		}else {
			Url += name + "=" +fmt.Sprint(f)
		}
	}

	return
}
