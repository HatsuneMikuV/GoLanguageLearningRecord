package tempconv

import "fmt"

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