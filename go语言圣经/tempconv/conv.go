package tempconv


// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c * 9 / 5 + 32)
}

// KToF converts a Kelvin temperature to Fahrenheit.
func KToF(k Kelvin) Fahrenheit {
	return Fahrenheit((Celsius(k) + AbsoluteZeroC) * 9 / 5 + 32)
}

// FToK converts a Fahrenheit temperature to Kelvin.
func FToK(f Fahrenheit) Kelvin {
	return Kelvin((f - 32) * 5 / 9 - Fahrenheit(AbsoluteZeroC))
}

// CToK converts a Celsius temperature to Kelvin.
func CToK(c Celsius) Kelvin {
	return Kelvin(c - AbsoluteZeroC)
}

// KToC converts a Kelvin temperature to Celsius.
func KToC(k Kelvin) Celsius {
	return Celsius(k * 9 / 5 + 32)
}