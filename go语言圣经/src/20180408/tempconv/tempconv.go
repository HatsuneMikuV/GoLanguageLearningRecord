package tempconv

import (
	"flag"
	"fmt"
)

type Celsius float64 //摄氏温度
type Fahrenheit float64 //华氏温度
type Kelvin float64 //绝对温度

const (
	AbsoluteZeroC Celsius = -273.15 // 絶對零度
	FreezingCCC Celsius = 0 // 結冰點溫度
	BoilingCCC Celsius = 100 // 沸水溫度
)
func (c Celsius) String() string {
	return fmt.Sprintf("%g°C", c)
}
func (f Fahrenheit) String() string {
	return fmt.Sprintf("%g°F", f)
}
func (k Kelvin) String() string {
	return fmt.Sprintf("%g°K", k)
}

type celsiusFlag struct{ Celsius }
//练习 7.6： 对tempFlag加入支持开尔文温度。
type fahrenheitFlag struct{ Fahrenheit }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}

	return fmt.Errorf("invalid temperature %q", s)
}
//练习 7.6： 对tempFlag加入支持开尔文温度。

func (f *fahrenheitFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
	switch unit {
	case "F", "°F":
		f.Fahrenheit = Fahrenheit(value)
		return nil
	case "C", "°C":
		f.Fahrenheit = CToF(Celsius(value))
		return nil
	}

	return fmt.Errorf("invalid temperature %q", s)
}

func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
//练习 7.6： 对tempFlag加入支持开尔文温度。

func FahrenheitFlag(name string, value Fahrenheit, usage string) *Fahrenheit {
	f := fahrenheitFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Fahrenheit
}